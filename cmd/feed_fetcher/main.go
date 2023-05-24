package main

import (
	"bytes"
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/mikkelstb/feeds"
)

var config_file string
var loop int

var logger *log.Logger
var cfg *feeds.Config
var client http.Client

func init() {
	flag.StringVar(&config_file, "config", "./config.json", "filepath for config file")
	flag.IntVar(&loop, "loop", 0, "number of hours between each new fetch")
}

/*
	Feedfetcher is a program for reading rss-feeds and storing them on a given archive
	The feeds to be read, and the archives to write, are both stored on a config file read from
*/

func main() {

	flag.Parse()

	fmt.Println("This is Feedfetcher")
	fmt.Printf("using configfile: %s\n", config_file)
	fmt.Printf("using loop interval of %d hours\n", loop)

	// Read configfile. Exit if unsuccessful
	var err error
	cfg, err = feeds.ReadConf(config_file)
	if err != nil {
		fmt.Println(err.Error())
		fmt.Println("error: config file not read, aborting")
		os.Exit(1)
	}

	// Set up logfile. Exit if unsuccessful
	logfile, err := os.OpenFile(cfg.Apps["fetcher"]["log"], os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0666)
	if err != nil {
		log.Fatal(err)
	}

	logger = log.New(logfile, "INFO: ", log.Ldate|log.Ltime|log.Lshortfile)

	//Set up http client
	client = http.Client{
		Transport: nil,
		Jar:       nil,
		Timeout:   10 * time.Second,
	}

	for {
		// Go through all sources
		for _, source_config := range cfg.Sources {
			if !source_config.Active {
				continue
			}
			logger.Printf("updating source %s\n", source_config.Name)
			updateFeeds(source_config)
		}
		logger.Printf("done update\n\n")
		if loop == 0 {
			break
		}
		logger.Printf("Sleeping for %d hours\n", loop)
		time.Sleep(time.Duration(loop) * time.Hour)
	}
	purgeDocumentServer()
	reIndex()
}

func updateFeeds(srcconfig feeds.SourceConfig) {

	source := NewSource(srcconfig)

	err := source.Process()
	if err != nil {
		logger.Println(err)
	}
	newsitems, errs := source.GetNewsitems()

	logger.Printf("Found %d articles\n", len(newsitems))
	if len(errs) > 0 {
		logger.Printf("Discarded %d articles:\n", len(errs))
		var errormap map[string]int = map[string]int{}
		for e := range errs {
			errormap[errs[e].Error()]++
		}
		for mes, count := range errormap {
			logger.Printf("%v: %d", mes, count)
		}
	}

	report := make(map[string]int, 0)

	logger.Println("writing items to repository")

	for ni := range newsitems {

		json_data, err := json.Marshal(newsitems[ni])
		if err != nil {
			logger.Fatal(err)
		}

		req, err := http.NewRequest(http.MethodPost, "http://localhost:3333/add/", bytes.NewReader(json_data))
		if err != nil {
			logger.Fatal(err)
		}

		res, err := client.Do(req)
		if err != nil {
			logger.Println(err)
			logger.Println(res.Status)
			continue
			//logger.Fatal(err)
		}
		defer res.Body.Close()
		logger.Println("status code:", res.StatusCode)
	}

	for status, freq := range report {
		logger.Printf("%s: %d\n", status, freq)
	}
}

func reIndex() {

	req, err := http.NewRequest(http.MethodGet, "http://localhost:3333/reindex/", nil)
	if err != nil {
		logger.Println(err)
	}
	res, err := client.Do(req)
	if err != nil {
		logger.Println("error during reindexing", err)
	}
	defer res.Body.Close()
	logger.Println("status code:", res.StatusCode)
}

func purgeDocumentServer() {

	req, err := http.NewRequest(http.MethodGet, "http://localhost:3333/purge/", nil)
	if err != nil {
		logger.Println(err)
	}
	res, err := client.Do(req)
	if err != nil {
		logger.Println("error during purge", err)
	}
	defer res.Body.Close()
	logger.Println("status code:", res.StatusCode)
}

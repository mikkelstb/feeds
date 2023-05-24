package main

import (
	"flag"
	"log"
	"net/http"
	"os"

	"github.com/mikkelstb/feeds"
)

var conffile = flag.String("c", "./config.json", "path for conffile")
var conf feeds.Config

func main() {
	flag.Parse()

	conf, err := feeds.ReadConf(*conffile)
	if err != nil {
		log.Fatal(err)
	}

	logfile, err := os.Create(conf.Logs["server"])
	if err != nil {
		log.Fatal(err)
	}
	log.SetOutput(logfile)

	http.HandleFunc("/", displaySearchPage)

	log.Fatal(http.ListenAndServe(":3300", nil))
}

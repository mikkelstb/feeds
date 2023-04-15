package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"
	"os"
	"strings"

	"github.com/mikkelstb/feeds"
	"github.com/mikkelstb/feeds/index"
)

var conffile = flag.String("c", "./config.json", "path for conffile")
var ii index.StandardIndex = make(index.StandardIndex)
var rep *feeds.Repository

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

	rep, err = feeds.InitRepository(conf.Repositories["archive"].Address)
	if err != nil {
		log.Fatal(err)
	}

	for _, s := range conf.Sources {
		if !s.Active {
			continue
		}

		for _, m := range []int{1, 2, 3} {

			articles, err := rep.GetArticlesByFeed(s.Id, 2023, m)
			if err != nil {
				log.Println(err)
				continue
			}

			for a := range articles {
				ii.AddDocument(index.DocumentFromNewsItem(articles[a]))
			}

		}
	}

	fmt.Println("no of terms: ", len(ii))

	http.HandleFunc("/postings/", getPostings)
	http.HandleFunc("/article/", getArticle)
	http.HandleFunc("/articles/", getArticles)

	log.Fatal(http.ListenAndServe(":3333", nil))
}

func getPostings(w http.ResponseWriter, r *http.Request) {
	fmt.Printf("got request %s\n", r.RequestURI)

	parts := strings.Split(r.RequestURI, "/")
	query, err := url.QueryUnescape(parts[len(parts)-1])
	if err != nil {
		io.WriteString(w, err.Error())
	}

	w = SetHeader(w)

	postings := ii.GetPostings(query)
	bytes, err := json.MarshalIndent(postings, "", "   ")
	if err != nil {
		io.WriteString(w, err.Error())
	} else {
		io.WriteString(w, string(bytes))
	}
}

func getArticle(w http.ResponseWriter, r *http.Request) {
	fmt.Printf("got request %s\n", r.RequestURI)

	parts := strings.Split(r.RequestURI, "/")
	query := parts[len(parts)-1]

	w = SetHeader(w)

	article, err := rep.GetArticleByID(query)
	if err != nil {
		io.WriteString(w, err.Error())
	}

	bytes, err := json.MarshalIndent(article, "", "   ")
	if err != nil {
		io.WriteString(w, err.Error())
	} else {
		io.WriteString(w, string(bytes))
	}
}

func SetHeader(w http.ResponseWriter) http.ResponseWriter {
	w.Header().Set("content-type", "application/json")
	w.Header().Set("encoding", "utf-8")
	return w
}

func getArticles(w http.ResponseWriter, r *http.Request) {
	fmt.Printf("got request %s\n", r.RequestURI)

	parts := strings.Split(r.RequestURI, "/")
	query, err := url.QueryUnescape(parts[len(parts)-1])
	if err != nil {
		io.WriteString(w, err.Error())
	}

	w = SetHeader(w)

	postings := ii.GetPostings(query)
	articles := make([]feeds.NewsItem, len(postings))

	for p := range postings {
		article, err := rep.GetArticleByID(postings[p].String())
		if err == nil {
			articles[p] = *article
		}
	}

	bytes, err := json.MarshalIndent(articles, "", "   ")
	if err != nil {
		io.WriteString(w, err.Error())
	} else {
		io.WriteString(w, string(bytes))
	}
}

package main

import (
	"encoding/json"
	"flag"
	"io"
	"log"
	"net/http"
	"net/url"
	"os"
	"sort"
	"strings"
	"time"

	"github.com/mikkelstb/feeds"
	"github.com/mikkelstb/feeds/vocabulary"
)

var conffile = flag.String("config", "./config.json", "path for conffile")

var indices map[string]*ExtendedIndex

// var ii index.StandardIndex
var rep *feeds.Repository
var conf *feeds.Config
var latest_items NewsQueue
var logger *log.Logger

func main() {
	flag.Parse()
	var err error
	conf, err = feeds.ReadConf(*conffile)
	if err != nil {
		logger.Fatal(err)
	}

	logfile, err := os.OpenFile(conf.Apps["server"]["log"], os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0666)
	if err != nil {
		logger.Fatal(err)
	}

	logger = log.New(logfile, "INFO: ", log.Ldate|log.Ltime|log.Lshortfile)
	logger.Println("starting document server")

	rep, err = feeds.InitRepository(conf.Repositories["archive"].Address)
	if err != nil {
		logger.Fatal(err)
	}

	// Init newsqueues
	latest_items = NewNewsQueue(20)

	//start_time := time.Now()
	initIndices()
	//index_duration := time.Since(start_time)

	//logger.Printf("Indexed %d articles in %f seconds\n", n_articles, index_duration.Seconds())
	//logger.Println("Vocabulary size:", ii.Size())

	http.HandleFunc("/postings/", getPostings)
	http.HandleFunc("/article/", getArticle)
	http.HandleFunc("/articles/", getArticles)
	http.HandleFunc("/spelling/", getSpellingSuggestion)
	http.HandleFunc("/reindex/", reindex)
	http.HandleFunc("/purge/", purgeNewsQueue)
	http.HandleFunc("/add/", addArticle)
	http.HandleFunc("/latest/", getLatest)

	logger.Fatal(http.ListenAndServe(":3333", nil))
}

/*
Returns json encoded array of docids
*/
func getPostings(w http.ResponseWriter, r *http.Request) {
	logger.Printf("got request %s\n", r.RequestURI)

	parts := strings.Split(r.RequestURI, "/")
	query, err := url.QueryUnescape(parts[len(parts)-1])
	if err != nil {
		io.WriteString(w, err.Error())
	}

	indexname, err := url.QueryUnescape(parts[len(parts)-2])
	if err != nil {
		io.WriteString(w, err.Error())
	}

	if !checkIndexName(indexname) {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	w = SetHeader(w)

	postings := indices[indexname].Index.GetPostings(query)
	bytes, err := json.MarshalIndent(postings, "", "   ")
	if err != nil {
		io.WriteString(w, err.Error())
	} else {
		io.WriteString(w, string(bytes))
	}
}

/*
Returns a single NewsItem based on DocID
*/
func getArticle(w http.ResponseWriter, r *http.Request) {
	logger.Printf("got request %s\n", r.RequestURI)

	parts := strings.Split(r.RequestURI, "/")
	query := parts[len(parts)-1]

	w = SetHeader(w)

	article, err := rep.GetArticleByID(query)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		io.WriteString(w, err.Error())
	}

	//bytes, err := json.MarshalIndent(article, "", "   ")
	bytes, err := article.MarshalJSON()
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

/*
GetArticles() is used for search queries
Function uses last part of url as search query,
and Html-escapes the query before calling the getPostings() function for the index
Returns json encoded list of feeds.Newsitem
*/
func getArticles(w http.ResponseWriter, r *http.Request) {
	logger.Printf("got request %s\n", r.RequestURI)

	parts := strings.Split(r.RequestURI, "/")
	query, err := url.QueryUnescape(parts[len(parts)-1])
	if err != nil {
		io.WriteString(w, err.Error())
	}

	indexname, err := url.QueryUnescape(parts[len(parts)-2])
	if err != nil {
		io.WriteString(w, err.Error())
	}

	if !checkIndexName(indexname) {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	w = SetHeader(w)

	postings := indices[indexname].Index.GetPostings(query)
	articles := make([]feeds.NewsItem, len(postings))

	for p := range postings {
		article, err := rep.GetArticleByID(postings[p].String())
		if err == nil {
			articles[p] = *article
		}
	}

	sort.Slice(articles, func(i, j int) bool {
		return articles[i].Docdate > articles[j].Docdate
	})

	bytes, err := json.MarshalIndent(articles, "", "   ")
	if err != nil {
		io.WriteString(w, err.Error())
	} else {
		io.WriteString(w, string(bytes))
	}
}

/*
Returns json encoded string of alternative spelling
*/
func getSpellingSuggestion(w http.ResponseWriter, r *http.Request) {
	logger.Printf("got request %s\n", r.RequestURI)

	parts := strings.Split(r.RequestURI, "/")

	query, err := url.QueryUnescape(parts[len(parts)-1])
	if err != nil {
		io.WriteString(w, err.Error())
	}

	indexname, err := url.QueryUnescape(parts[len(parts)-2])
	if err != nil {
		io.WriteString(w, err.Error())
	}

	if !checkIndexName(indexname) {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	suggestion := indices[indexname].Index.SpellingSuggestion(query, 1)

	bytes, err := json.MarshalIndent(suggestion, "", "   ")
	if err != nil {
		io.WriteString(w, err.Error())
	} else {
		io.WriteString(w, string(bytes))
	}
}

func initIndices() {

	indices = map[string]*ExtendedIndex{}

	//Go through config and add index
	for name, _ := range conf.Indices {
		if conf.Indices[name].Active {
			ei := NewExtendedIndex()
			ei.Name = name
			ei.Language = vocabulary.GetLanguage(conf.Indices[name].Language)
			indices[name] = ei
		}
	}

	added_articles := 0

	//Go through Sources in conf
	//If active, check for article files and
	//add them to indices that match the language
	for _, s := range conf.Sources {
		if !s.Active {
			continue
		}

		folders, err := rep.GetAvailableFolders(s.Id)
		if err != nil {
			logger.Println(err)
			continue
		}

		for _, f := range folders {

			articles, err := rep.GetArticlesByFolderName(f)
			if err != nil {
				logger.Println(err)
				continue
			}

			for a := range articles {
				for ei := range indices {
					lang := vocabulary.GetLanguage(articles[a].Language)
					if lang == indices[ei].Language {
						indices[ei].AddNewsItem(articles[a])
						//indices[ei].Index.AddDocument(index.DocumentFromNewsItem(articles[a]))
						added_articles++
					}
				}
			}
		}
	}
}

func reindex(w http.ResponseWriter, r *http.Request) {
	start_time := time.Now()

	initIndices()

	index_duration := time.Since(start_time)
	logger.Printf("reindexed %d indices in %f seconds\n", len(indices), index_duration.Seconds())

	for i := range indices {
		logger.Printf("index: %s, vocabulary size: %d", indices[i].Name, indices[i].Index.Size())
	}

	io.WriteString(w, http.StatusText(200))
}

/*
Adds Newsitem to repository and NewsQueue
*/
func addArticle(w http.ResponseWriter, r *http.Request) {

	data, err := io.ReadAll(r.Body)
	defer r.Body.Close()
	if err != nil {
		logger.Printf("got request %s\n", r.RequestURI)
		logger.Println("cannot read request body:", err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	ni := feeds.NewsItem{}
	err = json.Unmarshal(data, &ni)
	if err != nil {
		logger.Printf("got request %s\n", r.RequestURI)
		logger.Println("bad request:", err, string(data))
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	err = rep.AddNewsItem(ni, false)
	if err != nil {
		logger.Println("error adding to repository:", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	logger.Printf("got request %s, adding %s\n", r.RequestURI, ni.Id)
	latest_items.Add(&ni)
	w.WriteHeader(http.StatusOK)
}

/*
Runs purge on newsqueue
*/
func purgeNewsQueue(w http.ResponseWriter, r *http.Request) {

	latest_items.Purge()
	w.WriteHeader(http.StatusOK)
}

func getLatest(w http.ResponseWriter, r *http.Request) {
	data, err := json.Marshal(latest_items.Items)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
	io.WriteString(w, string(data))
}

func checkIndexName(name string) bool {
	for i := range indices {
		if name == i {
			return true
		}
	}
	return false
}

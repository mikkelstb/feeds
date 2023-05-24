package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"text/template"

	"github.com/mikkelstb/feeds"
)

var cnf *feeds.Config
var conffile = flag.String("config", "./config.json", "path for conffile")

type webpage struct {
	Title     string
	Form      map[string]string
	Articles  []feeds.NewsItem
	Languages map[string]bool
}

func getTemplates(names ...string) []string {
	templates := make([]string, len(names))
	for n := range names {
		templates[n] = filepath.Join(cnf.Apps["server"]["templates"], names[n])
	}
	return templates
}

func main() {

	var err error
	flag.Parse()

	cnf, err = feeds.ReadConf(*conffile)
	if err != nil {
		log.Fatal(err)
	}

	logfile, err := os.Create(cnf.Apps["server"]["log"])
	if err != nil {
		log.Fatal(err)
	}
	log.SetOutput(logfile)

	http.HandleFunc("/headline", displayLatestPage)
	http.HandleFunc("/headline/search/", displaySearchPage)
	http.HandleFunc("/headline/article/", displayArticlePage)
	http.HandleFunc("/headline/latest/", displayLatestPage)

	http.Handle("/headline/resources/", http.StripPrefix("/headline/resources/", http.FileServer(http.Dir(cnf.Apps["server"]["public_html"]))))
	log.Fatal(http.ListenAndServe(":3300", nil))
}

func displaySearchPage(w http.ResponseWriter, r *http.Request) {

	page := webpage{}
	page.Form = make(map[string]string)
	page.Title = "Test"

	//fmt.Println(r.RequestURI)

	page.Form["query"] = r.URL.Query().Get("searchtext")
	page.Form["selected_language"] = r.URL.Query().Get("language")

	page.Languages = make(map[string]bool)
	for lang := range cnf.Indices {
		if lang == page.Form["selected_language"] {
			page.Languages[lang] = true
		} else {
			page.Languages[lang] = false
		}
	}

	if page.Form["query"] != "" {

		page.Articles = getArticles(page.Form["query"], page.Form["selected_language"])

		// If no articles were found
		// Get spelling suggestion
		if len(page.Articles) == 0 {

			suggestion := getSpellingSuggestion(page.Form["query"])
			if suggestion != "" {
				page.Form["suggestion"] = suggestion
				page.Form["suggestion_link"] = "/headline/search?searchtext=" + suggestion
			}
		}
	}

	t, err := template.ParseFiles(getTemplates("list.html", "header.html", "top.html", "article.html")...)
	if err != nil {
		io.WriteString(w, err.Error())
	}
	err = t.Execute(w, page)
	if err != nil {
		io.WriteString(w, err.Error())
	}
}

func displayLatestPage(w http.ResponseWriter, r *http.Request) {

	page := webpage{}
	page.Title = "Latest"
	page.Articles = getLatest()

	page.Languages = make(map[string]bool)
	for lang := range cnf.Indices {
		if lang == "english" {
			page.Languages[lang] = true
		} else {
			page.Languages[lang] = false
		}
	}

	t, err := template.ParseFiles(getTemplates("list.html", "header.html", "top.html", "article.html")...)
	if err != nil {
		io.WriteString(w, err.Error())
	}
	err = t.Execute(w, page)
	if err != nil {
		fmt.Println(err)
	}
}

func displayArticlePage(w http.ResponseWriter, r *http.Request) {
	fmt.Println(r.RequestURI)

	parts := strings.Split(r.RequestURI, "/")
	query := parts[len(parts)-1]

	ni := getArticle(query)

	page := webpage{}
	if ni != nil {
		page.Title = fmt.Sprintf("article %s", ni.Id)
		page.Articles = append(page.Articles, *ni)
	} else {
		page.Title = "article not found"
	}

	page.Languages = make(map[string]bool)
	for lang := range cnf.Indices {
		if lang == "english" {
			page.Languages[lang] = true
		} else {
			page.Languages[lang] = false
		}
	}

	t, err := template.ParseFiles(getTemplates("single_article.html", "header.html", "top.html", "article_detail.html")...)
	if err != nil {
		io.WriteString(w, err.Error())
	}
	err = t.Execute(w, page)
	if err != nil {
		fmt.Println(err)
	}
}

/*
Returns newsitem based on query id, nil if not found
*/
func getArticle(docid string) *feeds.NewsItem {
	resp, err := http.Get("http://localhost:3333/article/" + docid)
	if err != nil {
		return nil
	}
	defer resp.Body.Close()
	ni := feeds.NewsItem{}

	data, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil
	}

	err = json.Unmarshal(data, &ni)
	if err != nil {
		return nil
	}
	return &ni

}

func getLatest() []feeds.NewsItem {
	resp, err := http.Get("http://localhost:3333/latest/")
	if err != nil {
		return nil
	}
	defer resp.Body.Close()

	articles := make([]feeds.NewsItem, 20)

	data, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil
	}

	err = json.Unmarshal(data, &articles)
	if err != nil {
		return nil
	}
	return articles
}

func getArticles(searchword, language string) []feeds.NewsItem {
	resp, err := http.Get("http://localhost:3333/articles/" + language + "/" + searchword)
	if err != nil {
		return nil
	}
	defer resp.Body.Close()

	articles := make([]feeds.NewsItem, 10)
	data, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil
	}

	err = json.Unmarshal(data, &articles)
	if err != nil {
		return nil
	}
	return articles
}

func getSpellingSuggestion(query string) string {
	resp, err := http.Get("http://localhost:3333/spelling/" + query)
	if err != nil {
		log.Println(err)
		return ""
	}
	defer resp.Body.Close()

	data, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Println(err)
		return ""
	}

	var suggestion string
	err = json.Unmarshal(data, &suggestion)
	if err != nil {
		log.Println(err)
		return ""
	}
	return suggestion
}

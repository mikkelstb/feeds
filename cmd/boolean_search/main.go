package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/mikkelstb/feeds"
	"github.com/mikkelstb/feeds/index"
)

func main() {

	rep, err := feeds.InitRepository("/Users/mikkel/feeds")
	if err != nil {
		log.Fatal(err)
	}

	ii := index.StandardIndex{}
	var docs int

	for s := 1; s < 18; s++ {
		if s == 10 || s == 17 {
			continue
		}
		art, err := rep.GetArticlesByFeed(s, 2023, 3)
		if err != nil {
			log.Fatal(err)
		}
		for a := range art {
			ii.AddDocument(index.DocumentFromNewsItem(art[a]))
			docs++
		}
	}

	fmt.Printf("indexed %d terms from %d articles\n", len(ii), docs)

	for {
		reader := bufio.NewReader(os.Stdin)
		input, err := reader.ReadString('\n')
		if err != nil {
			fmt.Println(err)
		}

		search_terms := strings.Split(input, " ")
		postings := make([]index.PostingsList, len(search_terms))

		for st := range search_terms {
			query := strings.Trim(search_terms[st], " \n")
			postings[st] = ii[query]
		}

		results := index.InterSect(postings...)

		fmt.Printf("found %d articles\n", len(results))

		for i, r := range results {

			if i > 10 {
				break
			}

			ni, _ := rep.GetArticleByID(r.String())
			fmt.Println(ni.Docdate, ni.Source)
			fmt.Println(ni.Headline)
			fmt.Println(ni.Story)
			fmt.Println(ni.Url)
			fmt.Println()
		}

	}
}

package main

import (
	"context"
	"encoding/xml"
	"fmt"
	"html"
	"io"
	"net/http"
	"regexp"
	"strings"
	"time"
	"unicode"
	"unicode/utf8"

	"github.com/microcosm-cc/bluemonday"
	"github.com/mikkelstb/feeds"
)

type RSSFeed struct {
	url              string
	docdate_layout   string
	fetch_time       time.Time
	rss_data         []byte
	story_min_length int

	XMLName     xml.Name `xml:"rss"`
	Title       string   `xml:"channel>title"`
	Description string   `xml:"channel>description"`

	Items []struct {
		Headline string `xml:"title"`
		Story    string `xml:"description"`
		LocalId  string `xml:"guid"`
		Url      string `xml:"link"`
		Docdate  string `xml:"pubDate"`
	} `xml:"channel>item"`

	html_cleaner bluemonday.Policy
}

// Words divided by . ie: bla.bla
var mashed_words_pattern regexp.Regexp = *regexp.MustCompile(`((\w+[\.:,;])(\w+))`)
var endword = *regexp.MustCompile(`\w+…$`)
var spaces = *regexp.MustCompile(`\s{2,}`)
var startbracket = *regexp.MustCompile(`\[.+?\]`)
var newlinepattern = *regexp.MustCompile(`\n+`)
var title_range_table = []*unicode.RangeTable{unicode.Letter, unicode.Punct, unicode.Number}

func NewRSSFeed(conf feeds.FeedConfig) RSSFeed {
	rf := RSSFeed{}
	rf.url = conf.Url
	rf.docdate_layout = conf.DocdateLayout
	rf.story_min_length = conf.StoryMinLength
	rf.fetch_time = time.Now()
	rf.html_cleaner = *bluemonday.StrictPolicy()
	return rf
}

// func NewRSSFeed(url, docdate_layout string) RSSFeed {
// 	rf := RSSFeed{}
// 	rf.url = url
// 	rf.docdate_layout = docdate_layout
// 	rf.html_cleaner = *bluemonday.StrictPolicy()
// 	rf.fetch_time = time.Now()

// 	return rf
// }

/*
	Connects to http server defined by url
	Function returns an error if either connection or reading of data fails
	On success the response will be stored into feed.rss_data
*/

func (feed *RSSFeed) Read() error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	request, err := http.NewRequestWithContext(ctx, http.MethodGet, feed.url, nil)
	if err != nil {
		return err
	}

	client := &http.Client{}
	response, err := client.Do(request)
	if err != nil {
		return err
	}
	defer response.Body.Close()

	if response.StatusCode != 200 {
		return fmt.Errorf("error response from server: %s", response.Status)
	}

	feed.rss_data, err = io.ReadAll(response.Body)
	if err != nil {
		return err
	}
	return nil
}

/*
	Parse reads the rss_data, and uses xml.Unmarshal to set itself up
*/

func (feed *RSSFeed) Parse() error {
	err := xml.Unmarshal(feed.rss_data, feed)
	return err
}

func (feed *RSSFeed) HasNext() bool {
	return len(feed.Items) != 0
}

/*
	GetNext returns the first newsitem in the Rssfeed
	It runs sanitize on headline and story, inserts fetchtime from feed
	If storytext is smaller than 16 bytes returns nil with errormessage
	Rerurns nil, error if feed.Items is empty
*/

func (feed *RSSFeed) GetNext() (*feeds.NewsItem, error) {

	if len(feed.Items) == 0 {
		panic(fmt.Errorf("warning: no more items in feed"))
	}

	nextitem := feed.Items[len(feed.Items)-1]

	// Slice off current item
	feed.Items = feed.Items[0 : len(feed.Items)-1]

	n := new(feeds.NewsItem)

	n.Headline = feed.sanitize(nextitem.Headline)
	n.Headline = trimEmoji(n.Headline)

	n.Story = feed.sanitize(nextitem.Story)
	n.Url = nextitem.Url
	n.LocalId = nextitem.LocalId

	dd, err := time.Parse(feed.docdate_layout, nextitem.Docdate)
	if err != nil {
		return nil, err
	}
	n.Docdate = dd.UTC().Format(time.RFC3339)
	n.FetchTime = feed.fetch_time.UTC().Format(time.RFC3339)

	if utf8.RuneCountInString(n.Story) < feed.story_min_length {
		return nil, fmt.Errorf("story length (%d) smaller than limit for feed (%d)", utf8.RuneCountInString(n.Story), feed.story_min_length)
	}

	return n, nil
}

/*
	private function to clean up most unnessesary symbols and html tags
*/

func (feed *RSSFeed) sanitize(field string) string {

	field = newlinepattern.ReplaceAllString(field, " ")
	field = feed.html_cleaner.Sanitize(field)
	field = html.UnescapeString(field)
	field = startbracket.ReplaceAllString(field, "")
	field = endword.ReplaceAllString(field, "")
	field = spaces.ReplaceAllString(field, " ")
	field = mashed_words_pattern.ReplaceAllString(field, "$2 $3")
	field = strings.Trim(field, " :,-.")

	return field
}

func trimEmoji(s string) string {
	return strings.TrimFunc(s, func(r rune) bool {
		return !unicode.IsOneOf(title_range_table, r)
		//return !unicode.IsLetter(r) && !unicode.IsPunct(r)
	})
}

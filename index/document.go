package index

import (
	"unicode/utf8"

	"github.com/mikkelstb/feeds"
	"github.com/mikkelstb/feeds/vocabulary"
	"golang.org/x/exp/slices"
)

/*
Document represents a text as a bag of words
Unlike NewsItem the text in Document is not meaningful for reading
Document has a hashmap of all words from a NewsItem and their frequency
*/
type Document struct {
	ID       DocID
	terms    map[string]int
	language vocabulary.Language
}

type DocID [12]rune
type PostingsList []DocID

func (di DocID) String() string {
	return string(di[:])
}

func (d DocID) MarshalJSON() ([]byte, error) {
	return []byte(`"` + d.String() + `"`), nil
}

func NewDocument(docid string) Document {
	d := Document{}
	d.ID = DocID([]rune(docid))
	d.terms = make(map[string]int)

	return d
}

// Convert NewsItem into Document
// Get Id from NewsItem
// Take all words from story and headline and add them to hash
func DocumentFromNewsItem(ni feeds.NewsItem) Document {
	d := Document{}
	d.ID = DocID([]rune(ni.Id[:12]))
	d.terms = make(map[string]int)

	d.language = vocabulary.GetLanguage(ni.Language)
	d.AddTerms(vocabulary.TokenizeText(ni.Headline, ni.Story))
	return d
}

func (d *Document) SetLanguage(lang string) {
	d.language = vocabulary.GetLanguage(lang)
}

// Adds terms to Document hash
func (d *Document) AddTerms(words []string) {

	for w := range words {
		term := vocabulary.CleanWord(words[w], d.language)

		if slices.Contains(vocabulary.SkipLists[d.language], term) {
			continue
		}

		if utf8.RuneCountInString(term) < vocabulary.MinWordLength[d.language] {
			continue
		}

		d.terms[term]++
	}
}

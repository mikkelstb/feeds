package index

import (
	"strings"

	"github.com/mikkelstb/feeds"
	"github.com/mikkelstb/feeds/vocabulary"
	"golang.org/x/exp/slices"
)

/*
Document has a hashmap of all words from a NewsItem and their frequency
*/
type Document struct {
	ID       DocID
	terms    map[string]int
	language vocabulary.Language
}

// type Term struct {
// 	Frequency int
// 	Pos       []int
// }

type DocID [16]rune

// func (tid DocID) String() string {
// 	sb := strings.Builder{}
// 	for x := range tid {
// 		sb.WriteRune(tid[x])
// 	}
// 	return sb.String()
// }

func (tid DocID) String() string {
	return string(tid[:])
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

func DocumentFromNewsItem(ni feeds.NewsItem) Document {
	d := Document{}
	d.ID = DocID([]rune(ni.Id[:16]))
	d.terms = make(map[string]int)

	d.language = vocabulary.GetLanguage(ni.Language)
	d.AddTerms(vocabulary.TokenizeText(ni.Headline, ni.Story))
	return d
}

func (d *Document) SetLanguage(lang string) {
	d.language = vocabulary.GetLanguage(lang)
}

func (d *Document) AddTerms(words []string) {

	for w := range words {
		term := vocabulary.CleanWord(words[w], d.language)

		if len([]rune(term)) < 3 {
			continue
		}
		if slices.Contains(vocabulary.SkipLists[d.language], term) {
			continue
		}
		d.terms[term]++
	}
}

type PostingsList []DocID

/*
StandardIndex stores all searchable terms as well as the docIDs
*/
type StandardIndex map[string]PostingsList

/*
adding a document to the InvertedIndex
terms already present will be updated with docid
*/
func (ii StandardIndex) AddDocument(d Document) {
	for t := range d.terms {
		ii[t] = append(ii[t], d.ID)
	}
}

func (ii StandardIndex) GetPostings(query string) PostingsList {
	query_terms := strings.Split(query, " ")
	if len(query_terms) == 0 {
		return nil
	}
	if len(query_terms) == 1 {
		return ii[query_terms[0]]
	} else {
		postings := make([]PostingsList, len(query_terms))
		for qt := range query_terms {
			postings[qt] = ii[query_terms[qt]]
		}
		return InterSect(postings...)
	}
}

func (ii StandardIndex) String() string {
	sb := strings.Builder{}
	for t := range ii {
		sb.WriteString(t)
		sb.WriteString(": ")
		for _, r := range ii[t] {
			sb.WriteString(r.String())
		}
		sb.WriteString("\n")
	}
	return sb.String()
}

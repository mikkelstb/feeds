package index

import (
	"strings"

	"github.com/mikkelstb/feeds/vocabulary"
)

/*
StandardIndex stores all searchable terms as well as the docIDs
*/
type StandardIndex struct {
	index map[string]PostingsList
}

func NewStandardIndex() StandardIndex {
	si := StandardIndex{}
	si.index = make(map[string]PostingsList)
	return si
}

/*
adding a document to the InvertedIndex
terms already present will be updated with docid
*/
func (ii *StandardIndex) AddDocument(d Document) {
	for t := range d.terms {
		ii.index[t] = append(ii.index[t], d.ID)
	}
}

/*
GetPostings returns a list of postings(doc-ids) based on one or more
search terms (divided by a space)
*/
func (ii StandardIndex) GetPostings(query string) PostingsList {
	query_terms := strings.Split(query, " ")
	if len(query_terms) == 0 {
		return nil
	}
	if len(query_terms) == 1 {
		return ii.index[query_terms[0]]
	} else {
		postings := make([]PostingsList, len(query_terms))
		for qt := range query_terms {
			postings[qt] = ii.index[query_terms[qt]]
		}
		return InterSect(postings...)
	}
}

func (ii StandardIndex) Size() int {
	return len(ii.index)
}

func (ii *StandardIndex) AddTerm(t string) {
	ii.index[t] = nil
}

func (ii StandardIndex) String() string {
	sb := strings.Builder{}
	for t := range ii.index {
		sb.WriteString(t)
		sb.WriteString(": ")
		for _, r := range ii.index[t] {
			sb.WriteString(r.String())
		}
		sb.WriteString("\n")
	}
	return sb.String()
}

/*
Returns first string of key with EditDistance=min_dist to query
If none exist return empty string
*/
func (ii StandardIndex) SpellingSuggestion(query string, min_dist int) string {
	for key := range ii.index {
		if vocabulary.EditDistance([]rune(query), []rune(key)) <= min_dist {
			return key
		}
	}
	return ""
}

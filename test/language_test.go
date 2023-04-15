package feeds_test

import (
	"fmt"
	"strings"
	"testing"

	"github.com/mikkelstb/feeds/vocabulary"
)

type pair struct {
	test   string
	answer string
}

func TestExtractTerms(t *testing.T) {
	text := string("EU har sendt mer enn dobbelt så mye penger til Russland:")

	terms := vocabulary.TokenizeText(text)
	fmt.Println(strings.Join(terms, ","))
}

func TestTrimWord(t *testing.T) {

	words := []pair{
		{test: ",test", answer: "test"},
		{test: ",test+-", answer: "test"},
		{test: "test:", answer: "test"},
		{test: "test---", answer: "test"},
		{test: "bla", answer: "bla"},
		{test: "", answer: ""},
		{test: "i", answer: "i"},
	}

	for w := range words {
		if vocabulary.TrimWord([]rune(words[w].test)) != words[w].answer {
			t.Error("word not trimmed correctly")
		}
	}
}

func TestEditDistance(t *testing.T) {

	if r := vocabulary.EditDistance([]rune("fast"), []rune("cats")); r != 3 {
		t.Errorf("editdistance was %d should have been 3", r)
	}

	if r := vocabulary.EditDistance([]rune("brittney"), []rune("britney")); r != 1 {
		t.Errorf("editdistance was %d should have been 1", r)
	}

	if r := vocabulary.EditDistance([]rune("apple"), []rune("aplpe")); r != 2 {
		t.Errorf("editdistance was %d should have been 2", r)
	}

	if r := vocabulary.EditDistance([]rune("dr"), []rune("politics")); r != 8 {
		t.Errorf("editdistance was %d should have been 8", r)
	}
}

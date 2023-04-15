package feeds_test

import (
	"testing"

	"github.com/mikkelstb/feeds/index"
)

func TestAddTerm(t *testing.T) {

	doc1 := index.NewDocument("04230228568ac35d")
	doc2 := index.NewDocument("04230228568aca8f")

	doc1.SetLanguage("da")
	doc2.SetLanguage("da")

	doc1.AddTerms([]string{"med", "hvis", "af", "kun", "imens"})
	doc2.AddTerms([]string{"hvor", "mig", "af", "bare", "imens", "hvor"})

	var ti = index.StandardIndex{}

	ti.AddDocument(doc1)
	ti.AddDocument(doc2)

	//fmt.Println(ti)

	if len(ti) != 6 {
		t.Errorf("term index is %d should be 6", len(ti))
	}

	//fmt.Println(ti)
}

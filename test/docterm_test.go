package feeds_test

import (
	"fmt"
	"testing"

	"github.com/mikkelstb/feeds/index"
)

func TestAddTerm(t *testing.T) {

	doc1 := index.NewDocument("042302288ac3")
	doc2 := index.NewDocument("042302288aca")

	doc1.SetLanguage("da")
	doc2.SetLanguage("da")

	doc1.AddTerms([]string{"med", "hvis", "af", "kun", "imens"})
	doc2.AddTerms([]string{"hvor", "mig", "af", "bare", "imens", "hvor"})

	var ti = index.NewStandardIndex()

	ti.AddDocument(doc1)
	ti.AddDocument(doc2)

	//fmt.Println(ti)

	if ti.Size() != 8 {
		t.Errorf("term index is %d should be 6", ti.Size())
		fmt.Println(ti.String())
	}

	//fmt.Println(ti)
}

func TestSpellingSuggestion(t *testing.T) {
	var ti = index.NewStandardIndex()
	ti.AddTerm("oslo")
	ti.AddTerm("køber")
	ti.AddTerm("københavn")
	ti.AddTerm("køben")

	result := ti.SpellingSuggestion("kobenhavn", 1)
	if result != "københavn" {
		t.Errorf("expected københavn got %s", result)
	}

	println(ti.SpellingSuggestion("slo", 1))
	println(ti.SpellingSuggestion("kober", 1))

}

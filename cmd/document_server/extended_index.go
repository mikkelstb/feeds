package main

import (
	"github.com/mikkelstb/feeds"
	"github.com/mikkelstb/feeds/index"
	"github.com/mikkelstb/feeds/vocabulary"
)

type ExtendedIndex struct {
	Index    index.StandardIndex
	Name     string
	Language vocabulary.Language
}

func NewExtendedIndex() *ExtendedIndex {
	ei := new(ExtendedIndex)
	ei.Index = index.NewStandardIndex()
	return ei
}

func (ei *ExtendedIndex) AddNewsItem(ni feeds.NewsItem) {
	doc := index.DocumentFromNewsItem(ni)
	ei.Index.AddDocument(doc)
}

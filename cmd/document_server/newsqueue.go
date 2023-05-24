package main

import (
	"sort"

	"github.com/mikkelstb/feeds"
)

type NewsQueue struct {
	Items []*feeds.NewsItem
	size  int
}

func NewNewsQueue(size int) NewsQueue {
	nq := NewsQueue{}
	nq.Items = []*feeds.NewsItem{}
	nq.size = size
	return nq
}

func (nq NewsQueue) String() string {
	var result string
	for ni := range nq.Items {
		if nq.Items[ni] != nil {
			result += (nq.Items[ni].Docdate)
			result += "\n"
		}
	}
	return result
}

// Insert ni into Items
func (nq *NewsQueue) Add(ni *feeds.NewsItem) {

	for n := range nq.Items {
		if nq.Items[n].Id == ni.Id {
			return
		}
	}

	nq.Items = append(nq.Items, ni)
}

/*
Sorts all newsitems descending, and cuts off
*/
func (nq *NewsQueue) Purge() {
	sort.Slice(nq.Items, func(i, j int) bool {
		return nq.Items[i].GetDocdate().After(nq.Items[j].GetDocdate())
	})
	if nq.size <= len(nq.Items) {
		nq.Items = nq.Items[:nq.size]
	}
}

func (nq NewsQueue) Search(query string) *feeds.NewsItem {
	for ni := range nq.Items {
		if nq.Items[ni].Id == query {
			return nq.Items[ni]
		}
	}
	return nil
}

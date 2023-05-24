package feeds_test

import (
	"os"
	"testing"

	"github.com/mikkelstb/feeds"
)

func TestInitRepository(t *testing.T) {

	// Create a tempoary folder and check if it is valid for a repository
	// Delete folder afterwards

	path := "./temp_repository"
	//t.Log("Using folder", path)

	err := os.MkdirAll(path, 0777)
	if err != nil {
		t.Error(err)
	}

	_, err = feeds.InitRepository(path)
	if err != nil {
		t.Error(err)
	}

}

func TestAddNewsItem(t *testing.T) {

	path := "./temp_repository"
	err := os.MkdirAll(path, 0777)
	if err != nil {
		t.Error(err)
	}

	r, err := feeds.InitRepository(path)
	if err != nil {
		t.Error(err)
	}

	ni := new(feeds.NewsItem)
	ni.Headline = "Test"
	ni.Story = "Testing story"
	ni.Docdate = "2023-04-05T18:00:00Z"
	ni.Country = "Da"
	ni.SourceName = "DR Nyheder"
	ni.SourceId = 5
	ni.Mediatype = "web"
	ni.LocalId = "0123456789"

	err = r.AddNewsItem(*ni, true)
	if err != nil {
		t.Error(err)
	}

	err = os.RemoveAll(path)
	if err != nil {
		t.Error(err)
	}
}

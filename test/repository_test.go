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
	ni.Feed = "DR Nyheder"
	ni.FeedId = 5
	ni.Mediatype = "web"
	ni.LocalId = "0123456789"

	err = r.AddNewsItem(*ni)
	if err != nil {
		t.Error(err)
	}

	err = os.RemoveAll(path)
	if err != nil {
		t.Error(err)
	}
}

func TestGetArticlesByFeed(t *testing.T) {
	path := "/Users/mikkel/feeds"

	rep, err := feeds.InitRepository(path)
	if err != nil {
		t.Error(err)
	}

	_, err = rep.GetArticlesByFeed(2, 2022, 9)
	if err != nil {
		t.Error(err)
	}

}

func TestGetArticleByID(t *testing.T) {
	path := "/Users/mikkel/feeds"

	rep, err := feeds.InitRepository(path)
	if err != nil {
		t.Error(err)
	}

	_, err = rep.GetArticleByID("04230220f9531822")
	if err != nil {
		t.Error(err)
	}
	_, err = rep.GetArticleByID("04230220f9531823")
	if err == nil {
		t.Error(err)
	}

	_, err = rep.GetArticleByID("32r23few")
	if err == nil {
		t.Error(err)
	}
}

// func generateArticle(feedid int, date string) feeds.NewsItem {

// 	ni := feeds.NewsItem{}
// 	ni.Headline = "Test"
// 	ni.Story = "Testing story"
// 	ni.Docdate = "2023-04-05T18:00:00Z"
// 	ni.Country = "Da"
// 	ni.Feed = "DR Nyheder"
// 	ni.FeedId = 5
// 	ni.Mediatype = "web"
// 	ni.LocalId = "0123456789"

// 	return ni
// }

package feeds

import (
	"crypto/md5"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io"
	"time"
)

type NewsItem struct {
	SourceName string              `json:"sourceName"`
	SourceId   int                 `json:"sourceID"`
	FeedName   string              `json:"feedName"`
	Mediatype  string              `json:"mediatype"`
	Headline   string              `json:"headline"`
	Story      string              `json:"story"`
	Url        string              `json:"url"`
	Language   string              `json:"language"`
	Country    string              `json:"country"`
	Docdate    string              `json:"docdate"`
	FetchTime  string              `json:"fetchTime"`
	Id         string              `json:"id"`
	LocalId    string              `json:"localId"`
	Categories map[string][]string `json:"category"`
}

// Helper struct in order to serialize NewsItem to JSON
type IdNewsItem NewsItem

/*
GetId() returns a "unique" four letter id based on local id.
If local id is not present, use headline
*/
func (ni NewsItem) GetId() string {
	id := md5.New()
	if ni.LocalId != "" {
		io.WriteString(id, ni.LocalId)
	} else {
		io.WriteString(id, ni.Headline)
	}
	return fmt.Sprintf("%02d%v%v", ni.SourceId, ni.GetDocdate().Format("060102"), hex.EncodeToString(id.Sum(nil))[0:4])
}

func (ni NewsItem) GetDocdate() time.Time {
	dd, _ := time.Parse(time.RFC3339, ni.Docdate)
	return dd
}

func (ni *NewsItem) SetDocdate(t time.Time) {
	ni.Docdate = t.Format(time.RFC3339)
}

func (ni NewsItem) MarshalJSON() ([]byte, error) {
	return json.Marshal(struct {
		Id string `json:"id"`
		IdNewsItem
	}{
		Id:         ni.GetId(),
		IdNewsItem: IdNewsItem(ni),
	})
}

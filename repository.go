package feeds

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"regexp"
)

var filenameregex = regexp.MustCompile(`^(\d{2})(\d{4})\d{2}[0-9a-f]{8}(.json)?$`)

type Repository struct {
	root string
}

/*
Checks if path is valid and a folder
stores the path as root of directory
*/
func InitRepository(path string) (*Repository, error) {

	r := new(Repository)
	r.root = path

	_, err := os.Open(r.root)
	if err != nil {
		return nil, err
	}
	return r, nil
}

/*
Newsitem shall be stored in the following manner:
folder: /sourceid(2 digits)/year-month(6 digits)/
filename: sourceid(2 digits) + (yymmdd) + hexid (8chars 0-9a-f) .json
*/
func (r *Repository) AddNewsItem(ni NewsItem) error {

	path := filepath.Join(
		r.root,
		fmt.Sprintf("%02d", ni.FeedId),
		fmt.Sprint(ni.GetDocdate().Year())+fmt.Sprintf("%02d", ni.GetDocdate().Month()),
	)

	err := os.MkdirAll(path, 0777)
	if err != nil {
		return err
	}

	file, err := os.Create(filepath.Join(path, ni.GetId()+".json"))
	if err != nil {
		return err
	}

	data, err := ni.ToJson()
	if err != nil {
		return err
	}

	_, err = file.Write(data)
	if err != nil {
		return err
	}

	return nil
}

/*
Returns a slice of NewsItem based on feed id, month and year
*/
func (r Repository) GetArticlesByFeed(feed_id, year, month int) ([]NewsItem, error) {

	queryfolder := filepath.Join(
		r.root,
		fmt.Sprintf("%02d", feed_id),
		fmt.Sprint(year)+fmt.Sprintf("%02d", month),
	)

	files, err := os.ReadDir(queryfolder)
	if err != nil {
		return nil, err
	}

	newsitems := make([]NewsItem, len(files))
	for file := range files {
		info, _ := files[file].Info()
		if filenameregex.MatchString(info.Name()) {
			data, err := os.ReadFile(queryfolder + "/" + files[file].Name())
			if err != nil {
				return nil, err
			}
			newsit := NewsItem{}
			err = json.Unmarshal(data, &newsit)
			if err != nil {
				return nil, err
			}
			newsitems[file] = newsit
		}
	}
	return newsitems, nil
}

/*
Returns a slice of NewsItem based on feed id, month and year
*/
func (r Repository) GetArticleByID(article_id string) (*NewsItem, error) {

	match := filenameregex.FindStringSubmatch(article_id)
	if match == nil {
		return nil, fmt.Errorf("article id invalid: %s", article_id)
	}

	fp := filepath.Join(
		r.root,
		match[1],
		"20"+match[2],
		match[0]+".json",
	)

	file, err := os.ReadFile(fp)
	if err != nil {
		return nil, err
	}
	ni := NewsItem{}
	json.Unmarshal(file, &ni)

	return &ni, nil
}

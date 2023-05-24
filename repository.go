package feeds

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"regexp"
)

var filenameregex = regexp.MustCompile(`^(\d{2})(\d{4})\d{2}[0-9a-f]{4}(.json)?$`)
var monthfolder = regexp.MustCompile(`^\d{4}\d{2}$`)

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
filename: sourceid(2 digits) + (yymmdd) + hexid (4 chars 0-9a-f) .json
*/
func (r *Repository) AddNewsItem(ni NewsItem, overwrite bool) error {

	path := filepath.Join(
		r.root,
		fmt.Sprintf("%02d", ni.SourceId),
		fmt.Sprint(ni.GetDocdate().Year())+fmt.Sprintf("%02d", ni.GetDocdate().Month()),
	)

	err := os.MkdirAll(path, 0777)
	if err != nil {
		return err
	}

	filename := filepath.Join(path, ni.GetId()+".json")

	if !overwrite {
		exist, _ := exists(filename)
		if exist {
			return fmt.Errorf("id %s already exists", ni.Id)
		}
	}

	file, err := os.Create(filename)
	if err != nil {
		return err
	}

	data, err := json.MarshalIndent(&ni, "", " ")
	if err != nil {
		return err
	}

	_, err = file.Write(data)
	if err != nil {
		return err
	}

	return nil
}

func (r Repository) GetAvailableFolders(source_id int) ([]string, error) {
	folders := make([]string, 0)

	sourcefolder := filepath.Join(
		r.root,
		fmt.Sprintf("%02d", source_id),
	)

	items, err := os.ReadDir(sourcefolder)
	if err != nil {
		return nil, err
	}
	for i := range items {
		if items[i].IsDir() && monthfolder.Match([]byte(items[i].Name())) {
			folders = append(folders, filepath.Join(r.root, fmt.Sprintf("%02d", source_id), items[i].Name()))
		}
	}
	return folders, nil
}

/*
Returns a slice of NewsItems based on feed id, month and year
*/
func (r Repository) GetArticlesByFeed(source_id, year, month int) ([]NewsItem, error) {

	foldername := filepath.Join(
		r.root,
		fmt.Sprintf("%02d", source_id),
		fmt.Sprint(year)+fmt.Sprintf("%02d", month),
	)

	return r.GetArticlesByFolderName(foldername)
}

/*
Returns a slice of NewsItems based on feed id, month and year
*/
func (r Repository) GetArticlesByFolderName(folder string) ([]NewsItem, error) {

	files, err := os.ReadDir(folder)
	if err != nil {
		return nil, err
	}

	newsitems := make([]NewsItem, len(files))
	for file := range files {
		info, _ := files[file].Info()
		if filenameregex.MatchString(info.Name()) {
			data, err := os.ReadFile(folder + "/" + files[file].Name())
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

/*
Return true, error=nil if file exists
*/
func exists(path string) (bool, error) {
	_, err := os.Stat(path)
	if err == nil {
		return true, nil
	}
	if os.IsNotExist(err) {
		return false, nil
	}
	return false, err
}

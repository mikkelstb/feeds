package feeds

import (
	"encoding/json"
	"os"
)

type Config struct {
	Apps         map[string]map[string]string `json:"apps"`
	Repositories map[string]RepositoryConfig  `json:"repositories"`
	Sources      []SourceConfig               `json:"sources"`
	Indices      map[string]IndexConfig       `json:"indices"`
}

type SourceConfig struct {
	Id        int                   `json:"id"`
	Name      string                `json:"name"`
	Country   string                `json:"country"`
	Language  string                `json:"language"`
	Mediatype string                `json:"mediatype"`
	Active    bool                  `json:"active"`
	Feeds     map[string]FeedConfig `json:"feeds"`
}

type FeedConfig struct {
	Active         bool   `json:"active"`
	Name           string `json:"name"`
	Datatype       string `json:"datatype"`
	Url            string `json:"url"`
	DocdateLayout  string `json:"docdate_layout"`
	StoryMinLength int    `json:"story_min_length"`
}

type RepositoryConfig struct {
	Type           string `json:"type"`
	Address        string `json:"address"`
	Active         bool   `json:"active"`
	EraseAfterDays int    `json:"erase_after_days"`
}

type IndexConfig struct {
	Active   bool   `json:"active"`
	Language string `json:"language"`
}

func ReadConf(filename string) (*Config, error) {
	cfg := new(Config)
	cfg_file, err := os.ReadFile(filename)
	if err != nil {
		return nil, err
	}

	err = json.Unmarshal(cfg_file, cfg)
	if err != nil {
		return nil, err
	}
	return cfg, nil
}

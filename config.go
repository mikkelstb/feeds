package feeds

import (
	"encoding/json"
	"os"
)

type Config struct {
	Logs         map[string]string           `json:"logs"`
	Repositories map[string]RepositoryConfig `json:"repositories"`
	Sources      []SourceConfig              `json:"sources"`
}

type SourceConfig struct {
	Active      bool              `json:"active"`
	Id          int               `json:"id"`
	Name        string            `json:"name"`
	Screen_name string            `json:"screen_name"`
	Country     string            `json:"country"`
	Language    string            `json:"language"`
	Mediatype   string            `json:"mediatype"`
	Feed        map[string]string `json:"feed"`
}

type RepositoryConfig struct {
	Type           string `json:"type"`
	Address        string `json:"address"`
	Active         bool   `json:"active"`
	EraseAfterDays int    `json:"erase_after_days"`
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

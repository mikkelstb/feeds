package main

import "github.com/mikkelstb/feeds"

type Source struct {
	config feeds.SourceConfig
	feeds  map[string]*RSSFeed
}

func (s Source) Name() string {
	return s.config.Name
}

/*
returns source struct with all active feeds initiated
*/
func NewSource(cfg feeds.SourceConfig) Source {
	s := Source{config: cfg}
	s.feeds = make(map[string]*RSSFeed, 0)

	for feedname := range cfg.Feeds {
		if !cfg.Feeds[feedname].Active {
			continue
		}
		feed := NewRSSFeed(cfg.Feeds[feedname])
		s.feeds[feedname] = &feed
	}
	return s
}

func (s *Source) Process() error {

	for _, f := range s.feeds {

		err := f.Read()
		if err != nil {
			return err
		}

		err = f.Parse()
		if err != nil {
			return err
		}
	}
	return nil
}

func (s *Source) GetNewsitems() ([]feeds.NewsItem, []error) {

	var articles []feeds.NewsItem
	var errs []error

	for feedname, feed := range s.feeds {

		for feed.HasNext() {

			article, err := feed.GetNext()
			if err != nil {
				errs = append(errs, err)
				continue
			}

			article.SourceId = s.config.Id
			article.SourceName = s.config.Name
			article.FeedName = feedname
			article.Mediatype = s.config.Mediatype
			article.Country = s.config.Country
			article.Language = s.config.Language

			articles = append(articles, *article)
		}
	}

	return articles, errs
}

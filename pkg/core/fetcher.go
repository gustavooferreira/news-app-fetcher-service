package core

import (
	"context"
	"fmt"
	"time"

	"github.com/gustavooferreira/news-app-fetcher-service/pkg/clients/artmgmt"
	"github.com/gustavooferreira/news-app-fetcher-service/pkg/clients/feedmgmt"
	"github.com/gustavooferreira/news-app-fetcher-service/pkg/core/log"
	"github.com/mmcdole/gofeed"
)

type Fetcher struct {
	Logger         log.Logger
	FeedClient     FeedsClient
	ArticlesClient ArticlesClient
	RSSParser      *gofeed.Parser

	WaitPeriod time.Duration
	quitChan   chan chan struct{}
}

func NewFetcher(logger log.Logger, feedsC FeedsClient, articlesC ArticlesClient, waitPeriod int) *Fetcher {
	quitChan := make(chan chan struct{}, 1)
	fetcher := &Fetcher{
		Logger:         logger,
		FeedClient:     feedsC,
		ArticlesClient: articlesC,
		RSSParser:      gofeed.NewParser(),
		WaitPeriod:     time.Duration(waitPeriod) * time.Second,
		quitChan:       quitChan,
	}

	return fetcher
}

func (f *Fetcher) ShutDown(ctx context.Context) error {
	quitResponseChan := make(chan struct{}, 1)
	f.quitChan <- quitResponseChan

	select {
	case <-ctx.Done():
		return ctx.Err()
	case <-quitResponseChan:
		return nil
	}
}

func (f *Fetcher) Run() error {
	quitFlag := false
	var quitResponseChan chan struct{}

	for {
		select {
		case quitResponseChan = <-f.quitChan:
			quitFlag = true
		case <-time.After(f.WaitPeriod):
			f.WorkerRun()
		}

		if quitFlag {
			break
		}
	}

	quitResponseChan <- struct{}{}
	return nil
}

func (f *Fetcher) WorkerRun() {
	f.Logger.Info("new fetching cycle starting")
	start := time.Now()

	// Get list of RSS feed URLs
	feeds, err := f.FeedClient.GetFeeds("", "", true)
	if err != nil {
		f.Logger.Error(fmt.Sprintf("error while sending request to feeds mgmt service: %s", err.Error()))
		return
	}

	// We'll fetch each feed one at a time, but we could have done this in parallel
	for _, feed := range feeds {
		articles, err := f.FetchUpdates(feed)

		err = f.ArticlesClient.AddArticles(articles)
		if err != nil {
			f.Logger.Error(fmt.Sprintf("error while sending request to articles mgmt service: %s", err.Error()))
			return
		}
	}

	end := time.Now()
	fetchDuration := end.Sub(start)
	f.Logger.Info("fetching cycle ending", log.Field("duration", fetchDuration.String()))
}

func (f *Fetcher) FetchUpdates(feed feedmgmt.Feed) (articles artmgmt.Articles, err error) {
	feedResult, err := f.RSSParser.ParseURL(feed.URL)
	if err != nil {
		return articles, err
	}

	if feedResult.Items == nil || len(feedResult.Items) == 0 {
		return artmgmt.Articles{}, nil
	}

	articles = make(artmgmt.Articles, 0, len(feedResult.Items))

	for _, item := range feedResult.Items {
		article := artmgmt.Article{
			GUID:        item.GUID,
			Title:       item.Title,
			Description: item.Description,
			Link:        item.Link,
			Provider:    feed.Provider,
			Category:    feed.Category,
		}

		if item.PublishedParsed != nil {
			article.PublishedTime = *item.PublishedParsed
		}

		articles = append(articles, article)
	}

	return articles, nil
}

package core

import (
	"context"
	"time"

	"github.com/gustavooferreira/news-app-fetcher-service/pkg/core/log"
)

type Fetcher struct {
	Logger         log.Logger
	FeedClient     FeedsClient
	ArticlesClient ArticlesClient

	WaitPeriod time.Duration
	quitChan   chan chan struct{}
}

func NewFetcher(logger log.Logger, feedsC FeedsClient, articlesC ArticlesClient, waitPeriod int) *Fetcher {
	quitChan := make(chan chan struct{}, 1)
	fetcher := &Fetcher{
		Logger:         logger,
		FeedClient:     feedsC,
		ArticlesClient: articlesC,
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

	// Fetch all of them in parallel

	// Send all of the results to articles mgmt API

	end := time.Now()
	fetchDuration := end.Sub(start)
	f.Logger.Info("fetching cycle ending", log.Field("duration", fetchDuration.String()))
}

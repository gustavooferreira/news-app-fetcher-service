package core

import (
	"context"

	"github.com/gustavooferreira/news-app-fetcher-service/pkg/clients/artmgmt"
	"github.com/gustavooferreira/news-app-fetcher-service/pkg/clients/feedmgmt"
)

// ShutDowner represents anything that can be shutdown.
type ShutDowner interface {
	ShutDown(ctx context.Context) error
}

// FeedsClient represents the feeds management service
type FeedsClient interface {
	GetFeeds(provider string, category string, enabled bool) (feeds feedmgmt.Feeds, err error)
}

// ArticlesClient represents the articles management service
type ArticlesClient interface {
	AddArticles(articles artmgmt.Articles) (err error)
}

package core_test

import (
	"testing"

	"github.com/gustavooferreira/news-app-fetcher-service/mocks"
	"github.com/gustavooferreira/news-app-fetcher-service/pkg/clients/feedmgmt"
	"github.com/gustavooferreira/news-app-fetcher-service/pkg/core"
	"github.com/gustavooferreira/news-app-fetcher-service/pkg/core/log"
	"github.com/stretchr/testify/mock"
)

func TestWorkerRun(t *testing.T) {

	logger := log.NullLogger{}

	feedsClient, articlesClient := setupMockClients()

	fetcher := core.NewFetcher(logger, feedsClient, articlesClient, 10)

	fetcher.WorkerRun()

	feedsClient.AssertCalled(t, "GetFeeds", "", "", true)

	// Mock gofeed Parser

	// Assert
}

func setupMockClients() (*mocks.FeedsClient, *mocks.ArticlesClient) {
	mockFeedsClient := &mocks.FeedsClient{}
	mockArticlesClient := &mocks.ArticlesClient{}

	call := mockFeedsClient.On("GetFeeds", "", "", true)
	call = call.Return(feedmgmt.Feeds{feedmgmt.Feed{URL: "url1", Provider: "provider1", Category: "category1"}}, nil)

	call = mockArticlesClient.On("AddArticles", mock.Anything)
	call = call.Return(nil)

	return mockFeedsClient, mockArticlesClient
}

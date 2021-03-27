package main

import (
	"fmt"
	"net/http"
	"os"
	"time"

	"github.com/gustavooferreira/news-app-fetcher-service/pkg/clients/artmgmt"
	"github.com/gustavooferreira/news-app-fetcher-service/pkg/clients/feedmgmt"
	"github.com/gustavooferreira/news-app-fetcher-service/pkg/core"
	"github.com/gustavooferreira/news-app-fetcher-service/pkg/core/lifecycle"
	"github.com/gustavooferreira/news-app-fetcher-service/pkg/core/log"
)

func main() {
	retCode := mainLogic()
	os.Exit(retCode)
}

func mainLogic() int {
	// Setup logger
	logger := core.NewAppLogger(os.Stdout, log.INFO)
	defer logger.Sync()

	logger.Info("APP starting")

	// Read config
	logger.Info("reading configuration", log.Field("type", "setup"))
	config := core.NewConfig()
	if err := config.LoadConfig(); err != nil {
		logger.Error(err.Error(), log.Field("type", "config"))
		return 1
	}

	// TODO: Set log level after reading config
	// something like this:
	// logger.SetLevel(config.Options.LogLevel)

	httpClient := &http.Client{
		Timeout: time.Second * time.Duration(config.Options.HTTPClientTimeout),
	}

	feedClient := feedmgmt.NewClient(config.FeedsMgmtService.Host, config.FeedsMgmtService.Port, httpClient)
	articleClient := artmgmt.NewClient(config.ArticlesMgmtService.Host, config.FeedsMgmtService.Port, httpClient)

	// Setup fetcher
	fetcher := core.NewFetcher(logger, feedClient, articleClient, config.Options.CyclePeriod)

	// Spawn SIGINT/SIGTERM listener
	go lifecycle.TerminateHandler(logger, fetcher)

	// Run fetcher
	err := fetcher.Run()
	if err != nil {
		logger.Error(fmt.Sprintf("unexpected error while running fetcher: %s", err))
		return 1
	}

	logger.Info("APP gracefully terminated")
	return 0
}

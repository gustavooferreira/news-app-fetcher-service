package main

import (
	"os"

	"github.com/gustavooferreira/news-app-fetcher-service/pkg/core"
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

	logger.Info("APP gracefully terminated")
	return 0
}

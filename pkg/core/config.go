package core

import (
	"fmt"
	"os"
	"strconv"
	"strings"

	"github.com/gustavooferreira/news-app-fetcher-service/pkg/core/log"
)

// TODO: We would replace this with a proper config library like Viper.

const AppPrefix = "NEWS_APP_FETCHER"

// Configuration holds the entire configuration
type Configuration struct {
	Options             OptionsConfiguration
	FeedsMgmtService    FeedsMgmtServiceConfiguration
	ArticlesMgmtService ArticlesMgmtServiceConfiguration
}

// OptionsConfiguration holds general configuration
type OptionsConfiguration struct {
	LogLevel log.Level

	// Time to wait between cycles (in seconds)
	CyclePeriod int

	HTTPClientTimeout int
}

// FeedsMgmtServiceConfiguration holds configuration related to the feeds management service
type FeedsMgmtServiceConfiguration struct {
	Host string
	Port int
}

// ArticlesMgmtServiceConfiguration holds configuration related to the articles management service
type ArticlesMgmtServiceConfiguration struct {
	Host string
	Port int
}

// NewConfig returns new default configuration
func NewConfig() (config Configuration) {
	config.setDefaults()
	return config
}

// LoadConfig loads and validates config (from env vars)
func (config *Configuration) LoadConfig() (err error) {

	if logLevel, ok := os.LookupEnv(AppPrefix + "_OPTIONS_LOG_LEVEL"); ok {
		config.Options.LogLevel, err = ParseLogLevel(logLevel)
		if err != nil {
			return fmt.Errorf("configuration error: [options loglevel] unrecognized log level")
		}
	}

	if cyclePeriod, ok := os.LookupEnv(AppPrefix + "_OPTIONS_CYCLEPERIOD"); ok {
		config.Options.CyclePeriod, err = strconv.Atoi(cyclePeriod)
		if err != nil || config.Options.CyclePeriod <= 0 {
			return fmt.Errorf("configuration error: [options cycleperiod] input not allowed <%s>", cyclePeriod)
		}
	}

	if httpClientTimeout, ok := os.LookupEnv(AppPrefix + "_OPTIONS_HTTPCLIENTTIMEOUT"); ok {
		config.Options.HTTPClientTimeout, err = strconv.Atoi(httpClientTimeout)
		if err != nil || config.Options.HTTPClientTimeout <= 0 {
			return fmt.Errorf("configuration error: [options httpclienttimeout] input not allowed <%s>", httpClientTimeout)
		}
	}

	if feedsHost, ok := os.LookupEnv(AppPrefix + "_FEEDSMGMTSERVICE_HOST"); ok {
		config.FeedsMgmtService.Host = feedsHost
	} else {
		return fmt.Errorf("configuration error: [feedsmgmtservice host] mandatory config parameter missing")
	}

	if feedsPort, ok := os.LookupEnv(AppPrefix + "_FEEDSMGMTSERVICE_PORT"); ok {
		config.FeedsMgmtService.Port, err = strconv.Atoi(feedsPort)
		if err != nil || config.FeedsMgmtService.Port <= 0 || config.FeedsMgmtService.Port > 1<<16-1 {
			return fmt.Errorf("configuration error: [feedsmgmtservice port] input not allowed <%s>", feedsPort)
		}
	}

	if articlesHost, ok := os.LookupEnv(AppPrefix + "_ARTICLESMGMTSERVICE_HOST"); ok {
		config.ArticlesMgmtService.Host = articlesHost
	} else {
		return fmt.Errorf("configuration error: [articlesmgmtservice host] mandatory config parameter missing")
	}

	if articlesPort, ok := os.LookupEnv(AppPrefix + "_ARTICLESMGMTSERVICE_PORT"); ok {
		config.ArticlesMgmtService.Port, err = strconv.Atoi(articlesPort)
		if err != nil || config.ArticlesMgmtService.Port <= 0 || config.ArticlesMgmtService.Port > 1<<16-1 {
			return fmt.Errorf("configuration error: [articlesmgmtservice port] input not allowed <%s>", articlesPort)
		}
	}

	return nil
}

// setDefaults sets the config default values.
func (config *Configuration) setDefaults() {

	// Options
	config.Options.LogLevel = log.INFO
	config.Options.CyclePeriod = 60
	config.Options.HTTPClientTimeout = 5

	// Feeds Management Service
	config.FeedsMgmtService.Port = 8080

	// Articles Management Service
	config.ArticlesMgmtService.Port = 8080
}

// ParseLogLevel parses a string and returns a log level enum.
func ParseLogLevel(level string) (logLevel log.Level, err error) {
	level = strings.ToLower(level)

	switch level {
	case "debug":
		logLevel = log.DEBUG
	case "info":
		logLevel = log.INFO
	case "warning":
		logLevel = log.WARN
	case "error":
		logLevel = log.ERROR
	default:
		return 0, fmt.Errorf("log level unrecognised")
	}

	return logLevel, nil
}

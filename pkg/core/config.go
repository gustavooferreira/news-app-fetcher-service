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
	Options OptionsConfiguration
}

// OptionsConfiguration holds general configuration
type OptionsConfiguration struct {
	LogLevel log.Level

	// Time to wait between cycles (in seconds)
	CyclePeriod int
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

	return nil
}

// setDefaults sets the config default values.
func (config *Configuration) setDefaults() {

	// Options
	config.Options.LogLevel = log.INFO
	config.Options.CyclePeriod = 60
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

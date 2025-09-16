package agentconfig

import (
	"flag"
	"os"
	"time"
)

var FlagRunAddr string
var ReportInterval time.Duration
var PollInterval time.Duration
var SecretKey string
var RateLimit string

func ParseFlags() {
	var reportIntervalStr string
	var pollIntervalStr string
	flag.StringVar(&FlagRunAddr, "a", "localhost:8080", "address and port to run server")
	flag.StringVar(&reportIntervalStr, "r", "10", "reporting interval (e.g. 30s, 5m)")
	flag.StringVar(&SecretKey, "k", "", "SecretKey")
	flag.StringVar(&pollIntervalStr, "p", "2", "pooll")
	flag.StringVar(&RateLimit, "l", "1", "Worker Poll")
	flag.Parse()
	if envRunAddr := os.Getenv("ADDRESS"); envRunAddr != "" {
		FlagRunAddr = envRunAddr
	}

	if envReport := os.Getenv("REPORT_INTERVAL"); envReport != "" {
		reportIntervalStr = envReport
	}
	if envPool := os.Getenv("POLL_INTERVAL"); envPool != "" {
		pollIntervalStr = envPool
	}
	if envSecretKey := os.Getenv("KEY"); envSecretKey != "" {
		SecretKey = envSecretKey
	}
	if envRateLimit := os.Getenv("RATE_LIMIT"); envRateLimit != "" {
		RateLimit = envRateLimit
	}

	ReportInterval, _ = time.ParseDuration(reportIntervalStr + "s")
	PollInterval, _ = time.ParseDuration(pollIntervalStr + "s")
}

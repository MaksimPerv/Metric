package agentconfig

import (
	"flag"
	"time"
)

var FlagRunAddr string
var ReportInterval time.Duration
var PollInterval time.Duration

func ParseFlags() {
	var reportIntervalStr string
	var pollIntervalStr string
	flag.StringVar(&FlagRunAddr, "a", "http://localhost:8080", "address and port to run server")
	flag.StringVar(&reportIntervalStr, "r", "10", "reporting interval (e.g. 30s, 5m)")
	flag.StringVar(&pollIntervalStr, "p", "2", "pooll")
	flag.Parse()
	ReportInterval, _ = time.ParseDuration(reportIntervalStr)
	PollInterval, _ = time.ParseDuration(pollIntervalStr)
}

package agentconfig

import (
	"flag"
	"time"
)

var FlagRunAddr string
var ReportInterval time.Duration
var PollInterval time.Duration

func ParseFlags() {
	flag.StringVar(&FlagRunAddr, "a", "http://localhost:8080", "address and port to run server")
	flag.DurationVar(&ReportInterval, "r", 10*time.Second, "reporting interval (e.g. 30s, 5m)")
	flag.DurationVar(&PollInterval, "p", 2*time.Second, "pooll")
	flag.Parse()
}

package config

import "time"

type Config struct {
	PollInterval   time.Duration
	ReportInterval time.Duration
}

func LoadConfig() *Config {
	return &Config{
		PollInterval:   2 * time.Second,
		ReportInterval: 10 * time.Second,
	}
}

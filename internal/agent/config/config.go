package config

import (
	"flag"
	"log"
	"time"

	"github.com/caarlos0/env/v6"
)

// Config is the configuration of the application.
// Address is the address of the server.
// PollInterval is the interval of the polling data.
// ReportInterval is the interval of the agent reports.
type Config struct {
	Address        string        `env:"ADDRESS" envDefault:"localhost:8080"`
	PollInterval   time.Duration `env:"POLL_INTERVAL" envDefault:"2s"`
	ReportInterval time.Duration `env:"REPORT_INTERVAL" envDefault:"10s"`
}

func NewAgent() Config {
	var cfg Config
	// parse command line flags
	flag.StringVar(&cfg.Address, "a", cfg.Address, "address of the server")
	flag.DurationVar(&cfg.PollInterval, "p", cfg.PollInterval, "interval of the polling data")
	flag.DurationVar(&cfg.ReportInterval, "r", cfg.ReportInterval, "interval of the agent reports")
	flag.Parse()
	// parse environment variables
	if err := env.Parse(&cfg); err != nil {
		log.Printf("can't parse environment variables: %v", err)
	}
	return cfg
}

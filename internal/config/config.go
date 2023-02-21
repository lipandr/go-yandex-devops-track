package config

import (
	"flag"
	"time"

	"github.com/caarlos0/env/v6"
)

// Config is the configuration of the application.
// Address is the address of the server.
// PollInterval is the interval of the polling data.
// ReportInterval is the interval of the agent reports.
// StoreInterval is the interval of writing data to the file.
// StoreFile is the path to the file where the data will be stored.
// Restore is the option to restore the data from the file.
type Config struct {
	Address        string        `env:"ADDRESS" envDefault:"localhost:8080"`
	PollInterval   time.Duration `env:"POLL_INTERVAL" envDefault:"2s"`
	ReportInterval time.Duration `env:"REPORT_INTERVAL" envDefault:"10s"`
	StoreInterval  time.Duration `env:"STORE_INTERVAL" envDefault:"300s"`
	StoreFile      string        `env:"STORE_FILE" envDefault:"/tmp/devops-metrics-db.json"`
	Restore        bool          `env:"RESTORE" envDefault:"false"`
}

func NewServer() Config {
	var cfg Config
	if err := env.Parse(&cfg); err != nil {
		panic(err)
	}
	flag.StringVar(&cfg.Address, "a", cfg.Address, "address of the server")
	flag.BoolVar(&cfg.Restore, "r", cfg.Restore, "restore the data from the file")
	flag.DurationVar(&cfg.StoreInterval, "s", cfg.StoreInterval, "interval of writing data to the file")
	flag.StringVar(&cfg.StoreFile, "f", cfg.StoreFile, "path to the file where the data will be stored")
	flag.Parse()

	return cfg
}

func NewAgent() Config {
	var cfg Config
	if err := env.Parse(&cfg); err != nil {
		panic(err)
	}
	flag.StringVar(&cfg.Address, "a", cfg.Address, "address of the server")
	flag.DurationVar(&cfg.PollInterval, "p", cfg.PollInterval, "interval of the polling data")
	flag.DurationVar(&cfg.ReportInterval, "r", cfg.ReportInterval, "interval of the agent reports")
	flag.Parse()

	return cfg
}

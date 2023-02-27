package config

import (
	"flag"
	"github.com/caarlos0/env/v6"
	"log"
	"time"
)

// Config is the configuration of the application.
// Address is the address of the server.
// StoreInterval is the interval of writing data to the file.
// StoreFile is the path to the file where the data will be stored.
// Restore is the option to restore the data from the file.
type Config struct {
	Address       string        `env:"ADDRESS" envDefault:"localhost:8080"`
	StoreInterval time.Duration `env:"STORE_INTERVAL" envDefault:"300s"`
	StoreFile     string        `env:"STORE_FILE" envDefault:"/tmp/devops-metrics-db.json"`
	Restore       bool          `env:"RESTORE" envDefault:"true"`
}

func NewServer() *Config {
	var cfg Config
	// parse environment variables
	if err := env.Parse(&cfg); err != nil {
		log.Printf("can't parse environment variables: %v", err)
	}
	// parse command line flags
	flag.StringVar(&cfg.Address, "a", cfg.Address, "address of the server")
	flag.BoolVar(&cfg.Restore, "r", cfg.Restore, "restore the data from the file")
	flag.DurationVar(&cfg.StoreInterval, "i", cfg.StoreInterval, "interval of writing data to the file")
	flag.StringVar(&cfg.StoreFile, "f", cfg.StoreFile, "path to the file where the data will be stored")
	flag.Parse()
	return &cfg
}

// TODO указатель на конфиг нужен?
// работает с указателем и без указателя

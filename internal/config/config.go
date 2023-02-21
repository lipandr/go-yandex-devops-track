package config

import "time"

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

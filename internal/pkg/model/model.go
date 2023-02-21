package model

import "sync"

type MetricType string

const (
	TypeCounter = "counter"
	TypeGauge   = "gauge"
)

var MetricNames = []string{
	"Alloc", "BuckHashSys", "Frees", "GCCPUFraction",
	"GCSys", "HeapAlloc", "HeapIdle", "HeapInuse",
	"HeapObjects", "HeapReleased", "HeapSys", "LastGC",
	"Lookups", "MCacheInuse", "MCacheSys", "MSpanInuse",
	"MSpanSys", "Mallocs", "NextGC", "NumForcedGC",
	"NumGC", "OtherSys", "PauseTotalNs", "StackInuse",
	"StackSys", "Sys", "TotalAlloc", "PollCount",
	"RandomValue",
}

// Metric defines the metrics data to collect.
type Metric struct {
	MType MetricType
	Delta int64
	Value float64
}

// MetricJSON defines the metrics data to collect in JSON format.
type MetricJSON struct {
	ID    string     `json:"id"`
	MType MetricType `json:"type"`
	Delta *int64     `json:"delta,omitempty"`
	Value *float64   `json:"value,omitempty"`
}

// MetricData defines the metrics in memory repository.
// key is the "ID" of metric
type MetricData struct {
	Data map[string]*Metric
	MU   *sync.RWMutex
}

// MetricWeb defines the metric data for web UI.
type MetricWeb struct {
	ID    string
	Value float64
}

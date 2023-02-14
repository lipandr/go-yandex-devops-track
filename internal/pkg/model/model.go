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
	ID    string     `json:"id"`
	MType MetricType `json:"type"`
	Delta int64      `json:"delta,omitempty"`
	Value float64    `json:"value,omitempty"`
}

// MetricData defines the metrics repository.
// key is the "ID" of metric
type MetricData struct {
	Data map[string]*Metric `json:"data"`
	MU   sync.RWMutex       `json:"-"`
}

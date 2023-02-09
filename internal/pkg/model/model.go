package model

import (
	"fmt"
	"sync/atomic"
)

const (
	TypeCounter = "counter"
	TypeGauge   = "gauge"
)

type Counter interface {
	IncrementCounter(delta int64) int64
	GetCounter() int64
	PutCounter(int64)
	CounterToString() string
}
type Gauge interface {
	Update(delta float64)
	GetGauge() float64
	PutGauge(float64)
	GaugeToString() string
}

// Metric defines the key as
type Metric struct {
	Type  string      `json:"type"`
	Name  string      `json:"name"`
	Value interface{} `json:"value"`
}

// MetricData defines the metrics list to be collected by the agent
// key is the "Type" of metric
type MetricData struct {
	Data map[string]map[string]interface{}
}

func (m *Metric) IncrementCounter(delta int64) int64 {
	tmp := m.Value.(int64)
	atomic.AddInt64(&tmp, delta)
	m.Value = tmp
	return m.Value.(int64)
}
func (m *Metric) Update(delta float64) {
	m.Value = delta
}
func (m *Metric) GetCounter() int64 {
	return m.Value.(int64)
}
func (m *Metric) GetGauge() float64 {
	return m.Value.(float64)
}
func (m *Metric) PutCounter(value int64) {
	m.Value = value
}
func (m *Metric) PutGauge(value float64) {
	m.Value = value
}
func (m *Metric) CounterToString() string {
	str := fmt.Sprintf("%d", m.Value)
	return str
}
func (m *Metric) GaugeToString() string {
	return fmt.Sprintf("%f", m.Value)
}

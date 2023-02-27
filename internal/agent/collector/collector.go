package collector

import (
	"fmt"
	"math/rand"
	"runtime"
	"sync"
	"sync/atomic"

	"github.com/lipandr/go-yandex-devops-track/internal/pkg/model"
)

type Collector struct {
	collector *model.MetricData
}

func New() *Collector {
	var collector model.MetricData
	collector.Data = make(map[string]*model.Metric)

	collector.Data["PollCount"] = &model.Metric{
		MType: model.TypeCounter,
		Delta: 0,
	}
	collector.MU = &sync.RWMutex{}
	return &Collector{collector: &collector}
}

func (c *Collector) Update() {
	c.collector.MU.Lock()
	defer c.collector.MU.Unlock()

	var rtm runtime.MemStats
	runtime.ReadMemStats(&rtm)

	counter := c.collector.Data["PollCount"].Delta
	atomic.AddInt64(&counter, 1)
	c.collector.Data["PollCount"] = &model.Metric{
		MType: model.TypeCounter,
		Delta: counter,
	}
	// Create a map of primitive values to avoid creating multiple model.Metric structs
	values := map[string]float64{
		"Alloc":         float64(rtm.Alloc),
		"BuckHashSys":   float64(rtm.BuckHashSys),
		"Frees":         float64(rtm.Frees),
		"GCCPUFraction": rtm.GCCPUFraction,
		"GCSys":         float64(rtm.GCSys),
		"HeapAlloc":     float64(rtm.HeapAlloc),
		"HeapIdle":      float64(rtm.HeapIdle),
		"HeapObjects":   float64(rtm.HeapObjects),
		"HeapReleased":  float64(rtm.HeapReleased),
		"HeapSys":       float64(rtm.HeapSys),
		"LastGC":        float64(rtm.LastGC),
		"Lookups":       float64(rtm.Lookups),
		"MCacheInuse":   float64(rtm.MCacheInuse),
		"MCacheSys":     float64(rtm.MCacheSys),
		"MSpanInuse":    float64(rtm.MSpanInuse),
		"MSpanSys":      float64(rtm.MSpanSys),
		"Mallocs":       float64(rtm.Mallocs),
		"NextGC":        float64(rtm.NextGC),
		"NumForcedGC":   float64(rtm.NumForcedGC),
		"NumGC":         float64(rtm.NumGC),
		"OtherSys":      float64(rtm.OtherSys),
		"HeapInuse":     float64(rtm.HeapInuse),
		"PauseTotalNs":  float64(rtm.PauseTotalNs),
		"StackInuse":    float64(rtm.StackInuse),
		"StackSys":      float64(rtm.StackSys),
		"Sys":           float64(rtm.Sys),
		"TotalAlloc":    float64(rtm.TotalAlloc),
		"RandomValue":   rand.Float64(),
	}
	// Iterate over values and assign model.Metric
	for k, v := range values {
		c.collector.Data[k] = &model.Metric{
			MType: model.TypeGauge,
			Value: v,
		}
	}
}

// Share data using an HTTP GET request.
// All data provided in the request URL.
// Deprecated since using JSON method.
func (c *Collector) Share() []string {
	c.collector.MU.RLock()
	defer c.collector.MU.RUnlock()

	var data []string
	for k, v := range c.collector.Data {
		str := fmt.Sprintf("%f", v.Value)
		if v.MType == model.TypeCounter {
			str = fmt.Sprintf("%d", v.Delta)
		}
		url := fmt.Sprintf("http://127.0.0.1:8080/update/%v/%v/%v", v.MType, k, str)
		data = append(data, url)
	}
	return data
}

func (c *Collector) ShareJSON() []model.MetricJSON {
	c.collector.MU.RLock()
	defer c.collector.MU.RUnlock()

	var data []model.MetricJSON

	for k, v := range c.collector.Data {
		switch v.MType {
		case model.TypeGauge:
			data = append(data, model.MetricJSON{
				ID:    k,
				MType: v.MType,
				Delta: nil,
				Value: &v.Value,
			})
		case model.TypeCounter:
			data = append(data, model.MetricJSON{
				ID:    k,
				MType: v.MType,
				Delta: &v.Delta,
				Value: nil,
			})
		}
	}
	return data
}

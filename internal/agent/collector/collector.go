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
	sync.RWMutex
}

func New() *Collector {
	var collector model.MetricData
	collector.Data = make(map[string]*model.Metric)
	collector.Data["PollCount"] = &model.Metric{
		ID:    "PollCount",
		MType: model.TypeCounter,
		Delta: 0,
	}
	return &Collector{collector: &collector}
}

func (c *Collector) UpdateMetrics() {
	var rtm runtime.MemStats
	runtime.ReadMemStats(&rtm)

	c.Lock()
	defer c.Unlock()
	counter := c.collector.Data["PollCount"].Delta
	atomic.AddInt64(&counter, 1)
	c.collector.Data["PollCount"] = &model.Metric{
		ID:    "PollCount",
		MType: model.TypeCounter,
		Delta: counter,
	}
	c.collector.Data["Alloc"] = &model.Metric{
		ID:    "Alloc",
		MType: model.TypeGauge,
		Value: float64(rtm.Alloc),
	}
	c.collector.Data["BuckHashSys"] = &model.Metric{
		ID:    "BuckHashSys",
		MType: model.TypeGauge,
		Value: float64(rtm.BuckHashSys),
	}
	c.collector.Data["Frees"] = &model.Metric{
		ID:    "Frees",
		MType: model.TypeGauge,
		Value: float64(rtm.Frees),
	}
	c.collector.Data["GCCPUFraction"] = &model.Metric{
		ID:    "GCCPUFraction",
		MType: model.TypeGauge,
		Value: rtm.GCCPUFraction,
	}
	c.collector.Data["GCSys"] = &model.Metric{
		ID:    "GCSys",
		MType: model.TypeGauge,
		Value: float64(rtm.GCSys),
	}
	c.collector.Data["HeapAlloc"] = &model.Metric{
		ID:    "HeapAlloc",
		MType: model.TypeGauge,
		Value: float64(rtm.HeapAlloc),
	}
	c.collector.Data["HeapIdle"] = &model.Metric{
		ID:    "HeapIdle",
		MType: model.TypeGauge,
		Value: float64(rtm.HeapIdle),
	}
	c.collector.Data["HeapObjects"] = &model.Metric{
		ID:    "HeapObjects",
		MType: model.TypeGauge,
		Value: float64(rtm.HeapObjects),
	}
	c.collector.Data["HeapReleased"] = &model.Metric{
		ID:    "HeapReleased",
		MType: model.TypeGauge,
		Value: float64(rtm.HeapReleased),
	}
	c.collector.Data["HeapSys"] = &model.Metric{
		ID:    "HeapSys",
		MType: model.TypeGauge,
		Value: float64(rtm.HeapSys),
	}
	c.collector.Data["LastGC"] = &model.Metric{
		ID:    "LastGC",
		MType: model.TypeGauge,
		Value: float64(rtm.LastGC),
	}
	c.collector.Data["Lookups"] = &model.Metric{
		ID:    "Lookups",
		MType: model.TypeGauge,
		Value: float64(rtm.Lookups),
	}
	c.collector.Data["MCacheInuse"] = &model.Metric{
		ID:    "MCacheInuse",
		MType: model.TypeGauge,
		Value: float64(rtm.MCacheInuse),
	}
	c.collector.Data["MCacheSys"] = &model.Metric{
		ID:    "MCacheSys",
		MType: model.TypeGauge,
		Value: float64(rtm.MCacheSys),
	}
	c.collector.Data["MSpanInuse"] = &model.Metric{
		ID:    "MSpanInuse",
		MType: model.TypeGauge,
		Value: float64(rtm.MSpanInuse),
	}
	c.collector.Data["MSpanSys"] = &model.Metric{
		ID:    "MSpanSys",
		MType: model.TypeGauge,
		Value: float64(rtm.MSpanSys),
	}
	c.collector.Data["Mallocs"] = &model.Metric{
		ID:    "Mallocs",
		MType: model.TypeGauge,
		Value: float64(rtm.Mallocs),
	}
	c.collector.Data["NextGC"] = &model.Metric{
		ID:    "NextGC",
		MType: model.TypeGauge,
		Value: float64(rtm.NextGC),
	}
	c.collector.Data["NumForcedGC"] = &model.Metric{
		ID:    "NumForcedGC",
		MType: model.TypeGauge,
		Value: float64(rtm.NumForcedGC),
	}
	c.collector.Data["NumGC"] = &model.Metric{
		ID:    "NumGC",
		MType: model.TypeGauge,
		Value: float64(rtm.NumGC),
	}
	c.collector.Data["OtherSys"] = &model.Metric{
		ID:    "OtherSys",
		MType: model.TypeGauge,
		Value: float64(rtm.OtherSys),
	}
	c.collector.Data["HeapInuse"] = &model.Metric{
		ID:    "HeapInuse",
		MType: model.TypeGauge,
		Value: float64(rtm.HeapInuse),
	}
	c.collector.Data["PauseTotalNs"] = &model.Metric{
		ID:    "PauseTotalNs",
		MType: model.TypeGauge,
		Value: float64(rtm.PauseTotalNs),
	}
	c.collector.Data["StackInuse"] = &model.Metric{
		ID:    "StackInuse",
		MType: model.TypeGauge,
		Value: float64(rtm.StackInuse),
	}
	c.collector.Data["StackSys"] = &model.Metric{
		ID:    "StackSys",
		MType: model.TypeGauge,
		Value: float64(rtm.StackSys),
	}
	c.collector.Data["Sys"] = &model.Metric{
		ID:    "Sys",
		MType: model.TypeGauge,
		Value: float64(rtm.Sys),
	}
	c.collector.Data["TotalAlloc"] = &model.Metric{
		ID:    "TotalAlloc",
		MType: model.TypeGauge,
		Value: float64(rtm.TotalAlloc),
	}
	c.collector.Data["RandomValue"] = &model.Metric{
		ID:    "RandomValue",
		MType: model.TypeGauge,
		Value: rand.Float64(),
	}
}
func (c *Collector) ShareMetrics() []string {
	c.RLock()
	defer c.RUnlock()

	var data []string
	for _, v := range c.collector.Data {
		str := fmt.Sprintf("%f", v.Value)
		if v.MType == model.TypeCounter {
			str = fmt.Sprintf("%d", v.Delta)
		}
		url := fmt.Sprintf("http://127.0.0.1:8080/update/%v/%v/%v", v.MType, v.ID, str)
		data = append(data, url)
	}
	return data

}

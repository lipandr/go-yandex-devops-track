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
	collector.Data = make(map[string]map[string]interface{})
	collector.Data[model.TypeCounter] = make(map[string]interface{})
	collector.Data[model.TypeGauge] = make(map[string]interface{})
	collector.Data[model.TypeCounter]["PollCount"] = int64(0)

	return &Collector{collector: &collector}
}

func (c *Collector) UpdateMetrics() {
	var rtm runtime.MemStats
	runtime.ReadMemStats(&rtm)

	c.Lock()
	defer c.Unlock()
	counter := c.collector.Data[model.TypeCounter]["PollCount"].(int64)
	atomic.AddInt64(&counter, 1)
	c.collector.Data[model.TypeCounter]["PollCount"] = counter

	c.collector.Data[model.TypeGauge]["Alloc"] = float64(rtm.Alloc)
	c.collector.Data[model.TypeGauge]["BuckHashSys"] = float64(rtm.BuckHashSys)
	c.collector.Data[model.TypeGauge]["Frees"] = float64(rtm.Frees)
	c.collector.Data[model.TypeGauge]["GCCPUFraction"] = rtm.GCCPUFraction
	c.collector.Data[model.TypeGauge]["GCSys"] = float64(rtm.GCSys)
	c.collector.Data[model.TypeGauge]["HeapAlloc"] = float64(rtm.HeapAlloc)
	c.collector.Data[model.TypeGauge]["HeapIdle"] = float64(rtm.HeapIdle)
	c.collector.Data[model.TypeGauge]["HeapObjects"] = float64(rtm.HeapObjects)
	c.collector.Data[model.TypeGauge]["HeapReleased"] = float64(rtm.HeapReleased)
	c.collector.Data[model.TypeGauge]["HeapSys"] = float64(rtm.HeapSys)
	c.collector.Data[model.TypeGauge]["LastGC"] = float64(rtm.LastGC)
	c.collector.Data[model.TypeGauge]["Lookups"] = float64(rtm.Lookups)
	c.collector.Data[model.TypeGauge]["MCacheInuse"] = float64(rtm.MCacheInuse)
	c.collector.Data[model.TypeGauge]["MCacheSys"] = float64(rtm.MCacheSys)
	c.collector.Data[model.TypeGauge]["MSpanInuse"] = float64(rtm.MSpanInuse)
	c.collector.Data[model.TypeGauge]["MSpanSys"] = float64(rtm.MSpanSys)
	c.collector.Data[model.TypeGauge]["Mallocs"] = float64(rtm.Mallocs)
	c.collector.Data[model.TypeGauge]["NextGC"] = float64(rtm.NextGC)
	c.collector.Data[model.TypeGauge]["NumForcedGC"] = float64(rtm.NumForcedGC)
	c.collector.Data[model.TypeGauge]["NumGC"] = float64(rtm.NumGC)
	c.collector.Data[model.TypeGauge]["OtherSys"] = float64(rtm.OtherSys)
	c.collector.Data[model.TypeGauge]["NumGC"] = float64(rtm.NumGC)
	c.collector.Data[model.TypeGauge]["PauseTotalNs"] = float64(rtm.PauseTotalNs)
	c.collector.Data[model.TypeGauge]["StackInuse"] = float64(rtm.StackInuse)
	c.collector.Data[model.TypeGauge]["StackSys"] = float64(rtm.StackSys)
	c.collector.Data[model.TypeGauge]["Sys"] = float64(rtm.Sys)
	c.collector.Data[model.TypeGauge]["TotalAlloc"] = float64(rtm.TotalAlloc)
	c.collector.Data[model.TypeGauge]["RandomValue"] = rand.Float64()
}
func (c *Collector) ShareMetrics() []string {
	c.RLock()
	defer c.RUnlock()

	var data []string
	m := c.collector
	for k, v := range m.Data {
		for nk, nv := range v {
			str := fmt.Sprintf("http://127.0.0.1:8080/update/%v/%v/%v", k, nk, nv)
			data = append(data, str)
		}
	}
	return data

}

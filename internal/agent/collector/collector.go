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
	c.collector.Data["Alloc"] = &model.Metric{
		MType: model.TypeGauge,
		Value: float64(rtm.Alloc),
	}

	c.collector.Data["BuckHashSys"] = &model.Metric{
		MType: model.TypeGauge,
		Value: float64(rtm.BuckHashSys),
	}
	c.collector.Data["Frees"] = &model.Metric{
		MType: model.TypeGauge,
		Value: float64(rtm.Frees),
	}
	c.collector.Data["GCCPUFraction"] = &model.Metric{
		MType: model.TypeGauge,
		Value: rtm.GCCPUFraction,
	}
	c.collector.Data["GCSys"] = &model.Metric{
		MType: model.TypeGauge,
		Value: float64(rtm.GCSys),
	}
	c.collector.Data["HeapAlloc"] = &model.Metric{
		MType: model.TypeGauge,
		Value: float64(rtm.HeapAlloc),
	}
	c.collector.Data["HeapIdle"] = &model.Metric{
		MType: model.TypeGauge,
		Value: float64(rtm.HeapIdle),
	}
	c.collector.Data["HeapObjects"] = &model.Metric{
		MType: model.TypeGauge,
		Value: float64(rtm.HeapObjects),
	}
	c.collector.Data["HeapReleased"] = &model.Metric{
		MType: model.TypeGauge,
		Value: float64(rtm.HeapReleased),
	}
	c.collector.Data["HeapSys"] = &model.Metric{
		MType: model.TypeGauge,
		Value: float64(rtm.HeapSys),
	}
	c.collector.Data["LastGC"] = &model.Metric{
		MType: model.TypeGauge,
		Value: float64(rtm.LastGC),
	}
	c.collector.Data["Lookups"] = &model.Metric{
		MType: model.TypeGauge,
		Value: float64(rtm.Lookups),
	}
	c.collector.Data["MCacheInuse"] = &model.Metric{
		MType: model.TypeGauge,
		Value: float64(rtm.MCacheInuse),
	}
	c.collector.Data["MCacheSys"] = &model.Metric{
		MType: model.TypeGauge,
		Value: float64(rtm.MCacheSys),
	}
	c.collector.Data["MSpanInuse"] = &model.Metric{
		MType: model.TypeGauge,
		Value: float64(rtm.MSpanInuse),
	}
	c.collector.Data["MSpanSys"] = &model.Metric{
		MType: model.TypeGauge,
		Value: float64(rtm.MSpanSys),
	}
	c.collector.Data["Mallocs"] = &model.Metric{
		MType: model.TypeGauge,
		Value: float64(rtm.Mallocs),
	}
	c.collector.Data["NextGC"] = &model.Metric{
		MType: model.TypeGauge,
		Value: float64(rtm.NextGC),
	}
	c.collector.Data["NumForcedGC"] = &model.Metric{
		MType: model.TypeGauge,
		Value: float64(rtm.NumForcedGC),
	}
	c.collector.Data["NumGC"] = &model.Metric{
		MType: model.TypeGauge,
		Value: float64(rtm.NumGC),
	}
	c.collector.Data["OtherSys"] = &model.Metric{
		MType: model.TypeGauge,
		Value: float64(rtm.OtherSys),
	}
	c.collector.Data["HeapInuse"] = &model.Metric{
		MType: model.TypeGauge,
		Value: float64(rtm.HeapInuse),
	}
	c.collector.Data["PauseTotalNs"] = &model.Metric{
		MType: model.TypeGauge,
		Value: float64(rtm.PauseTotalNs),
	}
	c.collector.Data["StackInuse"] = &model.Metric{
		MType: model.TypeGauge,
		Value: float64(rtm.StackInuse),
	}
	c.collector.Data["StackSys"] = &model.Metric{
		MType: model.TypeGauge,
		Value: float64(rtm.StackSys),
	}
	c.collector.Data["Sys"] = &model.Metric{
		MType: model.TypeGauge,
		Value: float64(rtm.Sys),
	}
	c.collector.Data["TotalAlloc"] = &model.Metric{
		MType: model.TypeGauge,
		Value: float64(rtm.TotalAlloc),
	}
	c.collector.Data["RandomValue"] = &model.Metric{
		MType: model.TypeGauge,
		Value: rand.Float64(),
	}
}

// Share data using a HTTP GET request.
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

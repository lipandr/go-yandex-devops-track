package controller

import (
	"time"

	"github.com/lipandr/go-yandex-devops-track/internal/agent/collector"
)

const pollInterval = 2 * time.Second

type Controller struct {
	collector *collector.Collector
}

func New(collector *collector.Collector) *Controller {
	return &Controller{
		collector: collector,
	}
}
func (c *Controller) CollectData() {
	ticker := time.NewTicker(pollInterval)
	quit := make(chan struct{})
	go func() {
		for {
			select {
			case <-ticker.C:
				c.collector.UpdateMetrics()
			case <-quit:
				ticker.Stop()
				return
			}
		}
	}()
}
func (c *Controller) ReportData() []string {
	return c.collector.ShareMetrics()
}

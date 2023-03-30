package controller

import (
	"github.com/lipandr/go-yandex-devops-track/internal/agent/config"
	"time"

	"github.com/lipandr/go-yandex-devops-track/internal/agent/collector"
)

type Controller struct {
	collector *collector.Collector
	config    *config.Config
}

func New(collector *collector.Collector, config *config.Config) *Controller {
	return &Controller{
		collector: collector,
		config:    config,
	}
}
func (c *Controller) CollectData() {
	ticker := time.NewTicker(c.config.PollInterval)
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

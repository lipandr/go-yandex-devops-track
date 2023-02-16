package controller

import (
	"time"

	"github.com/lipandr/go-yandex-devops-track/internal/agent/collector"
	"github.com/lipandr/go-yandex-devops-track/internal/agent/config"
	"github.com/lipandr/go-yandex-devops-track/internal/pkg/model"
)

type Controller struct {
	collector *collector.Collector
	config    *config.Config
}

func New(collector *collector.Collector, cfg *config.Config) *Controller {
	return &Controller{
		collector: collector,
		config:    cfg,
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
func (c *Controller) ReportJSON() []model.MetricJSON {
	return c.collector.ShareMetricsJSON()
}

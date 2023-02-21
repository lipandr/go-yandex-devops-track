package main

import (
	"context"
	"time"

	"github.com/lipandr/go-yandex-devops-track/internal/agent/collector"
	"github.com/lipandr/go-yandex-devops-track/internal/agent/controller"
	"github.com/lipandr/go-yandex-devops-track/internal/agent/handler/http"
	"github.com/lipandr/go-yandex-devops-track/internal/config"
)

func main() {
	cfg := config.NewAgent()
	ctx := context.Background()
	col := collector.New()
	ctl := controller.New(col, &cfg)
	ctl.CollectData()

	h := http.New(ctl, &cfg)

	ticker := time.NewTicker(cfg.ReportInterval)
	quit := make(chan struct{})
	for {
		select {
		case <-ticker.C:
			h.Run(ctx)
		case <-quit:
			ticker.Stop()
			return
		}
	}
}

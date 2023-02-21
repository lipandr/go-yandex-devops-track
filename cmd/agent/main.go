package main

import (
	"context"
	"time"

	"github.com/caarlos0/env/v6"

	"github.com/lipandr/go-yandex-devops-track/internal/agent/collector"
	"github.com/lipandr/go-yandex-devops-track/internal/agent/controller"
	"github.com/lipandr/go-yandex-devops-track/internal/agent/handler/http"
	"github.com/lipandr/go-yandex-devops-track/internal/config"
)

func main() {
	var cfg config.Config
	if err := env.Parse(&cfg); err != nil {
		panic(err)
	}
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

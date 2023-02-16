package main

import (
	"context"
	"time"

	"github.com/caarlos0/env/v6"

	"github.com/lipandr/go-yandex-devops-track/internal/agent/collector"
	"github.com/lipandr/go-yandex-devops-track/internal/agent/config"
	"github.com/lipandr/go-yandex-devops-track/internal/agent/controller"
	"github.com/lipandr/go-yandex-devops-track/internal/agent/handler"
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

	h := handler.New(ctl, &cfg)

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

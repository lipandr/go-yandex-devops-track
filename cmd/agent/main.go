package main

import (
	"context"
	"time"

	"github.com/lipandr/go-yandex-devops-track/internal/agent/collector"
	"github.com/lipandr/go-yandex-devops-track/internal/agent/controller"
	"github.com/lipandr/go-yandex-devops-track/internal/agent/handler"
)

const reportInterval = 10 * time.Second

func main() {
	ctx := context.Background()
	col := collector.New()
	ctl := controller.New(col)
	ctl.CollectData()

	h := handler.New(ctl)

	ticker := time.NewTicker(reportInterval)
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

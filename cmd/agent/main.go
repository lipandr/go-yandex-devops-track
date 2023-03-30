package main

import (
	"context"
	"github.com/lipandr/go-yandex-devops-track/internal/agent/config"
	"log"
	"time"

	"github.com/lipandr/go-yandex-devops-track/internal/agent/collector"
	"github.com/lipandr/go-yandex-devops-track/internal/agent/controller"
	"github.com/lipandr/go-yandex-devops-track/internal/agent/handler"
)

const reportInterval = 10 * time.Second

func main() {
	cfg := config.New()
	ctx := context.Background()
	col := collector.New()
	ctl := controller.New(col, cfg)
	ctl.CollectData()

	h := handler.New(ctl, cfg)

	log.Println("Starting agent....")
	log.Printf("Remote server address: %s", cfg.Address)
	log.Printf("Polling interval: %s", cfg.PollInterval)
	log.Printf("Report interval: %s", cfg.ReportInterval)

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

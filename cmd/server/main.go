package main

import (
	"context"
	"log"
	"net/http"
	"time"

	"github.com/caarlos0/env/v6"
	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"

	"github.com/lipandr/go-yandex-devops-track/internal/agent/config"
	"github.com/lipandr/go-yandex-devops-track/internal/server/controller"
	httpHandler "github.com/lipandr/go-yandex-devops-track/internal/server/handler/http"
	"github.com/lipandr/go-yandex-devops-track/internal/server/storage/memory"
)

func main() {
	var cfg config.Config
	if err := env.Parse(&cfg); err != nil {
		log.Fatal(err)
	}
	ctx := context.Background()
	repo := memory.New()
	ctl := controller.New(repo)
	h := httpHandler.New(ctx, ctl)

	server := &http.Server{
		Addr:         cfg.Address,
		Handler:      service(h),
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 5 * time.Second,
	}
	// Run the server
	if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		log.Fatal(err)
	}
}

// service is a http.Handler that implements the http.Handler interface.
func service(h *httpHandler.Handler) http.Handler {
	r := chi.NewRouter()
	//r.Use(middleware.RequestID)
	r.Use(middleware.Logger)

	r.Get("/value/*", h.GetMetricValue)
	r.Post("/value/", h.GetMetricValueJSON)
	r.Post("/update/", h.PutMetricJSON)
	r.Post("/update/*", h.PutMetric)
	r.Get("/", h.ListAllMetrics)

	return r
}

package main

import (
	"context"
	"github.com/lipandr/go-yandex-devops-track/internal/server/config"
	"log"
	"net/http"
	"time"

	"github.com/go-chi/chi"
	"github.com/lipandr/go-yandex-devops-track/internal/server/controller"
	httpHandler "github.com/lipandr/go-yandex-devops-track/internal/server/handler/http"
	"github.com/lipandr/go-yandex-devops-track/internal/server/storage/memory"
)

func main() {
	cfg := config.New()
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
	log.Printf("Starting server on port %s....", cfg.Address)
	if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		log.Fatal(err)
	}
}
func service(h *httpHandler.Handler) http.Handler {
	r := chi.NewRouter()
	//r.Use(middleware.RequestID)
	//r.Use(middleware.Logger)
	r.Get("/value/*", h.GetMetricValue)
	r.Post("/update/*", h.PutMetric)
	r.Get("/", h.ListAllMetrics)

	return r
}

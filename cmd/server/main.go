package main

import (
	"context"
	"log"
	"net/http"
	"time"

	"github.com/caarlos0/env/v6"
	"github.com/go-chi/chi"

	"github.com/lipandr/go-yandex-devops-track/internal/config"
	"github.com/lipandr/go-yandex-devops-track/internal/server/controller"
	httpHandler "github.com/lipandr/go-yandex-devops-track/internal/server/handler/http"
	"github.com/lipandr/go-yandex-devops-track/internal/server/storage/file"
	"github.com/lipandr/go-yandex-devops-track/internal/server/storage/memory"
)

func main() {
	var cfg config.Config
	if err := env.Parse(&cfg); err != nil {
		log.Fatal(err)
	}
	ctx := context.Background()
	inMemoryRepo := memory.New()
	ctl := controller.NewMemoryRepo(inMemoryRepo)
	// if Read is true, the data will be restored from the file.
	if cfg.Restore {
		r, err := file.NewFileReader(cfg.StoreFile)
		if err != nil {
			log.Fatal(err)
		}
		ctl = controller.NewFileRepo(inMemoryRepo, r)
		if err = ctl.Read(ctx); err != nil {
			log.Printf("restore error: %v", err)
		}
	}
	h := httpHandler.New(ctx, ctl)

	server := &http.Server{
		Addr:    cfg.Address,
		Handler: service(h),
	}
	// Run the server
	go func() {
		if err := server.ListenAndServe(); err != nil {
			log.Fatal(err)
		}
	}()
	// trying to save the data to the file at StoreInterval time interval
	ticker := time.NewTicker(cfg.StoreInterval)
	quit := make(chan struct{})
	for {
		select {
		case <-ticker.C:
			w, err := file.NewFileWriter(cfg.StoreFile)
			if err != nil {
				log.Fatal(err)
			}
			ctl = controller.NewFileRepo(inMemoryRepo, w)
			if err = ctl.Write(ctx); err != nil {
				log.Printf("write to file error: %v", err)
			}
		case <-quit:
			ticker.Stop()
			return
		}
	}
}

// service is a http.Handler that implements the http.Handler interface.
func service(h *httpHandler.Handler) http.Handler {
	r := chi.NewRouter()
	//r.Use(middleware.RequestID)
	//r.Use(middleware.Logger)

	r.Get("/value/*", h.GetValue)
	r.Post("/value/", h.GetValueJSON)
	r.Post("/update/", h.UpdateJSON)
	r.Post("/update/*", h.Update)
	r.Get("/", h.UIListAll)

	return r
}

package main

import (
	"context"
	"log"
	"net/http"
	"time"

	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"

	"github.com/lipandr/go-yandex-devops-track/internal/server/config"
	"github.com/lipandr/go-yandex-devops-track/internal/server/controller"
	httpHandler "github.com/lipandr/go-yandex-devops-track/internal/server/handler/http"
	"github.com/lipandr/go-yandex-devops-track/internal/server/storage/file"
	"github.com/lipandr/go-yandex-devops-track/internal/server/storage/memory"
)

func main() {
	cfg := config.NewServer()
	ctx := context.Background()
	inMemoryRepo := memory.New()
	ctl := controller.NewMemoryRepo(inMemoryRepo)
	// if Restore is true, attempt to restore data from the file.
	if cfg.Restore {
		r, err := file.NewFileReader(cfg.StoreFile)
		if err != nil {
			log.Printf("error: %v", err)
		} else {
			ctl = controller.NewFileRepo(inMemoryRepo, r)
			if err = ctl.Read(ctx); err != nil {
				log.Printf("restore data from file error: %v", err)
			}
			log.Printf("restored data from file: %s", cfg.StoreFile)
		}
	}
	h := httpHandler.New(ctx, ctl)

	server := &http.Server{
		Addr:    cfg.Address,
		Handler: service(h),
	}
	// Run the server
	go func() {
		log.Printf("starting server on %s", cfg.Address)
		log.Fatal(server.ListenAndServe())
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
	r.Use(middleware.Logger)
	//r.Use(middleware.NewCompressor(gzip.DefaultCompression).Handler)

	r.Get("/value/*", h.GetValue)
	r.Post("/value/", h.GetValueJSON)
	r.Post("/update/", h.UpdateJSON)
	r.Post("/update/*", h.Update)
	r.Get("/", h.UIListAll)

	return r
}

package handler

import (
	"context"
	"sync"

	"github.com/go-resty/resty/v2"

	"github.com/lipandr/go-yandex-devops-track/internal/agent/controller"
)

type Handler struct {
	controller *controller.Controller

	client *resty.Client
}

func New(controller *controller.Controller) *Handler {
	client := resty.New()
	return &Handler{
		controller: controller,
		client:     client,
	}
}

func (h *Handler) Run(_ context.Context) {
	var wg sync.WaitGroup

	urls := h.controller.ReportData()

	for _, u := range urls {
		wg.Add(1)
		go func(url string) {
			defer wg.Done()

			_, err := h.client.R().
				SetHeader("Content-MType", "text/plain").Post(url)
			if err != nil {
				return
			}

		}(u)
	}
}

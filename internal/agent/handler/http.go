package handler

import (
	"context"
	"fmt"
	"sync"

	"github.com/go-resty/resty/v2"

	"github.com/lipandr/go-yandex-devops-track/internal/agent/config"
	"github.com/lipandr/go-yandex-devops-track/internal/agent/controller"
)

type Handler struct {
	controller *controller.Controller
	config     *config.Config
	client     *resty.Client
}

func New(controller *controller.Controller, config *config.Config) *Handler {
	client := resty.New()
	return &Handler{
		controller: controller,
		config:     config,
		client:     client,
	}
}

func (h *Handler) Run(_ context.Context) {
	var wg sync.WaitGroup

	data := h.controller.ReportData()

	for _, d := range data {
		wg.Add(1)
		go func(str string) {
			defer wg.Done()
			url := fmt.Sprintf("http://%s/update/%s", h.config.Address, str)
			_, err := h.client.R().
				SetHeader("Content-MType", "text/plain").Post(url)
			if err != nil {
				return
			}

		}(d)
	}
}

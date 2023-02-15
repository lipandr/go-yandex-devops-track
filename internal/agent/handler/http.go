package handler

import (
	"context"
	"encoding/json"
	"github.com/lipandr/go-yandex-devops-track/internal/pkg/model"
	"log"
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

	data := h.controller.ReportJSON()

	for _, d := range data {
		wg.Add(1)
		go func(data model.MetricJSON) {
			defer wg.Done()
			b, err := json.Marshal(data)
			if err != nil {
				return
			}
			_, err = h.client.R().
				SetHeader("Content-Type", "application/json; charset=utf-8").
				SetBody(b).
				Post("http://localhost:8080/update")
			if err != nil {
				log.Printf("error: %v\n", err)
				return
			}

		}(d)
	}
}

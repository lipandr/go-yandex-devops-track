package http

import (
	"context"
	"fmt"
	"log"

	"github.com/go-resty/resty/v2"

	"github.com/lipandr/go-yandex-devops-track/internal/agent/controller"
	"github.com/lipandr/go-yandex-devops-track/internal/config"
)

// Handler is a struct that contains the data of the handler
type Handler struct {
	controller *controller.Controller
	client     *resty.Client
	config     *config.Config
}

// New returns a new handler.
func New(controller *controller.Controller, cfg *config.Config) *Handler {
	client := resty.New().
		SetHeader("Content-Type", "application/json").
		EnableTrace()
	return &Handler{
		controller: controller,
		client:     client,
		config:     cfg,
	}
}

func (h *Handler) Run(_ context.Context) {
	data := h.controller.ReportJSON()

	for _, val := range data {
		buf, err := h.client.JSONMarshal(val)
		if err != nil {
			log.Println(err)
			continue
		}
		_, err = h.client.R().
			SetBody(buf).
			Post(fmt.Sprintf("http://%s/update/", h.config.Address))
		if err != nil {
			log.Printf("failed to send data %v to %s: %v", val, h.config.Address, err)
		}
	}
}

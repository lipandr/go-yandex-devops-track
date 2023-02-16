package handler

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"sync"

	"github.com/go-resty/resty/v2"

	"github.com/lipandr/go-yandex-devops-track/internal/agent/config"
	"github.com/lipandr/go-yandex-devops-track/internal/agent/controller"
	"github.com/lipandr/go-yandex-devops-track/internal/pkg/model"
)

// Handler is a struct that contains the data of the handler
type Handler struct {
	controller *controller.Controller
	client     *resty.Client
	config     *config.Config
}

// New returns a new handler.
func New(controller *controller.Controller, cfg *config.Config) *Handler {
	client := resty.New()
	return &Handler{
		controller: controller,
		client:     client,
		config:     cfg,
	}
}

func (h *Handler) Run(_ context.Context) {
	var wg sync.WaitGroup

	data := h.controller.ReportJSON()

	for _, d := range data {
		wg.Add(1)
		go func(data model.MetricJSON) {
			defer wg.Done()

			var buf bytes.Buffer

			jsonEncoder := json.NewEncoder(&buf)
			err := jsonEncoder.Encode(data)
			if err != nil {
				return
			}
			//b, err := json.Marshal(data)
			//if err != nil {
			//	return
			//}
			resp, err := h.client.R().
				SetHeader("Content-Type", "application/json").
				SetBody(buf.Bytes()).
				Post(fmt.Sprintf("http://%s/update/", h.config.Address))
			if err != nil {
				//log.Printf("error: %v\n", err)
				return
			}
			defer func() {
				_ = resp.RawBody().Close()
			}()
		}(d)
	}
}

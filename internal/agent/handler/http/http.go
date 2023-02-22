package http

import (
	"bytes"
	"context"
	"encoding/json"
	"log"
	"net/http"
	"net/url"

	"github.com/lipandr/go-yandex-devops-track/internal/agent/config"
	"github.com/lipandr/go-yandex-devops-track/internal/agent/controller"
)

// Handler is a struct that contains the data of the handler
type Handler struct {
	controller *controller.Controller
	client     *http.Client
	config     *config.Config
}

// New returns a new handler.
func New(controller *controller.Controller, cfg *config.Config) *Handler {
	client := &http.Client{}
	return &Handler{
		controller: controller,
		client:     client,
		config:     cfg,
	}
}

func (h *Handler) Run(_ context.Context) {
	data := h.controller.ReportJSON()
	u := url.URL{
		Scheme: "http",
		Host:   h.config.Address,
		Path:   "/update/",
	}
	for _, val := range data {

		jsonBody, err := json.Marshal(val)
		if err != nil {
			log.Println(err)
			continue
		}
		responseBody := bytes.NewBuffer(jsonBody)
		request, err := http.NewRequest(http.MethodPost, u.String(), responseBody)
		if err != nil {
			log.Println(err)
			continue
		}
		request.Header.Set("Content-Type", "application/json")
		response, err := h.client.Do(request)
		if err != nil {
			log.Println(err)
			continue
		}

		if response.StatusCode != http.StatusOK {
			log.Println(response.StatusCode)
			continue
		}
		_ = response.Body.Close()
	}
}

package http

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"github.com/lipandr/go-yandex-devops-track/internal/agent/config"
	"github.com/lipandr/go-yandex-devops-track/internal/agent/controller"
	"log"
	"net/http"
)

// ClientHTTP is interface for http client
type ClientHTTP interface {
	Do(req *http.Request) (*http.Response, error)
}

var client ClientHTTP

func init() {
	client = &http.Client{}
}

// Handler is a struct that contains the data of the handler
type Handler struct {
	controller *controller.Controller
	client     ClientHTTP
	config     *config.Config
}

// New returns a new handler.
func New(controller *controller.Controller, cfg *config.Config) *Handler {
	return &Handler{
		controller: controller,
		client:     client,
		config:     cfg,
	}
}

// Run starts the internal agent logic.
func (h *Handler) Run(_ context.Context) {
	// Get data from the controller
	data := h.controller.ReportJSON()
	url := fmt.Sprintf("http://%s/update/", h.config.Address)
	// Make requests to the server
	for _, val := range data {
		response, err := Post(url, val)
		if err != nil {
			log.Printf("error: %v", err)
			continue
		}
		// Close the response body
		if err = response.Body.Close(); err != nil {
			log.Printf("error: %v", err)
		}
	}
}

func Post(url string, body interface{}) (*http.Response, error) {
	// Marshal the data as JSON
	jsonBytes, err := json.Marshal(body)
	if err != nil {
		return nil, err
	}
	// Create buffer to send the data
	responseBody := bytes.NewBuffer(jsonBytes)
	// New request
	request, err := http.NewRequest(http.MethodPost, url, responseBody)
	if err != nil {
		return nil, err
	}
	// Set the content type
	request.Header.Set("Content-Type", "application/json")

	return client.Do(request)
}

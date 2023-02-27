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

// ClientHTTP is interface for http Client
type ClientHTTP interface {
	Do(req *http.Request) (*http.Response, error)
}

var Client ClientHTTP

func init() {
	Client = &http.Client{}
}

// Handler is a struct that contains the data of the handler
type Handler struct {
	controller *controller.Controller
	client     ClientHTTP
	config     config.Config
}

// New returns a new handler.
func New(controller *controller.Controller, cfg config.Config) *Handler {
	return &Handler{
		controller: controller,
		client:     Client,
		config:     cfg,
	}
}

func (h *Handler) Run(_ context.Context) {
	// Get data from the controller
	data := h.controller.ReportJSON()
	u := url.URL{
		Scheme: "http",
		Host:   h.config.Address,
		Path:   "/update/",
	}
	// Make requests to the server
	for _, val := range data {
		response, err := Post(u, val)
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

func Post(url url.URL, body interface{}) (*http.Response, error) {
	// Marshal the data as JSON
	jsonBytes, err := json.Marshal(body)
	if err != nil {
		return nil, err
	}
	// Create buffer to send the data
	responseBody := bytes.NewBuffer(jsonBytes)
	// New request
	request, err := http.NewRequest(http.MethodPost, url.String(), responseBody)
	if err != nil {
		return nil, err
	}
	// Set the content type
	request.Header.Set("Content-Type", "application/json")

	return Client.Do(request)
}

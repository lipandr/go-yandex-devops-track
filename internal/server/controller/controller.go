package controller

import (
	"context"
	"errors"

	"github.com/lipandr/go-yandex-devops-track/internal/pkg/model"
	"github.com/lipandr/go-yandex-devops-track/internal/server/storage"
)

var ErrNotFound = errors.New("not found")

type memoryRepository interface {
	Get(ctx context.Context, name string) (string, error)
	Put(ctx context.Context, name string, metric *model.Metric) error
	GetAll(ctx context.Context) []model.MetricWeb
	GetAllJSON(ctx context.Context) []model.MetricJSON
}

type fileRepository interface {
	Write(metric *model.MetricJSON) error
	Read() (*model.MetricJSON, error)
	Close() error
}

type Controller struct {
	memory memoryRepository
	file   fileRepository
}

func NewMemoryRepo(repository memoryRepository) *Controller {
	return &Controller{
		memory: repository,
	}
}

func NewFileRepo(mem memoryRepository, repository fileRepository) *Controller {
	return &Controller{
		memory: mem,
		file:   repository,
	}
}

func (c *Controller) Get(ctx context.Context, name string) (string, error) {
	res, err := c.memory.Get(ctx, name)
	if err != nil && errors.Is(err, storage.ErrNotFound) {
		return "", ErrNotFound
	}
	return res, nil
}

func (c *Controller) Put(ctx context.Context, name string, metric *model.Metric) error {
	return c.memory.Put(ctx, name, metric)
}

func (c *Controller) GetAll(ctx context.Context) []model.MetricWeb {
	return c.memory.GetAll(ctx)
}

func (c *Controller) GetAllJSON(ctx context.Context) []model.MetricJSON {
	return c.memory.GetAllJSON(ctx)
}

func (c *Controller) Write(ctx context.Context) error {
	data := c.GetAllJSON(ctx)
	for _, metric := range data {
		if err := c.file.Write(&metric); err != nil {
			return err
		}
	}
	return c.file.Close()
}

func (c *Controller) Read(ctx context.Context) error {
	for {
		metricJSON, err := c.file.Read()
		if err != nil {
			break
		}
		name, metric := c.FromJSON(metricJSON)
		if err = c.memory.Put(ctx, name, metric); err != nil {
			return err
		}
	}
	return c.file.Close()
}

func (c *Controller) ToJSON(name string, metric model.Metric) *model.MetricJSON {
	metricJSON := model.MetricJSON{
		ID:    name,
		MType: metric.MType,
	}
	switch metric.MType {
	case model.TypeCounter:
		metricJSON.Delta = &metric.Delta
		metricJSON.Value = nil
	case model.TypeGauge:
		metricJSON.Value = &metric.Value
		metricJSON.Delta = nil
	}
	return &metricJSON
}

func (c *Controller) FromJSON(json *model.MetricJSON) (string, *model.Metric) {
	var metric model.Metric
	metric.MType = json.MType
	switch json.MType {
	case model.TypeCounter:
		metric.Delta = *json.Delta
		metric.Value = 0
	case model.TypeGauge:
		metric.Delta = 0
		metric.Value = *json.Value
	}
	return json.ID, &metric
}

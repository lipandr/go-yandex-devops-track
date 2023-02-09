package controller

import (
	"context"
	"errors"

	"github.com/lipandr/go-yandex-devops-track/internal/pkg/model"
	"github.com/lipandr/go-yandex-devops-track/internal/server/storage"
)

var ErrNotFound = errors.New("not found")

type metricsRepository interface {
	Get(ctx context.Context, metricType string, name string) (string, error)
	Put(ctx context.Context, metricType string, name string, value string) error
	GetAll(ctx context.Context) []model.Metric
}

type Controller struct {
	repo metricsRepository
}

func New(repository metricsRepository) *Controller {
	return &Controller{
		repo: repository,
	}
}

func (c *Controller) Get(ctx context.Context, metricType string, name string) (string, error) {
	res, err := c.repo.Get(ctx, metricType, name)
	if err != nil && errors.Is(err, storage.ErrNotFound) {
		return "", ErrNotFound
	}
	return res, nil
}

func (c *Controller) Put(ctx context.Context, metricType string, name string, value string) error {
	return c.repo.Put(ctx, name, metricType, value)
}

func (c *Controller) GetAll(ctx context.Context) []model.Metric {
	return c.repo.GetAll(ctx)
}

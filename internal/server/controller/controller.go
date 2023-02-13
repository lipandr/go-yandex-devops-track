package controller

import (
	"context"
	"errors"

	"github.com/lipandr/go-yandex-devops-track/internal/pkg/model"
	"github.com/lipandr/go-yandex-devops-track/internal/server/storage"
)

var ErrNotFound = errors.New("not found")

type metricsRepository interface {
	Get(ctx context.Context, metric *model.Metric) (string, error)
	Put(ctx context.Context, metric *model.Metric) error
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

func (c *Controller) Get(ctx context.Context, metric *model.Metric) (string, error) {
	res, err := c.repo.Get(ctx, metric)
	if err != nil && errors.Is(err, storage.ErrNotFound) {
		return "", ErrNotFound
	}
	return res, nil
}

func (c *Controller) Put(ctx context.Context, metric *model.Metric) error {
	return c.repo.Put(ctx, metric)
}

func (c *Controller) GetAll(ctx context.Context) []model.Metric {
	return c.repo.GetAll(ctx)
}

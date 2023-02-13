package memory

import (
	"context"
	"fmt"
	"sync"

	"github.com/lipandr/go-yandex-devops-track/internal/pkg/model"
	"github.com/lipandr/go-yandex-devops-track/internal/server/storage"
)

// Repository defines a memory metrics data repository.
type Repository struct {
	data model.MetricData
	sync.RWMutex
}

// New creates a new memory repository.
func New() *Repository {
	var memory model.MetricData
	memory.Data = make(map[string]*model.Metric)
	memory.Data["PollCount"] = &model.Metric{
		ID:    "PollCount",
		MType: model.TypeCounter,
		Delta: 0,
	}
	return &Repository{data: memory}
}

// Get retrieves metric value by name.
func (r *Repository) Get(_ context.Context, metric *model.Metric) (string, error) {
	r.RLock()
	defer r.RUnlock()

	if res, ok := r.data.Data[metric.ID]; ok {
		switch metric.MType {
		case model.TypeCounter:
			return fmt.Sprintf("%v", res.Delta), nil
		case model.TypeGauge:
			return fmt.Sprintf("%v", res.Value), nil
		}
	}
	return "", storage.ErrNotFound
}

// Put adds metric metadata for a given name.
func (r *Repository) Put(_ context.Context, metric *model.Metric) error {
	r.Lock()
	defer r.Unlock()
	if metric.ID == "PollCount" {
		cur := r.data.Data[metric.ID].Delta
		r.data.Data[metric.ID].Delta = cur + metric.Delta
		return nil
	}

	r.data.Data[metric.ID] = metric
	return nil
}

// GetAll retrieves all metrics.
func (r *Repository) GetAll(_ context.Context) []model.Metric {
	r.RLock()
	defer r.RUnlock()

	var data []model.Metric
	for _, v := range r.data.Data {
		tmp := model.Metric{
			ID:    v.ID,
			Value: v.Value,
		}
		if v.MType == model.TypeCounter {
			tmp.Value = float64(v.Delta)
		}
		data = append(data, tmp)
	}
	return data
}

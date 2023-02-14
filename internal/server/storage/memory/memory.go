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
}

// New creates a new memory repository.
func New() *Repository {
	var memory model.MetricData
	memory.Data = make(map[string]*model.Metric)
	memory.MU = &sync.RWMutex{}
	return &Repository{data: memory}
}

// Get retrieves metric value by name.
func (r *Repository) Get(_ context.Context, metric *model.Metric) (string, error) {
	r.data.MU.RLock()
	defer r.data.MU.RUnlock()

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
	r.data.MU.Lock()
	defer r.data.MU.Unlock()
	if metric.ID == "PollCount" {
		if _, ok := r.data.Data["PollCount"]; !ok {
			r.data.Data["PollCount"] = &model.Metric{
				ID:    "PollCount",
				MType: model.TypeCounter,
				Delta: metric.Delta,
			}
			return nil
		}
		r.data.Data[metric.ID].Delta += metric.Delta
		return nil
	}

	r.data.Data[metric.ID] = metric
	return nil
}

// GetAll retrieves all metrics.
func (r *Repository) GetAll(_ context.Context) []model.Metric {
	r.data.MU.RLock()
	defer r.data.MU.RUnlock()

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

package memory

import (
	"context"
	"fmt"
	"strconv"
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
	var memoryRepo model.MetricData
	memoryRepo.Data = make(map[string]map[string]interface{})
	memoryRepo.Data[model.TypeCounter] = make(map[string]interface{})
	memoryRepo.Data[model.TypeGauge] = make(map[string]interface{})
	return &Repository{data: memoryRepo}
}

// Get retrieves metric value by name.
func (r *Repository) Get(_ context.Context, metricType string, name string) (string, error) {
	r.RLock()
	defer r.RUnlock()
	if _, ok := r.data.Data[metricType]; ok {
		if res, ok := r.data.Data[metricType][name]; ok {
			return fmt.Sprintf("%v", res), nil
		}
	}
	return "", storage.ErrNotFound
}

// Put adds metric metadata for a given name.
func (r *Repository) Put(_ context.Context, metricType string, name string, value string) error {
	r.Lock()
	defer r.Unlock()

	var val interface{}
	switch metricType {
	case model.TypeGauge:
		fl, err := strconv.ParseFloat(value, 64)
		if err != nil {
			return err
		}
		val = fl

	case model.TypeCounter:
		in, err := strconv.Atoi(value)
		if err != nil {
			return err
		}
		val = in
	}
	r.data.Data[metricType][name] = val
	return nil
}

func (r *Repository) GetAll(_ context.Context) []model.Metric {
	r.RLock()
	defer r.RUnlock()

	var data []model.Metric
	for _, v := range r.data.Data {
		for nk, nv := range v {
			tmp := model.Metric{
				Name:  nk,
				Value: nv,
			}
			data = append(data, tmp)
		}
	}
	return data
}

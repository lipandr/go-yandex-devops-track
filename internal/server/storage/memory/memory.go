package memory

import (
	"context"
	"fmt"
	"sync"

	"github.com/lipandr/go-yandex-devops-track/internal/pkg/model"
	"github.com/lipandr/go-yandex-devops-track/internal/server/storage"
)

// Memory defines the in-memory metrics data repository.
type Memory struct {
	data model.MetricData
}

// New creates a new memory repository.
func New() *Memory {
	var memory model.MetricData
	memory.Data = make(map[string]*model.Metric)
	memory.MU = &sync.RWMutex{}
	return &Memory{data: memory}
}

// Get retrieves metric value by name.
func (r *Memory) Get(_ context.Context, name string) (string, error) {
	r.data.MU.RLock()
	defer r.data.MU.RUnlock()

	if res, ok := r.data.Data[name]; ok {
		val := ""
		switch r.data.Data[name].MType {
		case model.TypeCounter:
			val = fmt.Sprintf("%v", res.Delta)
		case model.TypeGauge:
			val = fmt.Sprintf("%v", res.Value)
		}
		if len(val) > 0 {
			return val, nil
		}
	}
	return "", storage.ErrNotFound
}

// Put adds metric metadata for a given name.
func (r *Memory) Put(_ context.Context, name string, metric *model.Metric) error {
	r.data.MU.Lock()
	defer r.data.MU.Unlock()

	switch metric.MType {
	case model.TypeCounter:
		// TODO: this code block should be moved somewhere else.
		if _, ok := r.data.Data[name]; !ok {
			r.data.Data[name] = &model.Metric{
				MType: model.TypeCounter,
				Delta: metric.Delta,
			}
		} else {
			r.data.Data[name].Delta += metric.Delta
		}
	case model.TypeGauge:
		r.data.Data[name] = metric
	}
	return nil
}

// GetAll retrieves all metrics for web UI.
func (r *Memory) GetAll(_ context.Context) []model.MetricWeb {
	r.data.MU.RLock()
	defer r.data.MU.RUnlock()

	var res []model.MetricWeb
	for k, v := range r.data.Data {
		tmp := model.MetricWeb{
			ID:    k,
			Value: v.Value,
		}
		// формируем метрику в формате ID и Value для последующего отображения в браузере
		// значения типа counter преобразовываем в значение Value float64.
		if v.MType == model.TypeCounter {
			tmp.Value = float64(v.Delta)
		}
		res = append(res, tmp)
	}
	return res
}

func (r *Memory) GetAllJSON(_ context.Context) []model.MetricJSON {
	r.data.MU.RLock()
	defer r.data.MU.RUnlock()

	var res []model.MetricJSON
	for k, v := range r.data.Data {
		tmp := model.MetricJSON{
			ID:    k,
			MType: v.MType,
		}
		switch v.MType {
		case model.TypeCounter:
			tmp.Delta = &v.Delta
			tmp.Value = nil
		case model.TypeGauge:
			tmp.Delta = nil
			tmp.Value = &v.Value
		default:
			tmp.Delta = nil
			tmp.Value = nil
		}
		res = append(res, tmp)
	}
	return res
}

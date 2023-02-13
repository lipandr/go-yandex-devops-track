package collector

import (
	"sync"
	"testing"

	"github.com/lipandr/go-yandex-devops-track/internal/pkg/model"
)

const metricCount = 29

func TestCollector_ShareMetrics(t *testing.T) {
	type fields struct {
		collector *model.MetricData
	}
	col := fields{
		collector: &model.MetricData{
			Data: make(map[string]*model.Metric),
			MU:   &sync.RWMutex{},
		},
	}
	col.collector.Data["PollCount"] = &model.Metric{
		ID:    "PollCount",
		MType: model.TypeCounter,
		Delta: 0,
	}
	tests := []struct {
		name   string
		fields fields
		want   int
	}{
		{
			name:   "check length metrics array and metrics ids to share",
			fields: col,
			want:   metricCount,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := &Collector{
				collector: tt.fields.collector,
			}

			c.UpdateMetrics()
			for _, n := range model.MetricNames {
				if _, ok := c.collector.Data[n]; !ok {
					t.Errorf("%s: %s not found", tt.name, n)
				}
			}
			got := len(c.ShareMetrics())
			if got != tt.want {
				t.Errorf("ShareMetrics() length = %v, want %v", got, tt.want)
			}
		})
	}
}

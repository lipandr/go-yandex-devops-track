package memory

import (
	"context"
	"strconv"
	"sync"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/lipandr/go-yandex-devops-track/internal/pkg/model"
	"github.com/lipandr/go-yandex-devops-track/internal/server/storage"
)

var (
	counterMetric = model.Metric{
		MType: model.TypeCounter,
		Delta: 10,
	}
	gaugeMetric = model.Metric{
		MType: model.TypeGauge,
		Value: 5.5,
	}
)

func TestNew(t *testing.T) {
	memory := New()
	require.NotNil(t, memory.data.Data, "error creating memory repository")
	require.NotNil(t, memory.data.MU, "error creating memory repository")

	mux := &sync.RWMutex{}
	require.Equal(t, mux, memory.data.MU, "the mutex is not properly initialized")
}

func TestMemory(t *testing.T) {
	t.Helper()

	data := New()
	ctx := context.Background()
	_, err := data.Get(ctx, "counterTest")
	require.Error(t, storage.ErrNotFound, err)

	testPut(t, data)
	testGet(t, data)
	testGetAll(t, data)
	testGetAllJSON(t, data)

}

func testPut(t *testing.T, data *Memory) {
	ctx := context.Background()
	err := data.Put(ctx, "counterTest", &counterMetric)
	require.NoError(t, err)
	err = data.Put(ctx, "gaugeTest", &gaugeMetric)
	require.NoError(t, err)
}

func testGet(t *testing.T, data *Memory) {
	ctx := context.Background()
	res, err := data.Get(ctx, "counterTest")
	require.NoError(t, err)
	require.Equal(t, strconv.Itoa(int(counterMetric.Delta)), res)

	res, err = data.Get(ctx, "gaugeTest")
	require.NoError(t, err)
	require.Equal(t, strconv.FormatFloat(gaugeMetric.Value, 'f', -1, 64), res)

	_, err = data.Get(ctx, "unknown")
	require.Error(t, storage.ErrNotFound, err)
	require.Equal(t, storage.ErrNotFound, err)
}

func testGetAll(t *testing.T, data *Memory) {
	ctx := context.Background()
	res := data.GetAll(ctx)
	require.Equal(t, 2, len(res))
	require.NotNil(t, res)
}

func testGetAllJSON(t *testing.T, data *Memory) {
	ctx := context.Background()
	res := data.GetAll(ctx)
	require.Equal(t, 2, len(res))
	require.NotNil(t, res)
}

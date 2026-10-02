package memstorage

import (
	"context"
	"math"
	"sync"
	"testing"

	models "github.com/nikitavaulin/metrics/internal/model"
	"github.com/stretchr/testify/require"
)

func TestUpdateMetrics(t *testing.T) {
	ctx := context.Background()
	s := New()
	delta, value := int64(3), 1.5
	metrics := []models.Metrics{
		{ID: "counter", MType: models.Counter, Delta: &delta},
		{ID: "gauge", MType: models.Gauge, Value: &value},
	}
	require.NoError(t, s.UpdateMetrics(ctx, nil))
	require.NoError(t, s.UpdateMetrics(ctx, metrics))
	delta, value = -1, 0
	require.NoError(t, s.UpdateMetrics(ctx, append(metrics, metrics[0])))
	got, err := s.GetList(ctx)
	require.NoError(t, err)
	require.Equal(t, map[string]any{"counter": int64(1), "gauge": float64(0)}, got)
	// Values are copied, so changing input pointers cannot change stored data.
	delta, value = 20, 20
	after, err := s.GetList(ctx)
	require.NoError(t, err)
	require.Equal(t, got, after)
}

func TestUpdateMetricsAtomicErrors(t *testing.T) {
	delta, value := int64(1), 2.0
	maxDelta, minDelta := int64(math.MaxInt64), int64(math.MinInt64)
	for _, tc := range []struct {
		name    string
		initial any
		metric  models.Metrics
	}{
		{"type conflict", int64(1), models.Metrics{ID: "existing", MType: models.Gauge, Value: &value}},
		{"overflow", int64(1), models.Metrics{ID: "existing", MType: models.Counter, Delta: &maxDelta}},
		{"underflow", int64(-1), models.Metrics{ID: "existing", MType: models.Counter, Delta: &minDelta}},
	} {
		t.Run(tc.name, func(t *testing.T) {
			s := New()
			require.NoError(t, s.Add(context.Background(), "existing", tc.initial))
			require.Error(t, s.UpdateMetrics(context.Background(), []models.Metrics{
				{ID: "new", MType: models.Counter, Delta: &delta}, tc.metric,
			}))
			got, err := s.GetList(context.Background())
			require.NoError(t, err)
			require.Equal(t, map[string]any{"existing": tc.initial}, got)
		})
	}
	t.Run("cancelled context", func(t *testing.T) {
		s := New()
		ctx, cancel := context.WithCancel(context.Background())
		cancel()
		require.ErrorIs(t, s.UpdateMetrics(ctx, []models.Metrics{
			{ID: "new", MType: models.Counter, Delta: &delta},
		}), context.Canceled)
		got, err := s.GetList(context.Background())
		require.NoError(t, err)
		require.Empty(t, got)
	})
}

func TestUpdateMetricsConcurrentCounters(t *testing.T) {
	s := New()
	delta := int64(1)
	metrics := []models.Metrics{{ID: "counter", MType: models.Counter, Delta: &delta}}
	var wg sync.WaitGroup
	for range 50 {
		wg.Add(1)
		go func() {
			defer wg.Done()
			if err := s.UpdateMetrics(context.Background(), metrics); err != nil {
				t.Error(err)
			}
		}()
	}
	wg.Wait()
	got, err := s.Get(context.Background(), "counter")
	require.NoError(t, err)
	require.Equal(t, int64(50), got)
}

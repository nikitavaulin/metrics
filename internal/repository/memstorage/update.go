package memstorage

import (
	"context"
	"fmt"
	"math"

	models "github.com/nikitavaulin/metrics/internal/model"
)

func (s *MemStorage) UpdateMetrics(ctx context.Context, metrics []models.Metrics) error {
	if len(metrics) == 0 {
		return nil
	}
	s.mu.Lock()
	defer s.mu.Unlock()

	updates := make(map[string]any, len(metrics))
	for _, metric := range metrics {
		if err := ctx.Err(); err != nil {
			return err
		}

		previous, exists := updates[metric.ID]
		if !exists {
			previous, exists = s.storage[metric.ID]
		}

		switch metric.MType {
		case models.Counter:
			var current int64
			if exists {
				var ok bool
				current, ok = previous.(int64)
				if !ok {
					return fmt.Errorf("metric %q already exists with a different type", metric.ID)
				}
			}

			delta := *metric.Delta
			if (delta > 0 && current > math.MaxInt64-delta) ||
				(delta < 0 && current < math.MinInt64-delta) {
				return fmt.Errorf("counter %q overflows int64", metric.ID)
			}
			updates[metric.ID] = current + delta

		case models.Gauge:
			if exists {
				if _, ok := previous.(float64); !ok {
					return fmt.Errorf("metric %q already exists with a different type", metric.ID)
				}
			}
			updates[metric.ID] = *metric.Value
		}
	}
	if err := ctx.Err(); err != nil {
		return err
	}
	if s.storage == nil {
		s.storage = make(map[string]any)
	}
	for name, value := range updates {
		s.storage[name] = value
	}
	return nil
}

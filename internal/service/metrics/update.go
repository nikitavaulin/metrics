package metricsservice

import (
	"context"

	models "github.com/nikitavaulin/metrics/internal/model"
)

func (s *MetricsService) UpdateMetrics(ctx context.Context, metrics []models.Metrics) error {
	if len(metrics) == 0 {
		return nil
	}
	return s.storage.UpdateMetrics(ctx, metrics)
}

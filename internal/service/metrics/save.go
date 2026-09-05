package metricsservice

import (
	"context"
	"fmt"
	"time"

	"github.com/nikitavaulin/metrics/internal/logger"
	models "github.com/nikitavaulin/metrics/internal/model"
	"go.uber.org/zap"
)

func (s *MetricsService) SaveToFile(ctx context.Context, filename string, intervalSec int) {
	if intervalSec <= 0 {
		logger.Log.Error(
			"invalid metrics save interval",
			zap.Int("interval_sec", intervalSec),
		)
		return
	}

	interval := time.Duration(intervalSec) * time.Second

	go func() {
		ticker := time.NewTicker(interval)
		defer ticker.Stop()
		for {
			select {
			case <-ticker.C:
				if err := s.saveToFile(filename); err != nil {
					logger.Log.Error("failed to save metrics", zap.Error(err))
					return
				}
				logger.Log.Info("current metrics have been saved to file", zap.String("file", filename))
			case <-ctx.Done():
				return
			}
		}
	}()
}

func (s *MetricsService) saveToFile(filename string) error {
	m, err := s.GetList()
	if err != nil {
		return fmt.Errorf("failed to get list of metrics: %w", err)
	}
	metrics, err := s.toModelMetricsSlice(m)
	if err != nil {
		return fmt.Errorf("failed to convert map of metrics to slice: %w", err)
	}

	if err := s.fs.SaveToJSON(filename, &metrics); err != nil {
		return fmt.Errorf("failed to save metrics to json file: %w", err)
	}

	return nil
}

func (s *MetricsService) toModelMetricsSlice(m map[string]any) ([]models.Metrics, error) {
	var result []models.Metrics
	for name, value := range m {
		var metric models.Metrics
		metric.ID = name
		switch val := value.(type) {
		case float64:
			metric.MType = models.Gauge
			metric.Value = &val
		case int64:
			metric.MType = models.Counter
			metric.Delta = &val
		default:
			return nil, fmt.Errorf("unknown mtype")
		}
		result = append(result, metric)
	}
	return result, nil
}

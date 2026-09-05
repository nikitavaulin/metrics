package metricsservice

import (
	"github.com/nikitavaulin/metrics/internal/logger"
	models "github.com/nikitavaulin/metrics/internal/model"
	"go.uber.org/zap"
)

func (s *MetricsService) LoadMetrics(filename string) {
	var metrics []models.Metrics
	if err := s.fs.LoadFromJSON(filename, &metrics); err != nil {
		logger.Log.Error("failed to load metrics from json file", zap.Error(err))
		return
	}

	for _, m := range metrics {
		if err := s.Add(m.ID, m); err != nil {
			logger.Log.Error(
				"failed to load metric from file",
				zap.String("metric", m.ID),
				zap.Error(err),
			)
		} else {
			logger.Log.Info("metric has been loaded", zap.String("metric", m.ID))
		}
	}

	logger.Log.Info("Successful metrics loading")
}

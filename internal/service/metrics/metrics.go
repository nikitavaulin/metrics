package metricsservice

import (
	"fmt"

	models "github.com/nikitavaulin/metrics/internal/model"
	"github.com/nikitavaulin/metrics/internal/repository"
)

type MetricsService struct {
	storage repository.Storage
}

func New(storage repository.Storage) *MetricsService {
	return &MetricsService{
		storage: storage,
	}
}

func (s *MetricsService) Add(name string, metric models.Metrics) error {
	prev, err := s.storage.GetOrNil(name)
	if err != nil {
		return fmt.Errorf("failed to get previous metric value: %w", err)
	}
	if prev == nil {
		return s.storage.Add(name, metric)
	}

	metric.ID = prev.ID
	// непонятно, что делать с хэшем

	switch metric.MType {
	case models.Gauge:
		return s.storage.Update(name, metric)

	case models.Counter:
		cnter := *prev.Delta + *metric.Delta
		metric.Delta = &cnter
		return s.storage.Update(name, metric)
	}
	return fmt.Errorf("unexpected error")
}

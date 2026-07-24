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

	// новое значение
	if prev == nil {
		if metric.MType == models.Counter {
			return s.storage.Add(name, *metric.Delta)
		}
		return s.storage.Add(name, *metric.Value)
	}

	// непонятно, что делать с ID и хэшем

	switch metric.MType {
	case models.Gauge:
		return s.storage.Update(name, *metric.Value)

	case models.Counter:
		prevValue, ok := prev.(int64)
		if !ok {
			return fmt.Errorf("failed to convert previous metric value to int")
		}
		counter := prevValue + *metric.Delta
		return s.storage.Update(name, counter)
	}

	return fmt.Errorf("metricsservice: unexpected error")
}

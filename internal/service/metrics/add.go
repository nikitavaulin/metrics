package metricsservice

import (
	"fmt"

	models "github.com/nikitavaulin/metrics/internal/model"
)

func (s *MetricsService) Add(name string, metric models.Metrics) error {
	if err := validateMetric(name, metric); err != nil {
		return fmt.Errorf("invalid input params: %w", err)
	}

	prev, err := s.storage.GetOrNil(name)
	if err != nil {
		return fmt.Errorf("failed to get previous metric value: %w", err)
	}

	var updateErr error

	// новое значение
	if prev == nil {
		if metric.MType == models.Counter {
			updateErr = s.storage.Add(name, *metric.Delta)
		} else {
			updateErr = s.storage.Add(name, *metric.Value)
		}
		if updateErr != nil {
			return updateErr
		}
		s.notifySave()
		return nil
	}

	switch metric.MType {
	case models.Gauge:
		updateErr = s.storage.Update(name, *metric.Value)

	case models.Counter:
		prevValue, ok := prev.(int64)
		if !ok {
			return fmt.Errorf("failed to convert previous metric value to int")
		}
		counter := prevValue + *metric.Delta
		updateErr = s.storage.Update(name, counter)
	default:
		return fmt.Errorf("metricsservice: unknown mtype")
	}
	if updateErr != nil {
		return updateErr
	}
	s.notifySave()
	return nil
}

func validateMetric(name string, metric models.Metrics) error {
	if len(name) == 0 {
		return fmt.Errorf("metric name must be not empty")
	}
	switch metric.MType {
	case models.Counter:
		if metric.Delta == nil {
			return fmt.Errorf("counter has empty delta value")
		}
	case models.Gauge:
		if metric.Value == nil {
			return fmt.Errorf("gauge has empty value")
		}
	default:
		return fmt.Errorf("unknown metric type")
	}
	return nil
}

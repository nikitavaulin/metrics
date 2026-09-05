package metricsservice

import (
	"fmt"

	models "github.com/nikitavaulin/metrics/internal/model"
)

func (s *MetricsService) Get(name string) (any, error) {
	return s.storage.Get(name)
}

func (s *MetricsService) GetList() (map[string]any, error) {
	return s.storage.GetList()
}

func (s *MetricsService) SetMValueFromStorage(m *models.Metrics) error {
	value, err := s.Get(m.ID)
	if err != nil {
		return err
	}

	switch m.MType {
	case models.Counter:
		delta, ok := value.(int64)
		if !ok {
			return fmt.Errorf("failed to convert mvalue to int64")
		}
		m.Delta = &delta

	case models.Gauge:
		val, ok := value.(float64)
		if !ok {
			return fmt.Errorf("failed to convert mvalue to float64")
		}
		m.Value = &val
	default:
		return fmt.Errorf("unknown mtype")
	}

	return nil
}

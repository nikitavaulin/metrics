package memstorage

import (
	"fmt"
	"sync"

	models "github.com/nikitavaulin/metrics/internal/model"
)

type MemStorage struct {
	storage sync.Map
}

func New() *MemStorage {
	return &MemStorage{}
}

func (s *MemStorage) Get(name string) (*models.Metrics, error) {
	metric, err := s.GetOrNil(name)
	if err != nil {
		return nil, err
	}
	if metric == nil {
		return nil, fmt.Errorf("metric with name=%s not found", name)
	}
	return metric, nil
}

func (s *MemStorage) GetOrNil(name string) (*models.Metrics, error) {
	m, ok := s.storage.Load(name)
	if !ok {
		return nil, nil
	}
	metric, ok := m.(models.Metrics)
	if !ok {
		return nil, fmt.Errorf("failed to convert object to metric")
	}
	return &metric, nil
}

func (s *MemStorage) Add(name string, metric models.Metrics) error {
	_, ok := s.storage.Load(name)
	if ok {
		return fmt.Errorf("metric is already exist")
	}
	s.storage.Store(name, metric)
	return nil
}

func (s *MemStorage) Update(name string, updated models.Metrics) error {
	_, ok := s.storage.Load(name)
	if !ok {
		return fmt.Errorf("metric with name=%s not found", name)
	}
	s.storage.Store(name, updated)
	return nil
}

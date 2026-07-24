package memstorage

import (
	"fmt"
	"sync"
)

type MemStorage struct {
	storage sync.Map
}

func New() *MemStorage {
	return &MemStorage{}
}

func (s *MemStorage) Get(name string) (any, error) {
	metric, err := s.GetOrNil(name)
	if err != nil {
		return nil, err
	}
	if metric == nil {
		return nil, fmt.Errorf("metric with name=%s not found", name)
	}
	return metric, nil
}

func (s *MemStorage) GetOrNil(name string) (any, error) {
	m, ok := s.storage.Load(name)
	if !ok {
		return nil, nil
	}
	return m, nil
}

func (s *MemStorage) Add(name string, metric any) error {
	_, ok := s.storage.Load(name)
	if ok {
		return fmt.Errorf("metric is already exist")
	}
	s.storage.Store(name, metric)
	return nil
}

func (s *MemStorage) Update(name string, updated any) error {
	_, ok := s.storage.Load(name)
	if !ok {
		return fmt.Errorf("metric with name=%s not found", name)
	}
	s.storage.Store(name, updated)
	return nil
}

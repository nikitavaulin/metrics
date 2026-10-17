package memstorage

import (
	"context"
	"fmt"
	"sync"

	"github.com/nikitavaulin/metrics/internal/commonerrors"
)

type MemStorage struct {
	mu      sync.RWMutex
	storage map[string]any
}

func New() *MemStorage {
	return &MemStorage{storage: make(map[string]any)}
}

func (s *MemStorage) PingDB(context.Context) error {
	return nil
}

func (s *MemStorage) Get(ctx context.Context, name string) (any, error) {
	metric, err := s.GetOrNil(ctx, name)
	if err != nil {
		return nil, err
	}
	if metric == nil {
		return nil, fmt.Errorf("metric with name=%s: %w", name, commonerrors.ErrNotFound)
	}
	return metric, nil
}

func (s *MemStorage) GetList(ctx context.Context) (map[string]any, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	metrics := make(map[string]any, len(s.storage))
	for name, value := range s.storage {
		metrics[name] = value
	}
	return metrics, nil
}

func (s *MemStorage) GetOrNil(ctx context.Context, name string) (any, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	m, ok := s.storage[name]
	if !ok {
		return nil, nil
	}
	return m, nil
}

func (s *MemStorage) Add(ctx context.Context, name string, metric any) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	_, ok := s.storage[name]
	if ok {
		return fmt.Errorf("metric is already exist")
	}
	if s.storage == nil {
		s.storage = make(map[string]any)
	}
	s.storage[name] = metric
	return nil
}

func (s *MemStorage) Update(ctx context.Context, name string, updated any) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	_, ok := s.storage[name]
	if !ok {
		return fmt.Errorf("metric with name=%s not found", name)
	}
	s.storage[name] = updated
	return nil
}

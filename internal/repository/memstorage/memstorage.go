package memstorage

import (
	"context"
	"errors"
	"fmt"
	"sync"

	"github.com/nikitavaulin/metrics/internal/commonerrors"
)

type MemStorage struct {
	storage sync.Map
}

func New() *MemStorage {
	return &MemStorage{}
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
	metrics := make(map[string]any)
	var err error
	s.storage.Range(func(key, value any) bool {
		name, ok := key.(string)
		if !ok {
			err = errors.New("failed to convert key to string")
			return false
		}
		metrics[name] = value
		return true
	})
	return metrics, err
}

func (s *MemStorage) GetOrNil(ctx context.Context, name string) (any, error) {
	m, ok := s.storage.Load(name)
	if !ok {
		return nil, nil
	}
	return m, nil
}

func (s *MemStorage) Add(ctx context.Context, name string, metric any) error {
	_, ok := s.storage.Load(name)
	if ok {
		return fmt.Errorf("metric is already exist")
	}
	s.storage.Store(name, metric)
	return nil
}

func (s *MemStorage) Update(ctx context.Context, name string, updated any) error {
	_, ok := s.storage.Load(name)
	if !ok {
		return fmt.Errorf("metric with name=%s not found", name)
	}
	s.storage.Store(name, updated)
	return nil
}

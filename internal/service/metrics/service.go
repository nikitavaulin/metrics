package metricsservice

import "github.com/nikitavaulin/metrics/internal/repository"

type MetricsService struct {
	storage repository.Storage
	fs      repository.FileStorage
}

func New(storage repository.Storage, fs repository.FileStorage) *MetricsService {
	return &MetricsService{
		storage: storage,
		fs:      fs,
	}
}

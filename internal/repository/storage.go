package repository

import models "github.com/nikitavaulin/metrics/internal/model"

type Storage interface {
	Add(name string, metric models.Metrics) error
	Get(name string) (*models.Metrics, error)
	GetOrNil(name string) (*models.Metrics, error)
	Update(name string, updated models.Metrics) error
}

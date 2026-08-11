package service

import models "github.com/nikitavaulin/metrics/internal/model"

type MetricsService interface {
	Add(name string, metric models.Metrics) error
	Get(name string) (any, error)
	GetList() (map[string]any, error)
}

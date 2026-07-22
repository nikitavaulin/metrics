package service

import models "github.com/nikitavaulin/metrics/internal/model"

type MetricsServie interface {
	Add(name string, metric models.Metrics) error
}

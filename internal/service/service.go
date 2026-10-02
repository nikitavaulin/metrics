package service

import (
	"context"

	models "github.com/nikitavaulin/metrics/internal/model"
)

type MetricsService interface {
	Add(ctx context.Context, name string, metric models.Metrics) error
	UpdateMetrics(ctx context.Context, metrics []models.Metrics) error
	Get(ctx context.Context, name string) (any, error)
	GetList(ctx context.Context) (map[string]any, error)
	SetMValueFromStorage(ctx context.Context, m *models.Metrics) error
}

type PingService interface {
	PingDB(ctx context.Context) error
}

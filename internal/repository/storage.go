package repository

import (
	"context"

	models "github.com/nikitavaulin/metrics/internal/model"
)

type Storage interface {
	PingDB(ctx context.Context) error
	Add(ctx context.Context, name string, metric any) error
	UpdateMetrics(ctx context.Context, metrics []models.Metrics) error
	Get(ctx context.Context, name string) (any, error)
	GetList(ctx context.Context) (map[string]any, error)
	GetOrNil(ctx context.Context, name string) (any, error)
	Update(ctx context.Context, name string, updated any) error
}

type FileStorage interface {
	SaveToJSON(filename string, data any) error
	LoadFromJSON(filename string, dest any) error
}

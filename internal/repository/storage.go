package repository

import "context"

type Storage interface {
	Add(ctx context.Context, name string, metric any) error
	Get(ctx context.Context, name string) (any, error)
	GetList(ctx context.Context) (map[string]any, error)
	GetOrNil(ctx context.Context, name string) (any, error)
	Update(ctx context.Context, name string, updated any) error
}

type FileStorage interface {
	SaveToJSON(filename string, data any) error
	LoadFromJSON(filename string, dest any) error
}

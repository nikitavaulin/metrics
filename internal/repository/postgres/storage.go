package postgresstorage

import (
	"context"
	"fmt"

	dbconfig "github.com/nikitavaulin/metrics/internal/config/db"
)

type Storage struct {
	pool *dbconfig.Pool
}

func New(pool *dbconfig.Pool) *Storage {
	return &Storage{
		pool: pool,
	}
}

func (s *Storage) PingDB(ctx context.Context) error {
	ctx, cancel := context.WithTimeout(ctx, s.pool.OperationTime())
	defer cancel()
	if err := s.pool.Ping(ctx); err != nil {
		return fmt.Errorf("failed to ping pool: %w", err)
	}
	return nil
}

package dbconfig

import (
	"context"
	"fmt"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
)

const (
	defaultOpTime time.Duration = 10 * time.Minute
)

type Pool struct {
	*pgxpool.Pool
	opTime time.Duration
}

func New(ctx context.Context, connStr string) (*Pool, error) {
	cfg, err := pgxpool.ParseConfig(connStr)
	if err != nil {
		return nil, fmt.Errorf("failed to parse pgx pool conn string:%w", err)
	}

	pool, err := pgxpool.NewWithConfig(ctx, cfg)
	if err != nil {
		return nil, fmt.Errorf("failed to create new pgx pool %w", err)
	}

	return &Pool{
		Pool:   pool,
		opTime: defaultOpTime,
	}, nil
}

func (p *Pool) OperationTime() time.Duration {
	if p.opTime == 0 {
		return defaultOpTime
	}
	return p.opTime
}

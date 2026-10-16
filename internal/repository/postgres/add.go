package postgresstorage

import (
	"context"
	"fmt"

	"github.com/nikitavaulin/metrics/internal/domain"
	models "github.com/nikitavaulin/metrics/internal/model"
	"github.com/nikitavaulin/metrics/internal/retry"
)

func (s *Storage) Add(ctx context.Context, name string, metric any) error {
	ctx, cancel := context.WithTimeout(ctx, s.pool.OperationTime())
	defer cancel()

	mtype, err := domain.GetMetricTypeByValue(metric)
	if err != nil {
		return fmt.Errorf("failed to get mtype by value: %w", err)
	}

	var query string
	switch mtype {
	case models.Counter:
		query = `INSERT INTO metrics (name, mtype, delta) VALUES ($1, $2, $3);`
	case models.Gauge:
		query = `INSERT INTO metrics (name, mtype, value) VALUES ($1, $2, $3);`
	default:
		return fmt.Errorf("unsupported metric type: %s", mtype)
	}

	err = retry.Do(ctx, func(ctx context.Context) error {
		_, err := s.pool.Exec(ctx, query, name, mtype, metric)
		return err
	}, IsRetriableError)
	if err != nil {
		return fmt.Errorf("failed to insert metric into db: %w", err)
	}

	return nil
}

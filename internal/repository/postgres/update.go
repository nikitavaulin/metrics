package postgresstorage

import (
	"context"
	"fmt"

	"github.com/nikitavaulin/metrics/internal/commonerrors"
	"github.com/nikitavaulin/metrics/internal/domain"
	models "github.com/nikitavaulin/metrics/internal/model"
	"github.com/nikitavaulin/metrics/internal/retry"
)

func (s *Storage) Update(ctx context.Context, name string, updated any) error {
	ctx, cancel := context.WithTimeout(ctx, s.pool.OperationTime())
	defer cancel()

	mtype, err := domain.GetMetricTypeByValue(updated)
	if err != nil {
		return fmt.Errorf("failed to get mtype by value: %w", err)
	}

	var query string
	switch mtype {
	case models.Counter:
		query = `UPDATE metrics SET delta = $1 WHERE name = $2 AND mtype = $3;`
	case models.Gauge:
		query = `UPDATE metrics SET value = $1 WHERE name = $2 AND mtype = $3;`
	default:
		return fmt.Errorf("unsupported metric type: %s", mtype)
	}

	err = retry.Do(ctx, func(ctx context.Context) error {
		result, err := s.pool.Exec(ctx, query, updated, name, mtype)
		if err != nil {
			return fmt.Errorf("failed to update metric in db: %w", err)
		}
		if result.RowsAffected() == 0 {
			return fmt.Errorf("metric with name=%s and mtype=%s: %w", name, mtype, commonerrors.ErrNotFound)
		}
		return nil
	}, IsRetriableError)
	if err != nil {
		return err
	}

	return nil
}

func (s *Storage) UpdateMetrics(ctx context.Context, metrics []models.Metrics) error {
	if len(metrics) == 0 {
		return nil
	}

	ctx, cancel := context.WithTimeout(ctx, s.pool.OperationTime())
	defer cancel()

	return retry.Do(ctx, func(ctx context.Context) error {
		return s.updateMetricsOnce(ctx, metrics)
	}, IsRetriableError)
}

func (s *Storage) updateMetricsOnce(ctx context.Context, metrics []models.Metrics) error {
	tx, err := s.pool.Begin(ctx)
	if err != nil {
		return fmt.Errorf("failed to begin tx: %w", err)
	}

	defer tx.Rollback(ctx)

	const query = `
		INSERT INTO metrics (name, mtype, delta, value)
		VALUES ($1, $2, $3, $4)
		ON CONFLICT (name) DO UPDATE SET
			delta = CASE WHEN EXCLUDED.mtype = 'counter'
				THEN metrics.delta + EXCLUDED.delta ELSE NULL END,
			value = EXCLUDED.value
		WHERE metrics.mtype = EXCLUDED.mtype;
	`

	for _, metric := range metrics {
		var delta *int64
		var value *float64
		if metric.MType == models.Counter {
			delta = metric.Delta
		} else {
			value = metric.Value
		}
		result, err := tx.Exec(ctx, query, metric.ID, metric.MType, delta, value)
		if err != nil {
			return fmt.Errorf("failed to update metric %q in db: %w", metric.ID, err)
		}
		if result.RowsAffected() == 0 {
			return fmt.Errorf("metric %q already exists with a different type", metric.ID)
		}
	}

	if err := tx.Commit(ctx); err != nil {
		// Commit may have succeeded before the connection was lost.
		// Retrying could apply counter increments twice.
		return &commitError{err: fmt.Errorf("failed to commit metrics update: %w", err)}
	}

	return nil
}

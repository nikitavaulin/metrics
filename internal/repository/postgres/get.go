package postgresstorage

import (
	"context"
	"errors"
	"fmt"
	"github.com/nikitavaulin/metrics/internal/retry"

	"github.com/jackc/pgx/v5"
	"github.com/nikitavaulin/metrics/internal/commonerrors"
	models "github.com/nikitavaulin/metrics/internal/model"
)

func (s *Storage) Get(ctx context.Context, name string) (any, error) {
	ctx, cancel := context.WithTimeout(ctx, s.pool.OperationTime())
	defer cancel()

	query := `
		SELECT name, mtype, delta, value FROM metrics WHERE name = $1;
	`

	var m models.Metrics

	err := retry.Do(ctx, func(ctx context.Context) error {
		m = models.Metrics{}
		return s.pool.QueryRow(ctx, query, name).Scan(&m.ID, &m.MType, &m.Delta, &m.Value)
	}, IsRetriableError)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, fmt.Errorf("metric with name=%s: %w", name, commonerrors.ErrNotFound)
		}
		return nil, fmt.Errorf("failed to get metric from db: %w", err)
	}

	if m.MType == models.Counter {
		return *m.Delta, nil
	}

	return *m.Value, nil
}

func (s *Storage) GetList(ctx context.Context) (map[string]any, error) {
	ctx, cancel := context.WithTimeout(ctx, s.pool.OperationTime())
	defer cancel()

	var metrics map[string]any

	query := `
		SELECT name, mtype, delta, value FROM metrics
	`

	err := retry.Do(ctx, func(ctx context.Context) error {
		attemptMetrics := make(map[string]any)
		rows, err := s.pool.Query(ctx, query)
		if err != nil {
			return fmt.Errorf("failed to get metrics from db: %w", err)
		}
		defer rows.Close()

		for rows.Next() {
			var m models.Metrics
			if err := rows.Scan(&m.ID, &m.MType, &m.Delta, &m.Value); err != nil {
				return fmt.Errorf("failed to scan metric from db: %w", err)
			}
			if m.MType == models.Counter {
				attemptMetrics[m.ID] = *m.Delta
			} else {
				attemptMetrics[m.ID] = *m.Value
			}
		}

		if err := rows.Err(); err != nil {
			return fmt.Errorf("rows err: %w", err)
		}
		metrics = attemptMetrics
		return nil
	}, IsRetriableError)
	if err != nil {
		return nil, err
	}

	return metrics, nil
}

func (s *Storage) GetOrNil(ctx context.Context, name string) (any, error) {
	value, err := s.Get(ctx, name)
	if err != nil {
		if errors.Is(err, commonerrors.ErrNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return value, nil
}

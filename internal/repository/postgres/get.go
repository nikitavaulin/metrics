package postgresstorage

import (
	"context"
	"errors"
	"fmt"

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

	row := s.pool.QueryRow(ctx, query, name)

	if err := row.Scan(&m.ID, &m.MType, &m.Delta, &m.Value); err != nil {
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

	metrics := make(map[string]any)

	query := `
		SELECT name, mtype, delta, value FROM metrics
	`

	rows, err := s.pool.Query(ctx, query)
	if err != nil {
		return nil, fmt.Errorf("failed to get metrics from db: %w", err)
	}
	defer rows.Close()

	for rows.Next() {
		var m models.Metrics
		if err := rows.Scan(&m.ID, &m.MType, &m.Delta, &m.Value); err != nil {
			return nil, fmt.Errorf("failed to scan metric from db: %w", err)
		}
		if m.MType == models.Counter {
			metrics[m.ID] = *m.Delta
		} else {
			metrics[m.ID] = *m.Value
		}
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("rows err: %w", err)
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

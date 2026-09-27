package postgresstorage

import (
	"context"
	"fmt"

	"github.com/nikitavaulin/metrics/internal/commonerrors"
	"github.com/nikitavaulin/metrics/internal/domain"
	models "github.com/nikitavaulin/metrics/internal/model"
)

func (s *Storage) Update(ctx context.Context, name string, updated any) error {
	ctx, cancel := context.WithTimeout(ctx, s.pool.OperationTime())
	defer cancel()

	mtype, err := domain.GetMetricTypeByValue(updated)
	if err != nil {
		return fmt.Errorf("failed to get mtype by value: %w", err)
	}

	var valueParamName string
	if mtype == models.Counter {
		valueParamName = "delta"
	} else {
		valueParamName = "value"
	}

	query := fmt.Sprintf(`
		UPDATE metrics SET %s = $1 WHERE name = $2 AND mtype = $3;
	`, valueParamName)

	result, err := s.pool.Exec(ctx, query, updated, name, mtype)
	if err != nil {
		return fmt.Errorf("failed to update metric in db: %w", err)
	}
	if result.RowsAffected() == 0 {
		return fmt.Errorf("metric with name=%s and mtype=%s: %w", name, mtype, commonerrors.ErrNotFound)
	}

	return nil
}

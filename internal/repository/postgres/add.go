package postgresstorage

import (
	"context"
	"fmt"

	"github.com/nikitavaulin/metrics/internal/domain"
	models "github.com/nikitavaulin/metrics/internal/model"
)

func (s *Storage) Add(ctx context.Context, name string, metric any) error {
	ctx, cancel := context.WithTimeout(ctx, s.pool.OperationTime())
	defer cancel()

	mtype, err := domain.GetMetricTypeByValue(metric)
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
		INSERT INTO metrics (name, mtype, %s)
		VALUES ($1, $2, $3);
	`, valueParamName)

	_, err = s.pool.Exec(ctx, query, name, mtype, metric)
	if err != nil {
		return fmt.Errorf("failed to insert metric into db: %w", err)
	}

	return nil
}

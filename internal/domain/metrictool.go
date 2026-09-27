package domain

import (
	"errors"

	models "github.com/nikitavaulin/metrics/internal/model"
)

func GetMetricTypeByValue(value any) (string, error) {
	switch value.(type) {
	case int64:
		return models.Counter, nil
	case float64:
		return models.Gauge, nil
	}
	return "", errors.New("unknown metric type")
}

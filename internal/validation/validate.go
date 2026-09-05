package validation

import (
	"errors"
	"fmt"
	"path/filepath"
	"strconv"
	"strings"

	models "github.com/nikitavaulin/metrics/internal/model"
)

func ValidateServerAddress(addr string) error {
	parts := strings.Split(addr, ":")
	if len(parts) != 2 {
		return fmt.Errorf("invalid address struct got: %s, want: %s", addr, "<addr>:<port>")
	}

	port, err := strconv.Atoi(parts[1])
	if err != nil {
		return fmt.Errorf("failed to convert port: %w", err)
	}

	if err := ValidateIntInBounds(port, 1, 65535); err != nil {
		return fmt.Errorf("invalid port value: %w", err)
	}

	return nil
}

func ValidateIntInBounds(value, minVal, maxVal int) error {
	if !(minVal <= value && value <= maxVal) {
		return fmt.Errorf("value '%d' not in bounds %d..%d", value, minVal, maxVal)
	}
	return nil
}

func ValidateMetricName(name string) error {
	if len(name) == 0 {
		return errors.New("metric name is empty")
	}
	return nil
}

func ValidateMetric(metric models.Metrics) error {
	if err := ValidateMetricName(metric.ID); err != nil {
		return err
	}
	switch metric.MType {
	case models.Counter:
		if metric.Delta == nil {
			return errors.New("counter metric delta is empty")
		}
	case models.Gauge:
		if metric.Value == nil {
			return errors.New("gauge metric value is empty")
		}
	default:
		return errors.New("unknown metric type")
	}
	return nil
}

func ValidateMetricType(mtype string) error {
	if mtype == "" {
		return fmt.Errorf("mtype is empty")
	}
	if mtype != models.Counter && mtype != models.Gauge {
		return fmt.Errorf("unknown metric type: %s", mtype)
	}
	return nil
}

func ValidateFileExt(fname, extWant string) error {
	ext := filepath.Ext(fname)
	if ext != extWant {
		return fmt.Errorf("missmatch of file extention: got: %s, want: %s", ext, extWant)
	}
	return nil
}

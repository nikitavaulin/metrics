package agent

import (
	"bytes"
	"compress/gzip"
	"fmt"
)

func compress(data []byte) ([]byte, error) {
	var b bytes.Buffer

	w := gzip.NewWriter(&b)

	if _, err := w.Write(data); err != nil {
		return nil, fmt.Errorf("failed write data to compress buffer: %w", err)
	}

	if err := w.Close(); err != nil {
		return nil, fmt.Errorf("failed compress data: %w", err)
	}

	return b.Bytes(), nil
}

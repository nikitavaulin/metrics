package gzipcompress

import (
	"compress/gzip"
	"fmt"
	"io"
)

type Reader struct {
	r  io.ReadCloser
	zr *gzip.Reader
}

func NewReader(r io.ReadCloser) (*Reader, error) {
	zr, err := gzip.NewReader(r)
	if err != nil {
		return nil, fmt.Errorf("failed to create gzip reader: %w", err)
	}
	return &Reader{zr: zr, r: r}, nil
}

func (r *Reader) Read(p []byte) (n int, err error) {
	return r.zr.Read(p)
}

func (r *Reader) Close() error {
	if err := r.r.Close(); err != nil {
		return fmt.Errorf("failed to close reader: %w", err)
	}
	return r.zr.Close()
}

package gzipcompress

import (
	"compress/gzip"
	"net/http"
)

type Writer struct {
	rw http.ResponseWriter
	zw *gzip.Writer
}

func NewWriter(rw http.ResponseWriter) *Writer {
	return &Writer{
		rw: rw,
		zw: gzip.NewWriter(rw),
	}
}

func (w *Writer) Write(b []byte) (n int, err error) {
	if w.shouldCompress() {
		return w.zw.Write(b)
	}
	return w.rw.Write(b)
}

func (w *Writer) Header() http.Header {
	return w.rw.Header()
}

func (w *Writer) WriteHeader(statusCode int) {
	if statusCode < 300 && w.shouldCompress() {
		w.rw.Header().Set("Content-Encoding", "gzip")
	}
	w.rw.WriteHeader(statusCode)
}

func (w *Writer) Close() error {
	return w.zw.Close()
}

func (w *Writer) shouldCompress() bool {
	contentType := w.rw.Header().Get("Content-Type")
	return contentType == "application/json" || contentType == "text/html"
}

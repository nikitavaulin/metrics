package gzipcompress

import (
	"compress/gzip"
	"net/http"
)

type Writer struct {
	rw          http.ResponseWriter
	zw          *gzip.Writer
	wroteHeader bool
}

func NewWriter(rw http.ResponseWriter) *Writer {
	return &Writer{
		rw: rw,
	}
}

func (w *Writer) Write(b []byte) (n int, err error) {
	if !w.wroteHeader {
		w.WriteHeader(http.StatusOK)
	}
	if w.zw != nil {
		return w.zw.Write(b)
	}
	return w.rw.Write(b)
}

func (w *Writer) Header() http.Header {
	return w.rw.Header()
}

func (w *Writer) WriteHeader(statusCode int) {
	if w.wroteHeader {
		return
	}
	w.wroteHeader = true
	if statusCode >= 200 && statusCode < 300 && statusCode != http.StatusNoContent && statusCode != http.StatusResetContent && w.shouldCompress() {
		w.rw.Header().Set("Content-Encoding", "gzip")
		w.zw = gzip.NewWriter(w.rw)
	}
	w.rw.WriteHeader(statusCode)
}

func (w *Writer) Close() error {
	if w.zw != nil {
		return w.zw.Close()
	}
	return nil
}

func (w *Writer) shouldCompress() bool {
	contentType := w.rw.Header().Get("Content-Type")
	return contentType == "application/json" || contentType == "text/html"
}

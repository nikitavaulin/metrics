package httpserver

import (
	"net/http"
)

type WriterResponse struct {
	http.ResponseWriter
	Status int
	Size   int
}

func NewResponseWriter(rw http.ResponseWriter) *WriterResponse {
	return &WriterResponse{
		ResponseWriter: rw,
	}
}

func (rw *WriterResponse) Write(b []byte) (int, error) {
	size, err := rw.ResponseWriter.Write(b)
	if err != nil {
		return 0, err
	}
	rw.Size += size
	return size, nil
}

func (rw *WriterResponse) WriteHeader(statusCode int) {
	rw.ResponseWriter.WriteHeader(statusCode)
	rw.Status = statusCode
}

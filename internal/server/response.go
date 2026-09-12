package httpserver

import (
	"encoding/json"
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

func JSONResponse(rw http.ResponseWriter, data any, statusCode int) {
	rw.Header().Set("Content-Type", "application/json")
	rw.WriteHeader(statusCode)
	json.NewEncoder(rw).Encode(&data)
}

func ErrorResponse(rw http.ResponseWriter, err error, msg string, statusCode int) {
	resp := map[string]string{
		"error":   err.Error(),
		"message": msg,
	}
	rw.Header().Set("Content-Type", "application/json")
	rw.WriteHeader(statusCode)
	json.NewEncoder(rw).Encode(&resp)
}

package httpserver

import (
	"net/http"
	"time"

	"github.com/nikitavaulin/metrics/internal/logger"
	"go.uber.org/zap"
)

func LogRequest(h http.Handler) http.Handler {
	logFunc := func(rw http.ResponseWriter, r *http.Request) {
		start := time.Now()

		logger.Log.Info(
			">>> Invoke HTTP request",
			zap.String("Method", r.Method),
			zap.String("URL", r.URL.Path),
		)

		writer := NewResponseWriter(rw)

		h.ServeHTTP(writer, r)
		latency := time.Since(start)

		logger.Log.Info(
			"<<< Send HTTP response",
			zap.Int("Status", writer.Status),
			zap.Int("Response size", writer.Size),
			zap.Duration("Latency", latency),
		)
	}
	return http.HandlerFunc(logFunc)
}

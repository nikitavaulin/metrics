package httpserver

import (
	"net/http"
	"strings"
	"time"

	"github.com/nikitavaulin/metrics/internal/gzipcompress"
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

func GzipResponse(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		const gzipName = "gzip"
		originalRW := w

		if strings.Contains(r.Header.Get("Accept-Encoding"), gzipName) {
			gzipRW := gzipcompress.NewWriter(w)
			originalRW = gzipRW
			defer gzipRW.Close()
		}

		next.ServeHTTP(originalRW, r)
	})
}

func GzipRequest(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Header.Get("Content-Encoding") == "gzip" {
			gzipR, err := gzipcompress.NewReader(r.Body)
			if err != nil {
				ErrorResponse(w, err, "invalid gzip request body", http.StatusBadRequest)
				return
			}
			r.Body = gzipR
			defer gzipR.Close()
		}

		next.ServeHTTP(w, r)
	})
}

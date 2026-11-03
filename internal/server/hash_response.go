package httpserver

import (
	"bytes"
	"net/http"

	"github.com/nikitavaulin/metrics/internal/security"
)

type hashResponseWriter struct {
	http.ResponseWriter
	body   bytes.Buffer
	status int
}

func (w *hashResponseWriter) WriteHeader(status int) {
	if w.status != 0 {
		return
	}
	if status >= 100 && status < 200 && status != http.StatusSwitchingProtocols {
		w.ResponseWriter.WriteHeader(status)
		return
	}
	w.status = status
}

func (w *hashResponseWriter) Write(body []byte) (int, error) {
	if w.status == 0 {
		w.WriteHeader(http.StatusOK)
	}
	if w.status == http.StatusNoContent || w.status == http.StatusNotModified || w.status < 200 {
		return 0, http.ErrBodyNotAllowed
	}
	return w.body.Write(body)
}

func HashResponse(next http.Handler, key *string) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if key == nil {
			next.ServeHTTP(w, r)
			return
		}

		buffered := &hashResponseWriter{ResponseWriter: w}
		next.ServeHTTP(buffered, r)
		w.Header().Set("HashSHA256", security.Hash(buffered.body.String(), *key))
		if buffered.status == 0 {
			buffered.status = http.StatusOK
		}
		w.WriteHeader(buffered.status)
		if buffered.body.Len() > 0 {
			_, _ = w.Write(buffered.body.Bytes())
		}
	})
}

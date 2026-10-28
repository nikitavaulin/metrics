package httpserver

import (
	"bytes"
	"crypto/hmac"
	"encoding/hex"
	"io"
	"net/http"

	"github.com/nikitavaulin/metrics/internal/security"
)

func HashRequest(next http.Handler, key *string) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		requestHash := r.Header.Get("HashSHA256")
		if key == nil {
			next.ServeHTTP(w, r)
			return
		}
		if requestHash == "" {
			http.Error(w, "missing request hash", http.StatusBadRequest)
			return
		}

		providedHash, err := hex.DecodeString(requestHash)
		if err != nil {
			http.Error(w, "invalid request hash", http.StatusBadRequest)
			return
		}
		body, err := io.ReadAll(r.Body)
		if err != nil {
			http.Error(w, "failed to read request body", http.StatusBadRequest)
			return
		}
		r.Body = struct {
			io.Reader
			io.Closer
		}{bytes.NewReader(body), r.Body}

		expectedHash, _ := hex.DecodeString(security.Hash(string(body), *key))
		if !hmac.Equal(providedHash, expectedHash) {
			http.Error(w, "request hash mismatch", http.StatusBadRequest)
			return
		}
		next.ServeHTTP(w, r)
	})
}

package httpserver

import (
	"compress/gzip"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/nikitavaulin/metrics/internal/security"
	"github.com/stretchr/testify/require"
)

func TestHashResponse(t *testing.T) {
	key, emptyKey := "secret", ""
	for _, tt := range []struct {
		name   string
		key    *string
		body   string
		status int
		gzip   bool
	}{
		{name: "no key", body: "response", status: http.StatusOK},
		{name: "signed", key: &key, body: "response", status: http.StatusCreated},
		{name: "empty key", key: &emptyKey, body: "response", status: http.StatusOK},
		{name: "empty body", key: &key, status: http.StatusNoContent},
		{name: "error", key: &key, body: "error", status: http.StatusBadRequest},
		{name: "gzip", key: &key, body: `{"value":123}`, status: http.StatusOK, gzip: true},
	} {
		t.Run(tt.name, func(t *testing.T) {
			next := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				w.Header().Set("Content-Type", "application/json")
				w.WriteHeader(tt.status)
				if tt.body != "" {
					_, _ = w.Write([]byte(tt.body[:len(tt.body)/2]))
					_, _ = w.Write([]byte(tt.body[len(tt.body)/2:]))
				}
			})
			r := httptest.NewRequest(http.MethodGet, "/", nil)
			if tt.gzip {
				r.Header.Set("Accept-Encoding", "gzip")
			}
			w := httptest.NewRecorder()
			GzipCompress(HashResponse(next, tt.key)).ServeHTTP(w, r)
			resp := w.Result()
			defer resp.Body.Close()
			require.Equal(t, tt.status, resp.StatusCode)
			var reader io.Reader = resp.Body
			if tt.gzip {
				require.Equal(t, "gzip", resp.Header.Get("Content-Encoding"))
				gz, err := gzip.NewReader(reader)
				require.NoError(t, err)
				defer gz.Close()
				reader = gz
			}
			body, err := io.ReadAll(reader)
			require.NoError(t, err)
			require.Equal(t, tt.body, string(body))
			if tt.key == nil {
				require.Empty(t, resp.Header.Get("HashSHA256"))
			} else {
				require.Equal(t, security.Hash(string(body), *tt.key), resp.Header.Get("HashSHA256"))
			}
		})
	}
}

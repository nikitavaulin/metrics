package httpserver

import (
	"bytes"
	"compress/gzip"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/nikitavaulin/metrics/internal/security"
	"github.com/stretchr/testify/require"
)

func TestHashRequest(t *testing.T) {
	key, emptyKey := "secret", ""
	const body = `{"value":123}`
	for _, tt := range []struct {
		name   string
		key    *string
		header string
		gzip   bool
		status int
	}{
		{name: "no key", header: "invalid", status: http.StatusNoContent},
		{name: "no header", key: &key, status: http.StatusNoContent},
		{name: "valid", key: &key, header: security.Hash(body, key), status: http.StatusNoContent},
		{name: "empty key", key: &emptyKey, header: security.Hash(body, emptyKey), status: http.StatusNoContent},
		{name: "wrong key", key: &key, header: security.Hash(body, "other"), status: http.StatusBadRequest},
		{name: "changed body", key: &key, header: security.Hash("other", key), status: http.StatusBadRequest},
		{name: "malformed hash", key: &key, header: "invalid", status: http.StatusBadRequest},
		{name: "gzip", key: &key, header: security.Hash(body, key), gzip: true, status: http.StatusNoContent},
	} {
		t.Run(tt.name, func(t *testing.T) {
			var input io.Reader = strings.NewReader(body)
			if tt.gzip {
				var compressed bytes.Buffer
				writer := gzip.NewWriter(&compressed)
				_, err := writer.Write([]byte(body))
				require.NoError(t, err)
				require.NoError(t, writer.Close())
				input = &compressed
			}
			r := httptest.NewRequest(http.MethodPost, "/updates/", input)
			if tt.header != "" {
				r.Header.Set("HashSHA256", tt.header)
			}
			if tt.gzip {
				r.Header.Set("Content-Encoding", "gzip")
			}
			called := false
			next := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				called = true
				got, err := io.ReadAll(r.Body)
				require.NoError(t, err)
				require.Equal(t, body, string(got))
				w.WriteHeader(http.StatusNoContent)
			})
			w := httptest.NewRecorder()
			GzipRequest(HashRequest(next, tt.key)).ServeHTTP(w, r)
			require.Equal(t, tt.status, w.Code)
			require.Equal(t, tt.status == http.StatusNoContent, called)
		})
	}
}

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
			GzipResponse(HashResponse(next, tt.key)).ServeHTTP(w, r)
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

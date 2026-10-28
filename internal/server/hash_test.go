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
		{name: "no header", key: &key, status: http.StatusBadRequest},
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
			GzipCompress(HashRequest(next, tt.key)).ServeHTTP(w, r)
			require.Equal(t, tt.status, w.Code)
			require.Equal(t, tt.status == http.StatusNoContent, called)
		})
	}
}

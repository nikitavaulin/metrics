package httpserver

import (
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	serverconfig "github.com/nikitavaulin/metrics/internal/config/server"
	"github.com/nikitavaulin/metrics/internal/security"
	"github.com/stretchr/testify/require"
)

func TestInvalidGzipResponseIsSigned(t *testing.T) {
	key := "secret"
	server := New(&serverconfig.Config{Key: &key})
	called := false
	server.mux.Post("/updates/", func(w http.ResponseWriter, r *http.Request) {
		called = true
		w.WriteHeader(http.StatusOK)
	})

	request := httptest.NewRequest(http.MethodPost, "/updates/", strings.NewReader("invalid gzip"))
	request.Header.Set("Content-Encoding", "gzip")
	request.Header.Set("Accept-Encoding", "gzip")
	request.Header.Set("HashSHA256", security.Hash("invalid gzip", key))
	writer := httptest.NewRecorder()
	server.handler().ServeHTTP(writer, request)

	response := writer.Result()
	defer response.Body.Close()
	body, err := io.ReadAll(response.Body)
	require.NoError(t, err)
	require.Equal(t, http.StatusBadRequest, response.StatusCode)
	require.False(t, called)
	require.NotEmpty(t, body)
	require.Equal(t, security.Hash(string(body), key), response.Header.Get("HashSHA256"))
}

package pinghandler

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/nikitavaulin/metrics/internal/repository/memstorage"
	pingservice "github.com/nikitavaulin/metrics/internal/service/ping"
	"github.com/stretchr/testify/require"
)

func TestPingMemoryStorage(t *testing.T) {
	handler := New(pingservice.New(memstorage.New()))
	response := httptest.NewRecorder()
	handler.PingDB(response, httptest.NewRequest(http.MethodGet, "/ping", nil))
	require.Equal(t, http.StatusOK, response.Code)
}

package pinghandler

import (
	"net/http"

	httpserver "github.com/nikitavaulin/metrics/internal/server"
)

func (h *PingHandler) PingDB(rw http.ResponseWriter, r *http.Request) {
	if err := h.pingService.PingDB(r.Context()); err != nil {
		httpserver.ErrorResponse(rw, err, "failed to ping db", http.StatusInternalServerError)
		return
	}
	rw.WriteHeader(http.StatusOK)
}

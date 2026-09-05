package httpserver

import (
	"net/http"

	"github.com/go-chi/chi/v5"
)

type Route struct {
	Method  string
	Path    string
	Handler http.HandlerFunc
}

type Router struct {
	*chi.Mux
}

func NewRouter() *Router {
	return &Router{
		Mux: chi.NewRouter(),
	}
}

func (r *Router) RegisterRoutes(routes []Route) {
	for _, route := range routes {
		r.Method(route.Method, route.Path, route.Handler)
	}
}

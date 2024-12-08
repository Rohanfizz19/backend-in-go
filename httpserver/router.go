package httpserver

import (
	"net/http"
	"net/http/pprof"
	"sync"

	"github.com/gorilla/mux"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

// Config for a single route.
// This defines the path on which route is defined
type RouteConfig struct {
	Path       string
	Handler    http.Handler
	Methods    []string //HTTP Methods
	Instrument bool     //Should we instrument metrics.
}

// Router is the interface exposed to add routes to the mux.
//
// AddRoute would add one route to the mux.
// This method can be invoked my multiple goroutines at the same time.
//
// mux returns the http Handler which implements multiplexer for various routes.
type Router interface {
	AddRoute(route RouteConfig)
	Mux() http.Handler
}

// implementation of the Router interface.
// gorilla mux is used but can be swapped by other router.
type router struct {
	gmux         *mux.Router
	instrumentor *PathInstrumentor
	sync.Mutex
}

func NewRouter() Router {
	r := &router{
		gmux:         mux.NewRouter().StrictSlash(true),
		instrumentor: NewInstrumentor(),
	}
	r.debugRoutes()
	r.promRoute()
	return r
}

// AddRoute can be invoked by multiple goroutines at the same time.
// This method adds single route to the mux
func (r *router) AddRoute(route RouteConfig) {
	r.Lock()
	defer r.Unlock()
	handler := route.Handler
	if route.Instrument {
		handler = r.instrumentor.Instrument(route.Path, handler)
	}
	rt := r.gmux.Handle(route.Path, handler)
	if route.Methods != nil && len(route.Methods) > 0 {
		rt.Methods(route.Methods...)
	}
}

func (r *router) Mux() http.Handler {
	return r.gmux
}

func (r *router) debugRoutes() {
	mux := r.gmux
	mux.HandleFunc("/debug/pprof/", pprof.Index)
	mux.HandleFunc("/debug/pprof/cmdline", pprof.Cmdline)
	mux.HandleFunc("/debug/pprof/profile", pprof.Profile)
	mux.HandleFunc("/debug/pprof/symbol", pprof.Symbol)
	mux.HandleFunc("/debug/pprof/trace", pprof.Trace)
	mux.Handle("/debug/pprof/goroutine", pprof.Handler("goroutine"))

}

func (r *router) promRoute() {
	r.gmux.Handle("/metrics", promhttp.Handler())
}

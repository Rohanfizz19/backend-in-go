package httpserver

import (
	"net/http"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

type PathInstrumentor struct {
	counter  *prometheus.CounterVec
	duration *prometheus.HistogramVec
}

func NewInstrumentor() *PathInstrumentor {

	//counter for api requests
	counter := prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Namespace: "foundation_gateway_service",
			Name:      "api_requests_total",
			Help:      "A counter for requests to the wrapped handler.",
		},
		[]string{"path", "code", "method"},
	)

	//This would have default buckets - might make sense to make that configurable.
	duration := prometheus.NewHistogramVec(
		prometheus.HistogramOpts{
			Namespace: "foundation_gateway_service",
			Name:      "request_duration_seconds",
			Help:      "A histogram of latencies for requests.",
		},
		[]string{"path", "method"},
	)
	prometheus.MustRegister(counter, duration)
	return &PathInstrumentor{counter: counter, duration: duration}
}

func (h *PathInstrumentor) Instrument(path string, f http.Handler) http.Handler {

	instrumented := promhttp.InstrumentHandlerDuration(h.duration.MustCurryWith(prometheus.Labels{"path": path}),
		promhttp.InstrumentHandlerCounter(h.counter.MustCurryWith(prometheus.Labels{"path": path}), f))

	return instrumented
}

func (h *PathInstrumentor) Shutdown() {
	prometheus.Unregister(h.duration)
	prometheus.Unregister(h.counter)
}

package metrics

import (
	"net/http"
	"strconv"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

const (
	metricName  = "projeto_korp_http_requests_total"
	metricsPath = "/metrics"
	projectPath = "/projeto-korp"
)

// Component owns application request metrics and their Prometheus exposition.
type Component struct {
	requests *prometheus.CounterVec
	registry *prometheus.Registry
}

// New creates an isolated metrics registry for one application instance.
func New() *Component {
	requests := prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: metricName,
			Help: "Total number of HTTP requests handled by the Projeto Korp application.",
		},
		[]string{"method", "route", "status"},
	)
	registry := prometheus.NewRegistry()
	registry.MustRegister(requests)

	return &Component{
		requests: requests,
		registry: registry,
	}
}

// Handler exposes /metrics and instruments the application router.
func (c *Component) Handler(application http.Handler) http.Handler {
	mux := http.NewServeMux()
	mux.Handle("GET "+metricsPath, promhttp.HandlerFor(c.registry, promhttp.HandlerOpts{}))
	mux.Handle(metricsPath+"/", c.Middleware(application))
	mux.Handle("/", c.Middleware(application))
	return mux
}

// Middleware records final application response statuses without coupling the
// Project Endpoint Handler to Prometheus.
func (c *Component) Middleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == metricsPath {
			next.ServeHTTP(w, r)
			return
		}

		writer := &statusWriter{ResponseWriter: w}
		next.ServeHTTP(writer, r)

		route := "unmatched"
		if r.URL.Path == projectPath {
			route = projectPath
		}
		c.requests.WithLabelValues(r.Method, route, strconv.Itoa(writer.statusCode())).Inc()
	})
}

type statusWriter struct {
	http.ResponseWriter
	status int
}

func (w *statusWriter) WriteHeader(status int) {
	if w.status != 0 {
		return
	}
	w.status = status
	w.ResponseWriter.WriteHeader(status)
}

func (w *statusWriter) Write(body []byte) (int, error) {
	if w.status == 0 {
		w.WriteHeader(http.StatusOK)
	}
	return w.ResponseWriter.Write(body)
}

func (w *statusWriter) statusCode() int {
	if w.status == 0 {
		return http.StatusOK
	}
	return w.status
}

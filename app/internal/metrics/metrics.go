package metrics

import (
	"net/http"
	"strconv"
	"time"

	"github.com/projeto-korp/app/internal/audit"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

const (
	metricName  = "projeto_korp_http_requests_total"
	metricsPath = "/metrics"
	projectPath = "/projeto-korp"
)

const (
	cacheOperationsMetric      = "projeto_korp_statistics_cache_operations_total"
	cacheDurationMetric        = "projeto_korp_statistics_cache_operation_duration_seconds"
	dependencyOperationsMetric = "projeto_korp_data_dependency_operations_total"
	dependencyDurationMetric   = "projeto_korp_data_dependency_operation_duration_seconds"
)

// Component owns application request metrics and their Prometheus exposition.
type Component struct {
	requests             *prometheus.CounterVec
	cacheOperations      *prometheus.CounterVec
	cacheDuration        *prometheus.HistogramVec
	dependencyOperations *prometheus.CounterVec
	dependencyDuration   *prometheus.HistogramVec
	registry             *prometheus.Registry
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
	cacheOperations := prometheus.NewCounterVec(prometheus.CounterOpts{Name: cacheOperationsMetric, Help: "Statistics cache operation outcomes."}, []string{"operation", "result"})
	cacheDuration := prometheus.NewHistogramVec(prometheus.HistogramOpts{Name: cacheDurationMetric, Help: "Statistics cache operation duration in seconds."}, []string{"operation"})
	dependencyOperations := prometheus.NewCounterVec(prometheus.CounterOpts{Name: dependencyOperationsMetric, Help: "Application data dependency operation outcomes."}, []string{"dependency", "operation", "result"})
	dependencyDuration := prometheus.NewHistogramVec(prometheus.HistogramOpts{Name: dependencyDurationMetric, Help: "Application data dependency operation duration in seconds."}, []string{"dependency", "operation"})
	registry := prometheus.NewRegistry()
	registry.MustRegister(requests)
	registry.MustRegister(cacheOperations, cacheDuration, dependencyOperations, dependencyDuration)

	return &Component{
		requests:             requests,
		cacheOperations:      cacheOperations,
		cacheDuration:        cacheDuration,
		dependencyOperations: dependencyOperations,
		dependencyDuration:   dependencyDuration,
		registry:             registry,
	}
}

func (c *Component) ObserveCacheLookup(result string, duration time.Duration) {
	c.cacheOperations.WithLabelValues("lookup", result).Inc()
	c.cacheDuration.WithLabelValues("lookup").Observe(duration.Seconds())
}

func (c *Component) ObserveCacheWrite(result string, duration time.Duration) {
	c.cacheOperations.WithLabelValues("write", result).Inc()
	c.cacheDuration.WithLabelValues("write").Observe(duration.Seconds())
}

func (c *Component) ObserveRepository(operation, result string, duration time.Duration) {
	c.dependencyOperations.WithLabelValues("postgresql", operation, result).Inc()
	c.dependencyDuration.WithLabelValues("postgresql", operation).Observe(duration.Seconds())
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

		route := audit.CanonicalRoute(r.URL.Path)
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

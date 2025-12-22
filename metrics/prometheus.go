package metrics

import (
	"net/http"
	"time"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

var (
	// TotalRequests counts all HTTP requests by method and endpoint.
	TotalRequests = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "http_requests_total",
			Help: "Total number of HTTP requests",
		},
		[]string{"method", "endpoint"},
	)

	// RequestDuration captures latency distributions.
	RequestDuration = prometheus.NewHistogramVec(
		prometheus.HistogramOpts{
			Name: "http_request_duration_seconds",
			Help: "Request duration in seconds",
			// Use Prometheus defaults; adjust buckets if needed.
		},
		[]string{"method", "endpoint"},
	)
)

func init() {
	prometheus.MustRegister(TotalRequests)
	prometheus.MustRegister(RequestDuration)
}

// MetricsMiddleware records counters and latency.
func MetricsMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		TotalRequests.WithLabelValues(r.Method, r.URL.Path).Inc()
		next.ServeHTTP(w, r)
		RequestDuration.WithLabelValues(r.Method, r.URL.Path).Observe(time.Since(start).Seconds())
	})
}

// Handler exposes /metrics endpoint.
func Handler() http.Handler {
	return promhttp.Handler()
}

package metrics

import (
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

// wrapper for ResponseWriter class
type responseWriter struct {
	http.ResponseWriter
	statusCode int
}

func (r *responseWriter) WriteHeader(status int) {
	r.statusCode = status
	r.ResponseWriter.WriteHeader(status)
}

var (
	// The Prometheus metrics that will be exposed.
	httpHits = prometheus.NewCounter(
		prometheus.CounterOpts{
			Name: "http_hit_total",
			Help: "Total number of HTTP hits.",
		},
	)

	trafficAccumulationMetric = prometheus.NewCounter(
		prometheus.CounterOpts{
			Name: "traffic_total",
			Help: "Total traffic in GB.",
		},
	)

	httpHitsSuccess = prometheus.NewCounter(
		prometheus.CounterOpts{
			Name: "success_http_hit_total",
			Help: "Total number of successfull HTTP hits.",
		},
	)

	httpHitsFailed = prometheus.NewCounter(
		prometheus.CounterOpts{
			Name: "failed_http_hit_total",
			Help: "Total number of failed HTTP hits.",
		},
	)

	httpStatusNotFoundCounter = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "not_found_response_status_counter",
			Help: "Counter for the 404 status of the HTTP response.",
		},
		[]string{"status"})

	uniqueVisitorCounter = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "unique_visitor_counter",
			Help: "Counter for each unique visitor.",
		},
		[]string{"visitor"})

	// Add all metrics that will be resisted
	metricsList = []prometheus.Collector{
		httpHits,
		httpHitsSuccess,
		httpHitsFailed,
		uniqueVisitorCounter,
		httpStatusNotFoundCounter,
		trafficAccumulationMetric,
	}

	// Prometheus Registry to register metrics.
	prometheusRegistry = prometheus.NewRegistry()
)

func init() {
	// Register metrics that will be exposed.
	prometheusRegistry.MustRegister(metricsList...)
}

func MetricsHandler() http.Handler {
	return promhttp.HandlerFor(prometheusRegistry, promhttp.HandlerOpts{})
}

func MetricProxy(f func(http.ResponseWriter, *http.Request)) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		route := mux.CurrentRoute(r)
		path, _ := route.GetPathTemplate()

		userAgent := r.Header.Get("User-Agent")
		rw := &responseWriter{w, http.StatusOK}

		f(rw, r) // original function call

		httpHits.Inc()

		requestSize, _ := strconv.Atoi((r.Header.Get("Content-Length")))
		trafficAccumulationMetric.Add(float64(requestSize) / 1e9)

		uniqueVisitorCounter.WithLabelValues(r.RemoteAddr + " (" + r.Host + ") via " + userAgent).Inc()

		if rw.statusCode >= 200 && rw.statusCode < 400 {
			httpHitsSuccess.Inc()
		}

		if rw.statusCode >= 400 && rw.statusCode < 600 {
			httpHitsFailed.Inc()
		}

		if rw.statusCode == 404 {
			httpStatusNotFoundCounter.WithLabelValues(path).Inc()
		}
	}
}

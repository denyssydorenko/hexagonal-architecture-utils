package metrics

import (
	"github.com/labstack/echo/v4"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

var (
	namespace   string = "hexagonal_architecture_utils"
	APIDuration *prometheus.HistogramVec
	DBDuration  *prometheus.HistogramVec
	APIRequests *prometheus.CounterVec
	DBQueries   *prometheus.CounterVec
)

func InitMetrics() {
	APIDuration = promauto.NewHistogramVec(prometheus.HistogramOpts{
		Namespace: namespace,
		Subsystem: "api",
		Name:      "histogram_seconds",
		Help:      "A histogram of api latencies.",
		Buckets:   prometheus.DefBuckets,
	}, []string{"endpoint", "success"})
	DBDuration = promauto.NewHistogramVec(prometheus.HistogramOpts{
		Namespace: namespace,
		Subsystem: "postgres",
		Name:      "query_duration_histogram_seconds",
		Help:      "A histogram of query latencies.",
		Buckets:   prometheus.DefBuckets,
	}, []string{"host", "database", "query", "success"})
	APIRequests = promauto.NewCounterVec(prometheus.CounterOpts{
		Namespace: namespace,
		Subsystem: "api",
		Name:      "requests_total",
		Help:      "A counter for requests.",
	}, []string{"endpoint", "success"})
	DBQueries = promauto.NewCounterVec(prometheus.CounterOpts{
		Namespace: namespace,
		Subsystem: "postgres",
		Name:      "queries_total",
		Help:      "A counter for queries.",
	}, []string{"host", "database", "query", "success"})
}

// Metrics
// @Summary Service prometheus metrics
// @Description Shows service prometheus metrics
// @Tags api
// @Produce plain
// @Success 200 {string} string
// @Router /api/metrics [get]
func Metrics() echo.HandlerFunc {
	return echo.WrapHandler(promhttp.Handler())
}

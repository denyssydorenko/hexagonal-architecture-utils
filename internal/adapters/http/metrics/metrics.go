package metrics

import (
	"github.com/labstack/echo/v4"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

var (
	namespace  string = "hexagonal_architecture_utils"
	DBDuration *prometheus.HistogramVec
	DBQueries  *prometheus.CounterVec
)

func InitMetrics() {
	DBDuration = promauto.NewHistogramVec(prometheus.HistogramOpts{
		Namespace: namespace,
		Subsystem: "postgres",
		Name:      "query_duration_histogram_seconds",
		Help:      "A histogram of query latencies.",
		Buckets:   prometheus.DefBuckets,
	}, []string{"host", "database", "query", "success"})
	DBQueries = promauto.NewCounterVec(prometheus.CounterOpts{
		Namespace: namespace,
		Subsystem: "postgres",
		Name:      "queries_total",
		Help:      "A counter for queries.",
	}, []string{"host", "database", "query", "success"})
}

// Metrics godoc
// @Summary Shows service prometheus metrics
// @Description Shows service prometheus metrics
// @Tags api
// @ID metrics
// @Produce plain
// @Success 200 {string} string
// @Router /api/metrics [get]
func Metrics() echo.HandlerFunc {
	return echo.WrapHandler(promhttp.Handler())
}

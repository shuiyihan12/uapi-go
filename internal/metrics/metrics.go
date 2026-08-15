// Package metrics provides Prometheus-based monitoring metrics for UAPI
// request volume, latency, errors and in-flight requests.
package metrics

import (
	"sync"
	"time"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
)

// Metrics holds the monitoring metrics.
type Metrics struct {
	// RequestTotal counts completed requests by service, operation and status.
	RequestTotal *prometheus.CounterVec
	// RequestDuration observes request latency by service and operation.
	RequestDuration *prometheus.HistogramVec
	// RequestErrors counts request errors by service, operation and error type.
	RequestErrors *prometheus.CounterVec
	// ActiveRequests gauges currently in-flight requests by service and operation.
	ActiveRequests *prometheus.GaugeVec
}

var (
	instance *Metrics
	once     sync.Once
)

// GetMetrics returns the global metrics instance (singleton).
func GetMetrics() *Metrics {
	once.Do(func() {
		instance = &Metrics{
			RequestTotal: promauto.NewCounterVec(
				prometheus.CounterOpts{
					Name: "uapi_requests_total",
					Help: "Total number of UAPI requests",
				},
				[]string{"service", "operation", "status"},
			),
			RequestDuration: promauto.NewHistogramVec(
				prometheus.HistogramOpts{
					Name:    "uapi_request_duration_seconds",
					Help:    "Duration of UAPI requests in seconds",
					Buckets: prometheus.DefBuckets,
				},
				[]string{"service", "operation"},
			),
			RequestErrors: promauto.NewCounterVec(
				prometheus.CounterOpts{
					Name: "uapi_request_errors_total",
					Help: "Total number of UAPI request errors",
				},
				[]string{"service", "operation", "error_type"},
			),
			ActiveRequests: promauto.NewGaugeVec(
				prometheus.GaugeOpts{
					Name: "uapi_active_requests",
					Help: "Number of active UAPI requests",
				},
				[]string{"service", "operation"},
			),
		}
	})
	return instance
}

// NewMetrics creates the monitoring metrics (deprecated; use GetMetrics).
func NewMetrics() *Metrics {
	return GetMetrics()
}

// RecordRequest records metrics for one completed request.
func (m *Metrics) RecordRequest(service, operation string, duration time.Duration, err error) {
	status := "success"
	if err != nil {
		status = "error"
		m.RequestErrors.WithLabelValues(service, operation, "unknown").Inc()
	}

	m.RequestTotal.WithLabelValues(service, operation, status).Inc()
	m.RequestDuration.WithLabelValues(service, operation).Observe(duration.Seconds())
}

// IncActiveRequests increments the in-flight request gauge.
func (m *Metrics) IncActiveRequests(service, operation string) {
	m.ActiveRequests.WithLabelValues(service, operation).Inc()
}

// DecActiveRequests decrements the in-flight request gauge.
func (m *Metrics) DecActiveRequests(service, operation string) {
	m.ActiveRequests.WithLabelValues(service, operation).Dec()
}

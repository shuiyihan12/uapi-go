package client

import "time"

// Metrics is the observability hook invoked by EnterpriseSOAPClient around
// every SOAP call. Implementations are supplied by the deployment: the
// gateway daemon injects its Prometheus collector, while SDK consumers may
// pass their own or omit it entirely (a no-op implementation is used then).
// The parameter types are deliberately scalar so implementations need no
// Prometheus or zap dependency.
type Metrics interface {
	// RecordRequest records one completed request.
	RecordRequest(service, operation string, duration time.Duration, err error)
	// IncActiveRequests increments the in-flight gauge.
	IncActiveRequests(service, operation string)
	// DecActiveRequests decrements the in-flight gauge.
	DecActiveRequests(service, operation string)
}

// noopMetrics discards everything; used when no Metrics implementation is
// configured.
type noopMetrics struct{}

func (noopMetrics) RecordRequest(string, string, time.Duration, error) {}
func (noopMetrics) IncActiveRequests(string, string)                   {}
func (noopMetrics) DecActiveRequests(string, string)                   {}

// Package client wraps Travelport UAPI SOAP calls with auth pass-through,
// trace injection and observability (metrics, leveled logging). This project
// does not retry: any single failed SOAP call returns immediately and the
// caller decides how to handle it (e.g. surfacing a GDS business error).
package client

import (
	"context"
	"encoding/xml"
	"errors"
	"fmt"
	"time"

	"github.com/shuiyihan12/uapi-go/pkg/logging"
	"github.com/shuiyihan12/uapi-go/pkg/requestctx"
	"go.uber.org/zap"
)

// EnterpriseSOAPClient layers metrics collection and leveled logging on top
// of SOAPClient.
type EnterpriseSOAPClient struct {
	*SOAPClient
	logger  logging.Logger
	metrics Metrics
	service string
}

// EnterpriseConfig aggregates the SOAP configuration and service name used
// to build the enterprise client.
type EnterpriseConfig struct {
	SOAPConfig  SOAPConfig
	ServiceName string
	Logger      logging.Logger
}

// NewEnterpriseSOAPClient builds an EnterpriseSOAPClient and initializes its
// metrics.
func NewEnterpriseSOAPClient(config EnterpriseConfig) (*EnterpriseSOAPClient, error) {
	logger := config.Logger
	if logger == nil {
		logger = logging.NewDefaultLogger()
	}

	config.SOAPConfig.Logger = logger
	soapClient, err := NewSOAPClient(config.SOAPConfig)
	if err != nil {
		return nil, err
	}

	m := config.SOAPConfig.Metrics
	if m == nil {
		m = noopMetrics{}
	}

	return &EnterpriseSOAPClient{
		SOAPClient: soapClient,
		logger:     logger,
		metrics:    m,
		service:    config.ServiceName,
	}, nil
}

// CallWithObservability performs one SOAP call with metrics collection and
// leveled logging (no retries).
//
// Any error (network, timeout, system or GDS business error) is returned
// directly without retrying.
func (c *EnterpriseSOAPClient) CallWithObservability(ctx context.Context, operation string, request interface{}) ([]byte, error) {
	start := time.Now()

	// Track the in-flight request count.
	c.metrics.IncActiveRequests(c.service, operation)
	defer c.metrics.DecActiveRequests(c.service, operation)

	// Log the request start.
	c.logger.WithContext(ctx).Info("Starting SOAP operation",
		zap.String("service", c.service),
		zap.String("operation", operation),
		zap.String("endpoint", c.resolveEndpoint(requestctx.NormalizeRegion(requestctx.Region(ctx)))))

	response, err := c.SOAPClient.Call(ctx, operation, request)

	// Record request completion.
	duration := time.Since(start)
	c.metrics.RecordRequest(c.service, operation, duration, err)

	if err != nil {
		if fault, ok := asSOAPFault(err); ok && !fault.Retryable() {
			// Business/client errors (e.g. unavailable inventory, invalid
			// parameters) are explicit business conditions returned by the
			// GDS, not failures of this service. Log them as Warn (no
			// stack trace) for business metrics while keeping ERROR-level
			// alerts and dashboards clean.
			c.logger.WithContext(ctx).Warn("SOAP operation returned business error",
				zap.String("service", c.service),
				zap.String("operation", operation),
				zap.Duration("duration", duration),
				zap.String("code", fault.Code),
				zap.String("type", fault.Type),
				zap.String("provider_service", fault.Service),
				zap.String("fault_code", fault.FaultCode),
				zap.String("description", fault.Description))
		} else {
			// Genuine system/network/timeout/parsing errors are logged as
			// Error with stack traces for troubleshooting.
			c.logger.WithContext(ctx).Error("SOAP operation failed",
				zap.String("service", c.service),
				zap.String("operation", operation),
				zap.Duration("duration", duration),
				zap.Error(err))
		}
		return nil, err
	}

	c.logger.WithContext(ctx).Info("SOAP operation completed successfully",
		zap.String("service", c.service),
		zap.String("operation", operation),
		zap.Duration("duration", duration))

	return response, nil
}

// asSOAPFault extracts *SOAPFaultError from the error chain, if present.
func asSOAPFault(err error) (*SOAPFaultError, bool) {
	var fault *SOAPFaultError
	if errors.As(err, &fault) {
		return fault, true
	}
	return nil, false
}

// CallPortType is a type-safe wrapper around CallWithObservability: it makes
// one SOAP call and parses the raw response body into the strongly typed
// response T of the WSDL PortType operation (e.g. *hotel.HotelDetailsRsp).
//
// Service layers use it to replace string-named operations with strongly
// typed methods mirroring the WSDL *PortType, keeping observability (metrics,
// leveled logging) while eliminating the hardcoding risk of hand-written
// operation strings. This project does not retry; every error is returned
// directly.
func CallPortType[T any](c *EnterpriseSOAPClient, ctx context.Context, operation string, request any) (*T, error) {
	body, err := c.CallWithObservability(ctx, operation, request)
	if err != nil {
		return nil, err
	}

	var rsp T
	if err := xml.Unmarshal(body, &rsp); err != nil {
		return nil, fmt.Errorf("%s: failed to parse response: %w", operation, err)
	}
	return &rsp, nil
}

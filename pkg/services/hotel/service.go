// Package hotel provides the SOAP client implementation for the Travelport
// Hotel service.
//
// All request/response types uniformly use the WSDL-generated hotel package
// types (pkg/generated/hotel); no hand-written SOAP models are maintained
// anymore:
//   - strongly typed Port operations are in hotel_port.go (HotelServicePort);
//   - the REST business surface (search /search, details /details, including
//     automatic pagination) is in hotel_soap_api.go.
package hotel

import (
	"fmt"

	"github.com/shuiyihan12/uapi-go/pkg/client"
	"github.com/shuiyihan12/uapi-go/pkg/logging"
)

// HotelService is the SOAP client for the Travelport Hotel service.
type HotelService struct {
	client   *client.EnterpriseSOAPClient
	logger   logging.Logger
	endpoint string
}

// NewHotelService builds a Hotel service client from the given SOAP configuration and logger.
func NewHotelService(config client.SOAPConfig, logger logging.Logger) (*HotelService, error) {
	enterpriseConfig := client.EnterpriseConfig{
		SOAPConfig:  config,
		ServiceName: "hotel-service",
		Logger:      logger,
	}

	enterpriseClient, err := client.NewEnterpriseSOAPClient(enterpriseConfig)
	if err != nil {
		return nil, fmt.Errorf("failed to create hotel service client: %v", err)
	}

	return &HotelService{
		client:   enterpriseClient,
		logger:   logger,
		endpoint: config.BaseEndpoint,
	}, nil
}

// Close closes the underlying SOAP client connection and releases its resources.
func (s *HotelService) Close() error {
	return s.client.Close()
}

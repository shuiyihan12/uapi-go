package hotel

import (
	"context"

	"github.com/shuiyihan12/uapi-go/pkg/client"
	hotelxsd "github.com/shuiyihan12/uapi-go/pkg/generated/hotel"
)

// HotelServicePort mirrors the *PortType operations of hotel_v55_0 one-to-one,
// making it easy to swap implementations and plug in test stubs.
//
// Request and response types now uniformly use the WSDL-generated hotel package
// types: the generator has fixed resolution of <xs:element ref="...">
// cross-element references (which previously generated child elements as
// interface{} and crashed on serialization, forcing these requests to use the
// hand-written models in pkg/services/hotel), so the hand-written models were
// removed and everything moved back to the generated types. Request-level
// business parameters (TargetBranch, BillingPointOfSaleInfo, etc.) are provided
// explicitly by the caller in the request body; response types are written back
// as-is, with no per-field mapping.
//
// The existing SearchAvailability/Details (used by the REST business surface)
// are kept; this interface targets scenarios that need to call strongly typed
// GDS operations directly, eliminating the hard-coded risk of hand-written
// operation strings.
type HotelServicePort interface {
	// HotelSearchAvailability corresponds to HotelSearchServicePortType and performs a hotel availability search.
	HotelSearchAvailability(ctx context.Context, req *hotelxsd.HotelSearchAvailabilityReq) (*hotelxsd.BaseHotelSearchRsp, error)
	// HotelMediaLinks corresponds to HotelMediaLinksServicePortType and retrieves hotel media links.
	HotelMediaLinks(ctx context.Context, req *hotelxsd.HotelMediaLinksReq) (*hotelxsd.HotelMediaLinksRsp, error)
	// HotelDetails corresponds to HotelDetailsServicePortType and retrieves hotel details and rates.
	HotelDetails(ctx context.Context, req *hotelxsd.HotelDetailsReq) (*hotelxsd.HotelDetailsRsp, error)
	// HotelRetrieve corresponds to HotelRetrieveServicePortType and retrieves an existing hotel reservation by its record locator.
	HotelRetrieve(ctx context.Context, req *hotelxsd.HotelRetrieveReq) (*hotelxsd.HotelRetrieveRsp, error)
	// HotelRules corresponds to HotelRulesServicePortType and retrieves hotel Fare rules.
	HotelRules(ctx context.Context, req *hotelxsd.HotelRulesReq) (*hotelxsd.HotelRulesRsp, error)
	// HotelUpsellSearch corresponds to HotelUpsellSearchServicePortType and searches hotel upsell offers.
	HotelUpsellSearch(ctx context.Context, req *hotelxsd.HotelUpsellDetailsReq) (*hotelxsd.HotelUpsellDetailsRsp, error)
	// HotelKeywords corresponds to HotelKeywordsServicePortType and searches hotels by keyword.
	HotelKeywords(ctx context.Context, req *hotelxsd.HotelKeywordReq) (*hotelxsd.HotelKeywordRsp, error)
	// HotelSuperShopper corresponds to HotelSuperShopperServicePortType and performs a hotel super shopper search.
	HotelSuperShopper(ctx context.Context, req *hotelxsd.HotelSuperShopperReq) (*hotelxsd.HotelSuperShopperRsp, error)
}

// Compile-time assertion: *HotelService must satisfy the HotelServicePort interface.
var _ HotelServicePort = (*HotelService)(nil)

// HotelSearchAvailability issues the HotelSearchAvailability SOAP call and returns the strongly typed response.
func (s *HotelService) HotelSearchAvailability(ctx context.Context, req *hotelxsd.HotelSearchAvailabilityReq) (*hotelxsd.BaseHotelSearchRsp, error) {
	return callPort[hotelxsd.BaseHotelSearchRsp](s.client, ctx, "HotelSearchAvailability", req)
}

// HotelMediaLinks issues the HotelMediaLinks SOAP call and returns the strongly typed response.
func (s *HotelService) HotelMediaLinks(ctx context.Context, req *hotelxsd.HotelMediaLinksReq) (*hotelxsd.HotelMediaLinksRsp, error) {
	return callPort[hotelxsd.HotelMediaLinksRsp](s.client, ctx, "HotelMediaLinks", req)
}

// HotelDetails issues the HotelDetails SOAP call and returns the strongly typed response.
func (s *HotelService) HotelDetails(ctx context.Context, req *hotelxsd.HotelDetailsReq) (*hotelxsd.HotelDetailsRsp, error) {
	return callPort[hotelxsd.HotelDetailsRsp](s.client, ctx, "HotelDetails", req)
}

// HotelRetrieve issues the HotelRetrieve SOAP call and returns the strongly typed response.
func (s *HotelService) HotelRetrieve(ctx context.Context, req *hotelxsd.HotelRetrieveReq) (*hotelxsd.HotelRetrieveRsp, error) {
	return callPort[hotelxsd.HotelRetrieveRsp](s.client, ctx, "HotelRetrieve", req)
}

// HotelRules issues the HotelRules SOAP call and returns the strongly typed response.
func (s *HotelService) HotelRules(ctx context.Context, req *hotelxsd.HotelRulesReq) (*hotelxsd.HotelRulesRsp, error) {
	return callPort[hotelxsd.HotelRulesRsp](s.client, ctx, "HotelRules", req)
}

// HotelUpsellSearch issues the HotelUpsellSearch SOAP call and returns the strongly typed response.
func (s *HotelService) HotelUpsellSearch(ctx context.Context, req *hotelxsd.HotelUpsellDetailsReq) (*hotelxsd.HotelUpsellDetailsRsp, error) {
	return callPort[hotelxsd.HotelUpsellDetailsRsp](s.client, ctx, "HotelUpsellSearch", req)
}

// HotelKeywords issues the HotelKeywords SOAP call and returns the strongly typed response.
func (s *HotelService) HotelKeywords(ctx context.Context, req *hotelxsd.HotelKeywordReq) (*hotelxsd.HotelKeywordRsp, error) {
	return callPort[hotelxsd.HotelKeywordRsp](s.client, ctx, "HotelKeywords", req)
}

// HotelSuperShopper issues the HotelSuperShopper SOAP call and returns the strongly typed response.
func (s *HotelService) HotelSuperShopper(ctx context.Context, req *hotelxsd.HotelSuperShopperReq) (*hotelxsd.HotelSuperShopperRsp, error) {
	return callPort[hotelxsd.HotelSuperShopperRsp](s.client, ctx, "HotelSuperShopper", req)
}

// callPort is a package-local convenience wrapper around client.CallPortType that
// performs a single SOAP call and decodes it into the strongly typed response T.
func callPort[T any](c *client.EnterpriseSOAPClient, ctx context.Context, operation string, req any) (*T, error) {
	return client.CallPortType[T](c, ctx, operation, req)
}

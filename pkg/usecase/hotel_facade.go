package usecase

import (
	"context"
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/shuiyihan12/uapi-go/pkg/generated/common55"
	hotelxsd "github.com/shuiyihan12/uapi-go/pkg/generated/hotel"
	hotelenums "github.com/shuiyihan12/uapi-go/pkg/generated/hotel/enums"
	"github.com/shuiyihan12/uapi-go/pkg/manager"
	hotelsvc "github.com/shuiyihan12/uapi-go/pkg/services/hotel"
	"github.com/shuiyihan12/uapi-go/pkg/trace"
)

// HotelFacade orchestrates the hotel domain use cases.
//
// The request side keeps no separate DTOs: the public API's request bodies
// deserialize directly into the WSDL-generated hotel package types (the
// generated types carry both json and xml tags). This layer only does three
// things — input validation, enum normalization, and non-empty checks of
// request-level required fields (TargetBranch, BillingPointOfSaleInfo, etc.)
// plus TraceId injection — with no field shuffling. The business fields
// above are always provided explicitly by the caller in the request body.
type HotelFacade struct {
	manager *manager.ServiceManager
}

// NewHotelFacade creates the hotel use-case layer.
func NewHotelFacade(serviceManager *manager.ServiceManager) *HotelFacade {
	return &HotelFacade{
		manager: serviceManager,
	}
}

// HotelSearchAvailabilityOutput is the result returned to the web layer.
type HotelSearchAvailabilityOutput struct {
	Response *hotelxsd.BaseHotelSearchRsp `json:"response"` // full search response carried by the WSDL-generated model
}

// HotelDetailsOutput is the result returned to the web layer.
type HotelDetailsOutput struct {
	Response *hotelxsd.HotelDetailsRsp `json:"response"` // full details response carried by the WSDL-generated model
}

// validateHotelStay verifies that the check-in date is not in the past and
// the check-out date is after check-in, returning a ValidationError
// otherwise.
func validateHotelStay(stay hotelxsd.HotelStay) error {
	checkin, err := parseHotelDate(string(stay.CheckinDate), "checkinDate")
	if err != nil {
		return err
	}
	checkout, err := parseHotelDate(string(stay.CheckoutDate), "checkoutDate")
	if err != nil {
		return err
	}

	now := time.Now().In(time.Local)
	today := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, time.Local)
	if checkin.Before(today) {
		return NewValidationError("checkinDate", "不能早于今天")
	}
	if checkout.Before(today) {
		return NewValidationError("checkoutDate", "不能早于今天")
	}
	if !checkout.After(checkin) {
		return NewValidationError("checkoutDate", "必须晚于 checkinDate")
	}
	return nil
}

// parseHotelDate parses a YYYY-MM-DD string into a time.Time in the local
// timezone; blank or malformed values return a ValidationError.
func parseHotelDate(value, fieldName string) (time.Time, error) {
	if strings.TrimSpace(value) == "" {
		return time.Time{}, NewValidationError(fieldName, "不能为空")
	}
	parsed, err := time.ParseInLocation("2006-01-02", value, time.Local)
	if err != nil {
		return time.Time{}, NewValidationError(fieldName, "必须是 YYYY-MM-DD 格式")
	}
	return parsed, nil
}

// normalizeRateRuleDetail normalizes a caller-supplied rate rule detail
// value to None/Complete/RatePlansOnly; invalid values return a
// ValidationError.
func normalizeRateRuleDetail(value string) (string, error) {
	switch strings.ToLower(strings.TrimSpace(value)) {
	case "":
		return "", nil
	case "none":
		return "None", nil
	case "complete":
		return "Complete", nil
	case "rateplansonly", "rate_plans_only", "rate-plans-only":
		return "RatePlansOnly", nil
	default:
		return "", NewValidationError("rateRuleDetail", "只能是 None、Complete 或 RatePlansOnly")
	}
}

// SearchAvailability performs a hotel availability search.
//
// req is the generated type deserialized straight from the public API's JSON
// body; after in-place validation this method hands it to the service layer
// for SOAP serialization. Request-level business parameters such as
// TargetBranch and BillingPointOfSaleInfo are provided explicitly by the
// caller in the request body.
func (f *HotelFacade) SearchAvailability(ctx context.Context, req *hotelxsd.HotelSearchAvailabilityReq) (*HotelSearchAvailabilityOutput, error) {
	if f.manager == nil {
		return nil, errors.New("service manager is nil")
	}
	if req == nil {
		return nil, NewValidationError("body", "请求体不能为空")
	}

	if req.HotelSearchLocation == nil {
		return nil, NewValidationError("hotelSearchLocation", "不能为空")
	}
	loc := req.HotelSearchLocation
	if loc.HotelLocation == nil && len(loc.VendorLocation) == 0 && loc.HotelAddress == nil &&
		loc.ReferencePoint == nil && loc.CoordinateLocation == nil && loc.Distance == nil {
		return nil, NewValidationError("hotelSearchLocation", "至少提供一种定位方式（hotelLocation / vendorLocation / hotelAddress / referencePoint / coordinateLocation / distance）")
	}
	if err := validateHotelStay(req.HotelStay); err != nil {
		return nil, err
	}

	ctx, err := normalizeHotelPortReq(ctx, &req.BaseReq)
	if err != nil {
		return nil, err
	}

	hotelService, err := manager.Get[*hotelsvc.HotelService](f.manager, "hotel")
	if err != nil {
		return nil, fmt.Errorf("failed to get hotel service: %w", err)
	}

	response, err := hotelService.SearchAvailability(ctx, req)
	if err != nil {
		return nil, err
	}

	return &HotelSearchAvailabilityOutput{
		Response: response,
	}, nil
}

// Details retrieves hotel details.
//
// req is the generated type deserialized straight from the public API's JSON
// body; after in-place validation and enum normalization this method hands
// it to the service layer for SOAP serialization (with automatic paging).
// Request-level business parameters such as TargetBranch and
// BillingPointOfSaleInfo are provided explicitly by the caller in the
// request body.
func (f *HotelFacade) Details(ctx context.Context, req *hotelxsd.HotelDetailsReq) (*HotelDetailsOutput, error) {
	if f.manager == nil {
		return nil, errors.New("service manager is nil")
	}
	if req == nil {
		return nil, NewValidationError("body", "请求体不能为空")
	}

	if strings.TrimSpace(string(req.HotelProperty.HotelChain)) == "" || strings.TrimSpace(string(req.HotelProperty.HotelCode)) == "" {
		return nil, NewValidationError("hotelProperty", "hotelChain 和 hotelCode 不能为空")
	}
	if req.HotelDetailsModifiers != nil {
		if err := validateHotelStay(req.HotelDetailsModifiers.HotelStay); err != nil {
			return nil, err
		}
		if v := req.HotelDetailsModifiers.RateRuleDetail; v != nil {
			normalized, err := normalizeRateRuleDetail(string(*v))
			if err != nil {
				return nil, err
			}
			casted := hotelenums.TypeRateRuleDetail(normalized)
			req.HotelDetailsModifiers.RateRuleDetail = &casted
		}
	}

	ctx, err := normalizeHotelPortReq(ctx, &req.BaseReq)
	if err != nil {
		return nil, err
	}

	hotelService, err := manager.Get[*hotelsvc.HotelService](f.manager, "hotel")
	if err != nil {
		return nil, fmt.Errorf("failed to get hotel service: %w", err)
	}

	response, err := hotelService.Details(ctx, req)
	if err != nil {
		return nil, err
	}

	return &HotelDetailsOutput{
		Response: response,
	}, nil
}

// The 6 methods below correspond to strongly typed HotelServicePort
// operations with no REST equivalents and are forwarded as pass-throughs.
// The request bodies are already deserialized by the API layer into the
// WSDL-generated hotel models (no DTOs, no hand-written models); this layer
// only performs the non-empty checks of request-level required fields
// (TargetBranch etc.) and injects the trace context, with no field
// shuffling. The existing /search and /details (the REST business surface)
// already cover HotelSearchAvailability and HotelDetails, so they are not
// duplicated here to avoid two entry points for one operation.

// normalizeHotelPortReq validates and normalizes the request-level fields
// shared by the 6 pass-through operations and returns the ctx with the trace
// ID injected.
//
// All 6 pass-through request types embed common55.BaseReq (HotelSuperShopperReq
// indirectly via BaseSearchReq), so they are handled uniformly as
// *common55.BaseReq instead of repeating the same validation per operation.
//
// Trace injection now uses trace.Ensure plus the generated type's
// InjectInfrastructure to write TraceId: before the migration it read
// TransactionID from the hand-written model, but that field does not exist
// in the XSD (BaseCoreReq only has TraceId) — it was an invented field of
// the hand-written model. Now this service's trace ID is passed down to
// Travelport for two-way reconciliation. Business fields such as
// TargetBranch and BillingPointOfSaleInfo still must be provided explicitly
// by the caller in the request body; this layer does not inject them.
func normalizeHotelPortReq(ctx context.Context, base *common55.BaseReq) (context.Context, error) {
	if base.TargetBranch == nil {
		return ctx, NewValidationError("targetBranch", "不能为空")
	}
	branch := common55.TypeBranchCode(strings.TrimSpace(string(*base.TargetBranch)))
	if branch == "" {
		return ctx, NewValidationError("targetBranch", "不能为空")
	}
	base.TargetBranch = &branch

	ctx, traceID := trace.Ensure(ctx)
	base.InjectInfrastructure(traceID)
	return ctx, nil
}

// MediaLinks retrieves hotel media links (HotelMediaLinksServicePortType).
func (f *HotelFacade) MediaLinks(ctx context.Context, req *hotelxsd.HotelMediaLinksReq) (*hotelxsd.HotelMediaLinksRsp, error) {
	if f.manager == nil {
		return nil, errors.New("service manager is nil")
	}
	if req == nil {
		return nil, NewValidationError("body", "请求体不能为空")
	}
	ctx, err := normalizeHotelPortReq(ctx, &req.BaseReq)
	if err != nil {
		return nil, err
	}

	hotelService, err := manager.Get[*hotelsvc.HotelService](f.manager, "hotel")
	if err != nil {
		return nil, fmt.Errorf("failed to get hotel service: %w", err)
	}
	return hotelService.HotelMediaLinks(ctx, req)
}

// Retrieve retrieves an existing hotel booking by record locator
// (HotelRetrieveServicePortType).
func (f *HotelFacade) Retrieve(ctx context.Context, req *hotelxsd.HotelRetrieveReq) (*hotelxsd.HotelRetrieveRsp, error) {
	if f.manager == nil {
		return nil, errors.New("service manager is nil")
	}
	if req == nil {
		return nil, NewValidationError("body", "请求体不能为空")
	}
	ctx, err := normalizeHotelPortReq(ctx, &req.BaseReq)
	if err != nil {
		return nil, err
	}

	hotelService, err := manager.Get[*hotelsvc.HotelService](f.manager, "hotel")
	if err != nil {
		return nil, fmt.Errorf("failed to get hotel service: %w", err)
	}
	return hotelService.HotelRetrieve(ctx, req)
}

// Rules retrieves hotel rate rules (HotelRulesServicePortType).
func (f *HotelFacade) Rules(ctx context.Context, req *hotelxsd.HotelRulesReq) (*hotelxsd.HotelRulesRsp, error) {
	if f.manager == nil {
		return nil, errors.New("service manager is nil")
	}
	if req == nil {
		return nil, NewValidationError("body", "请求体不能为空")
	}
	ctx, err := normalizeHotelPortReq(ctx, &req.BaseReq)
	if err != nil {
		return nil, err
	}

	hotelService, err := manager.Get[*hotelsvc.HotelService](f.manager, "hotel")
	if err != nil {
		return nil, fmt.Errorf("failed to get hotel service: %w", err)
	}
	return hotelService.HotelRules(ctx, req)
}

// UpsellSearch searches for hotel upsell offers
// (HotelUpsellSearchServicePortType).
func (f *HotelFacade) UpsellSearch(ctx context.Context, req *hotelxsd.HotelUpsellDetailsReq) (*hotelxsd.HotelUpsellDetailsRsp, error) {
	if f.manager == nil {
		return nil, errors.New("service manager is nil")
	}
	if req == nil {
		return nil, NewValidationError("body", "请求体不能为空")
	}
	ctx, err := normalizeHotelPortReq(ctx, &req.BaseReq)
	if err != nil {
		return nil, err
	}

	hotelService, err := manager.Get[*hotelsvc.HotelService](f.manager, "hotel")
	if err != nil {
		return nil, fmt.Errorf("failed to get hotel service: %w", err)
	}
	return hotelService.HotelUpsellSearch(ctx, req)
}

// Keywords retrieves hotels by keyword (HotelKeywordsServicePortType).
func (f *HotelFacade) Keywords(ctx context.Context, req *hotelxsd.HotelKeywordReq) (*hotelxsd.HotelKeywordRsp, error) {
	if f.manager == nil {
		return nil, errors.New("service manager is nil")
	}
	if req == nil {
		return nil, NewValidationError("body", "请求体不能为空")
	}
	ctx, err := normalizeHotelPortReq(ctx, &req.BaseReq)
	if err != nil {
		return nil, err
	}

	hotelService, err := manager.Get[*hotelsvc.HotelService](f.manager, "hotel")
	if err != nil {
		return nil, fmt.Errorf("failed to get hotel service: %w", err)
	}
	return hotelService.HotelKeywords(ctx, req)
}

// SuperShopper performs a hotel super shopper search
// (HotelSuperShopperServicePortType).
func (f *HotelFacade) SuperShopper(ctx context.Context, req *hotelxsd.HotelSuperShopperReq) (*hotelxsd.HotelSuperShopperRsp, error) {
	if f.manager == nil {
		return nil, errors.New("service manager is nil")
	}
	if req == nil {
		return nil, NewValidationError("body", "请求体不能为空")
	}
	ctx, err := normalizeHotelPortReq(ctx, &req.BaseReq)
	if err != nil {
		return nil, err
	}

	hotelService, err := manager.Get[*hotelsvc.HotelService](f.manager, "hotel")
	if err != nil {
		return nil, fmt.Errorf("failed to get hotel service: %w", err)
	}
	return hotelService.HotelSuperShopper(ctx, req)
}

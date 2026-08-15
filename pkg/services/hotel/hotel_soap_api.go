package hotel

import (
	"context"
	"encoding/xml"
	"fmt"
	"strings"
	"time"

	common55 "github.com/shuiyihan12/uapi-go/pkg/generated/common55"
	hotelxsd "github.com/shuiyihan12/uapi-go/pkg/generated/hotel"
	"github.com/shuiyihan12/uapi-go/pkg/trace"
	"go.uber.org/zap"
)

const (
	// defaultHotelDetailsPageTimeout is the business-level total time limit for
	// HotelDetails recursive paging: starting from the first request, the whole
	// "first page + subsequent NextResultReference paging" loop must not exceed
	// this value cumulatively; on timeout, paging stops and the data collected
	// so far is returned. It bounds the cumulative duration of multiple SOAP
	// calls and is orthogonal to HTTP transport-level timeouts such as
	// UAPI_REQUEST_TIMEOUT (which bound a single call).
	defaultHotelDetailsPageTimeout = 40 * time.Second
	// defaultHotelDetailsMaxPages caps the total number of pages fetched per
	// Details call (first page included).
	defaultHotelDetailsMaxPages = 20
)

// SearchAvailability issues the HotelSearchAvailability SOAP call and returns
// the strongly typed response of the generated types.
//
// Both request and response use the WSDL-generated hotel package types (via
// BaseSearchReq embedding BaseReq; request-level fields such as TargetBranch
// and BillingPointOfSaleInfo are provided explicitly by the caller in the JSON
// request body, and TraceId is injected as a fallback by the facade layer via
// InjectInfrastructure when the request body leaves it empty). The response is
// written back as-is, with no per-field mapping.
func (s *HotelService) SearchAvailability(ctx context.Context, req *hotelxsd.HotelSearchAvailabilityReq) (*hotelxsd.BaseHotelSearchRsp, error) {
	ctx, _ = trace.Ensure(ctx)

	responseBody, err := s.client.CallWithObservability(ctx, "HotelSearchAvailability", req)
	if err != nil {
		return nil, fmt.Errorf("hotel search availability failed: %w", err)
	}

	var response hotelxsd.BaseHotelSearchRsp
	if err := xml.Unmarshal(responseBody, &response); err != nil {
		return nil, fmt.Errorf("failed to parse hotel search availability response: %v", err)
	}
	return &response, nil
}

// Details issues the HotelDetails SOAP call and returns the strongly typed
// response of the generated types, automatically paging for the remaining
// rates within the timeout and page-count limits.
func (s *HotelService) Details(ctx context.Context, req *hotelxsd.HotelDetailsReq) (*hotelxsd.HotelDetailsRsp, error) {
	ctx, _ = trace.Ensure(ctx)

	responseBody, err := s.client.CallWithObservability(ctx, "HotelDetails", req)
	if err != nil {
		return nil, fmt.Errorf("hotel details failed: %w", err)
	}

	var response hotelxsd.HotelDetailsRsp
	if err := xml.Unmarshal(responseBody, &response); err != nil {
		return nil, fmt.Errorf("failed to parse hotel details response: %v", err)
	}

	if s.logger != nil {
		s.logger.Raw(fmt.Sprintf("[GDS RESPONSE] operation=HotelDetails trace_id=%s page=1 bytes=%d\n%s",
			trace.ID(ctx), len(responseBody), string(responseBody)))
	}

	if err := s.loadRemainingDetailPages(ctx, req, &response); err != nil {
		return nil, err
	}
	return &response, nil
}

// loadRemainingDetailPages automatically pages for the remaining details within
// the timeout and page-count limits.
//
// During paging, HotelRateDetail entries are deduplicated by their RatePlanType
// field: the GDS may repeatedly return RatePlanTypes already seen while paging
// (duplicated messages); if a page hits an already-collected RatePlanType, that
// is treated as a duplicate/loop signal and paging exits immediately to avoid an
// infinite loop. The cumulative time limit (defaultHotelDetailsPageTimeout) and
// next-page token deduplication serve as additional fallbacks.
func (s *HotelService) loadRemainingDetailPages(ctx context.Context, req *hotelxsd.HotelDetailsReq, firstRsp *hotelxsd.HotelDetailsRsp) error {
	start := time.Now()
	seenNextRefs := make(map[string]struct{})
	seenRatePlans := collectRatePlanTypes(firstRsp)

	currentRsp := firstRsp
	for page := 2; page <= defaultHotelDetailsMaxPages; page++ {
		nextRef := currentRsp.NextResultReference
		if nextRef == nil || strings.TrimSpace(string(nextRef.Value)) == "" {
			return nil
		}
		if time.Since(start) >= defaultHotelDetailsPageTimeout {
			if s.logger != nil {
				s.logger.WithContext(ctx).Warn("hotel details paging aborted by timeout",
					zap.Duration("elapsed", time.Since(start)),
					zap.Duration("limit", defaultHotelDetailsPageTimeout),
					zap.Int("page", page))
			}
			return nil
		}

		nextKey := nextResultReferenceKey(nextRef)
		if _, exists := seenNextRefs[nextKey]; exists {
			return nil
		}
		seenNextRefs[nextKey] = struct{}{}

		// The generated request and response types share common55.NextResultReference,
		// so the paging token is passed straight through.
		req.NextResultReference = nextRef
		responseBody, err := s.client.CallWithObservability(ctx, "HotelDetails", req)
		if err != nil {
			return fmt.Errorf("hotel details page %d failed: %v", page, err)
		}

		var pageRsp hotelxsd.HotelDetailsRsp
		if err := xml.Unmarshal(responseBody, &pageRsp); err != nil {
			return fmt.Errorf("failed to parse hotel details page %d response: %v", page, err)
		}

		// Log each page's raw SOAP XML to help troubleshoot Travelport rate paging issues.
		if s.logger != nil {
			s.logger.Raw(fmt.Sprintf("[GDS RESPONSE] operation=HotelDetails trace_id=%s page=%d bytes=%d\n%s",
				trace.ID(ctx), page, len(responseBody), string(responseBody)))
		}

		// Deduplication: a page containing an already-seen RatePlanType means the
		// GDS is returning duplicated data in a loop; exit immediately.
		if duplicatesRatePlanType(&pageRsp, seenRatePlans) {
			if s.logger != nil {
				s.logger.WithContext(ctx).Warn("hotel details paging stopped: duplicate RatePlanType detected",
					zap.Int("page", page))
			}
			return nil
		}

		// Merge this page's rate details into the first response and extend the
		// set of seen RatePlanTypes.
		mergeRateDetails(firstRsp, &pageRsp, seenRatePlans)

		currentRsp = &pageRsp
	}

	return nil
}

// collectRatePlanTypes collects all RatePlanType values already present in the
// response, used for paging deduplication checks.
func collectRatePlanTypes(rsp *hotelxsd.HotelDetailsRsp) map[string]struct{} {
	seen := make(map[string]struct{})
	if rsp == nil || rsp.RequestedHotelDetails == nil {
		return seen
	}
	for i := range rsp.RequestedHotelDetails.HotelRateDetail {
		if t := string(rsp.RequestedHotelDetails.HotelRateDetail[i].RatePlanType); t != "" {
			seen[t] = struct{}{}
		}
	}
	return seen
}

// duplicatesRatePlanType reports whether the page contains an already-collected
// RatePlanType (the GDS looping signal).
func duplicatesRatePlanType(rsp *hotelxsd.HotelDetailsRsp, seen map[string]struct{}) bool {
	if rsp == nil || rsp.RequestedHotelDetails == nil {
		return false
	}
	for i := range rsp.RequestedHotelDetails.HotelRateDetail {
		if t := string(rsp.RequestedHotelDetails.HotelRateDetail[i].RatePlanType); t != "" {
			if _, exists := seen[t]; exists {
				return true
			}
		}
	}
	return false
}

// mergeRateDetails appends the page's rate details to the first response and
// merges the RatePlanTypes seen on this page into the seen set.
func mergeRateDetails(first *hotelxsd.HotelDetailsRsp, page *hotelxsd.HotelDetailsRsp, seen map[string]struct{}) {
	if first.RequestedHotelDetails == nil || page.RequestedHotelDetails == nil {
		return
	}
	first.RequestedHotelDetails.HotelRateDetail = append(
		first.RequestedHotelDetails.HotelRateDetail,
		page.RequestedHotelDetails.HotelRateDetail...,
	)
	for i := range page.RequestedHotelDetails.HotelRateDetail {
		if t := string(page.RequestedHotelDetails.HotelRateDetail[i].RatePlanType); t != "" {
			seen[t] = struct{}{}
		}
	}
}

// nextResultReferenceKey builds the deduplication key for a next-page reference
// (provider code + token value).
func nextResultReferenceKey(ref *common55.NextResultReference) string {
	if ref == nil {
		return ""
	}
	pc := ""
	if ref.ProviderCode != nil {
		pc = string(*ref.ProviderCode)
	}
	return strings.TrimSpace(pc) + "|" + strings.TrimSpace(string(ref.Value))
}

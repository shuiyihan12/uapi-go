package system

import (
	"context"
	"encoding/json"
	"encoding/xml"
	"strings"
	"testing"

	common32 "github.com/shuiyihan12/uapi-go/pkg/generated/common32"
	systemxsd "github.com/shuiyihan12/uapi-go/pkg/generated/system"
	"github.com/shuiyihan12/uapi-go/pkg/trace"
)

// TestSystemInfoReqDecodeCamelCase verifies that the external API decodes with
// the camelCase contract (consistent with hotel), and that
// billingPointOfSaleInfo and targetBranch explicitly passed by the caller are
// correctly received, no longer failing with unknown field.
func TestSystemInfoReqDecodeCamelCase(t *testing.T) {
	body := `{"billingPointOfSaleInfo":{"originApplication":"UAPI","cIDBNumber":12345},"targetBranch":"P000000"}`
	var req systemxsd.SystemInfoReq
	if err := json.Unmarshal([]byte(body), &req); err != nil {
		t.Fatalf("decode camelCase failed: %v", err)
	}
	if req.BillingPointOfSaleInfo.OriginApplication != "UAPI" {
		t.Fatalf("billing not decoded: %+v", req.BillingPointOfSaleInfo)
	}
	if req.BillingPointOfSaleInfo.CIDBNumber == nil || *req.BillingPointOfSaleInfo.CIDBNumber != 12345 {
		t.Fatalf("cIDBNumber not decoded: %v", req.BillingPointOfSaleInfo.CIDBNumber)
	}
	if req.TargetBranch == nil || string(*req.TargetBranch) != "P000000" {
		t.Fatalf("targetBranch not decoded: %v", req.TargetBranch)
	}
}

// TestSystemInfoReqInjectAndMarshal verifies that the server side injects only
// the trace identifier (TraceId), that authorization/business fields such as
// the billing point of sale and TargetBranch are provided by the caller (API
// user), and that the serialized SOAP message faithfully reflects the caller's
// values.
func TestSystemInfoReqInjectAndMarshal(t *testing.T) {
	ctx, _ := trace.Ensure(context.Background())
	branch := common32.TypeBranchCode("P000000")
	req := systemxsd.SystemInfoReq{}
	req.TargetBranch = &branch
	req.BillingPointOfSaleInfo = common32.BillingPointOfSaleInfo{
		OriginApplication: "UAPI",
		CIDBNumber:        int64Ptr(12345),
	}
	prepareReq(ctx, &req)

	out, err := xml.Marshal(req)
	if err != nil {
		t.Fatalf("marshal failed: %v", err)
	}
	t.Logf("marshaled SystemInfoReq:\n%s", string(out))
	xmlStr := string(out)

	// The root element uses the system_v32_0 namespace (only the URI must match;
	// the GDS does not care about the prefix name).
	if !strings.Contains(xmlStr, `xmlns="http://www.travelport.com/schema/system_v32_0"`) {
		t.Errorf("missing system namespace on root element:\n%s", xmlStr)
	}
	// The billing point of sale and TargetBranch come from the caller (not
	// injected by the server side) and must appear faithfully.
	for _, want := range []string{
		`BillingPointOfSaleInfo xmlns="http://www.travelport.com/schema/common_v32_0"`,
		`OriginApplication="UAPI"`,
		`CIDBNumber="12345"`,
		`TargetBranch="P000000"`,
	} {
		if !strings.Contains(xmlStr, want) {
			t.Errorf("marshaled XML missing %q:\n%s", want, xmlStr)
		}
	}
	// The trace identifier is injected into an attribute (unqualified; it must
	// not carry a namespace prefix).
	if !strings.Contains(xmlStr, `TraceId="`) {
		t.Errorf("traceId not injected:\n%s", xmlStr)
	}
	if strings.Contains(xmlStr, `system_v32_0:TraceId`) || strings.Contains(xmlStr, `system_v32_0:CacheName`) {
		t.Errorf("attribute must NOT be namespace-qualified:\n%s", xmlStr)
	}
}

// TestExternalCacheAccessReqAttrsUnqualified verifies that for requests with
// attributes (ExternalCacheAccessReq), attributes such as CacheName are
// unqualified when serialized; otherwise the GDS is guaranteed to reject them.
func TestExternalCacheAccessReqAttrsUnqualified(t *testing.T) {
	ctx, _ := trace.Ensure(context.Background())
	req := systemxsd.ExternalCacheAccessReq{
		CacheName:     "MyCache",
		RetrieveEntry: []systemxsd.RetrieveEntryInline{{Key: "k1"}},
	}
	prepareReq(ctx, &req)

	out, err := xml.Marshal(req)
	if err != nil {
		t.Fatalf("marshal failed: %v", err)
	}
	t.Logf("marshaled ExternalCacheAccessReq:\n%s", string(out))
	xmlStr := string(out)
	for _, want := range []string{`CacheName="MyCache"`, `RetrieveEntry`, `Key="k1"`} {
		if !strings.Contains(xmlStr, want) {
			t.Errorf("marshaled XML missing %q:\n%s", want, xmlStr)
		}
	}
	if strings.Contains(xmlStr, `system_v32_0:CacheName`) {
		t.Errorf("attribute CacheName must NOT be namespace-qualified:\n%s", xmlStr)
	}
}

func int64Ptr(i int64) *int64 { return &i }

package util

import (
	"context"
	"encoding/xml"
	"strings"
	"testing"

	common55 "github.com/shuiyihan12/uapi-go/pkg/generated/common55"
	utilxsd "github.com/shuiyihan12/uapi-go/pkg/generated/util"
	"github.com/shuiyihan12/uapi-go/pkg/trace"
)

// TestCalculateTaxReqInjectAndMarshal verifies that util requests get only the
// trace identifier injected by prepareReq before sending, that authorization/
// business fields such as the billing point of sale are provided by the caller
// (API user), and that attributes are unqualified (otherwise the GDS is
// guaranteed to reject them).
func TestCalculateTaxReqInjectAndMarshal(t *testing.T) {
	ctx, _ := trace.Ensure(context.Background())
	req := utilxsd.CalculateTaxReq{}
	req.BillingPointOfSaleInfo = common55.BillingPointOfSaleInfo{
		OriginApplication: "UAPI",
		CIDBNumber:        int64Ptr(12345),
	}
	prepareReq(ctx, &req)

	out, err := xml.Marshal(req)
	if err != nil {
		t.Fatalf("marshal failed: %v", err)
	}
	t.Logf("marshaled CalculateTaxReq:\n%s", string(out))
	xmlStr := string(out)

	for _, want := range []string{
		`xmlns="http://www.travelport.com/schema/util_v55_0"`,
		`BillingPointOfSaleInfo xmlns="http://www.travelport.com/schema/common_v55_0"`,
		`OriginApplication="UAPI"`,
		`CIDBNumber="12345"`,
		`TraceId="`,
	} {
		if !strings.Contains(xmlStr, want) {
			t.Errorf("marshaled XML missing %q:\n%s", want, xmlStr)
		}
	}
	// No attribute may carry a namespace prefix (this was previously the root
	// cause of GDS rejections and lost values during response parsing).
	if strings.Contains(xmlStr, `util_v55_0:TraceId`) || strings.Contains(xmlStr, `common_v55_0:CIDBNumber`) {
		t.Errorf("attribute must NOT be namespace-qualified:\n%s", xmlStr)
	}
}

func int64Ptr(i int64) *int64 { return &i }

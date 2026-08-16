package sdk_test

import (
	"context"
	"fmt"
	"log"
	"time"

	common55 "github.com/shuiyihan12/uapi-go/pkg/generated/common55"
	hotelxsd "github.com/shuiyihan12/uapi-go/pkg/generated/hotel"
	"github.com/shuiyihan12/uapi-go/pkg/logging"
	"github.com/shuiyihan12/uapi-go/pkg/requestctx"
	"github.com/shuiyihan12/uapi-go/sdk"
)

// ExampleNew shows the full SDK flow: assemble a client, pick a domain
// service, and issue a request with per-call credentials. Construction is
// offline; only the actual SOAP call touches the network.
func ExampleNew() {
	client, err := sdk.New(
		sdk.WithEndpoint("https://apac.universal-api.travelport.com/B2BGateway/connect/uAPI"),
		sdk.WithTimeouts(45*time.Second, 90*time.Second, 90*time.Second),
		sdk.WithKeepAlivePool(100, 100),
		sdk.WithLogger(logging.Noop()),
	)
	if err != nil {
		log.Fatal(err)
	}
	defer client.Close()

	hotel, err := client.Hotel()
	if err != nil {
		log.Fatal(err)
	}

	// Credentials are request-level, mirroring the REST gateway: pass the
	// Travelport Authorization value and the region on the call context.
	ctx := context.Background()
	ctx = requestctx.WithAuthorization(ctx, "Basic <credentials>")
	ctx = requestctx.WithRegion(ctx, "apac")

	// Requests are built from the WSDL-generated types (pkg/generated/...).
	branch := common55.TypeBranchCode("P000000")
	req := &hotelxsd.HotelSearchAvailabilityReq{}
	req.TargetBranch = &branch

	// A real call would be: out, err := hotel.SearchAvailability(ctx, req)
	fmt.Println("client ready, request assembled:", hotel != nil && req.TargetBranch != nil)
	// Output: client ready, request assembled: true
}

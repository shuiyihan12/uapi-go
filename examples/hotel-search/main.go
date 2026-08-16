// Command hotel-search demonstrates a full business flow through the SDK:
// build a typed hotel availability request from the WSDL-generated types,
// call the Hotel service with request-level credentials, and walk the
// strongly typed response.
//
// Usage:
//
//	export UAPI_EXAMPLE_AUTHORIZATION='Basic <credentials>'
//	export UAPI_EXAMPLE_TARGET_BRANCH='<your branch code>'
//	go run ./examples/hotel-search [IATA-city-code]   # default NYC
package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"time"

	common55 "github.com/shuiyihan12/uapi-go/pkg/generated/common55"
	hotelxsd "github.com/shuiyihan12/uapi-go/pkg/generated/hotel"
	"github.com/shuiyihan12/uapi-go/pkg/logging"
	"github.com/shuiyihan12/uapi-go/pkg/requestctx"
	"github.com/shuiyihan12/uapi-go/sdk"
)

func main() {
	auth := os.Getenv("UAPI_EXAMPLE_AUTHORIZATION")
	branch := os.Getenv("UAPI_EXAMPLE_TARGET_BRANCH")
	if auth == "" || branch == "" {
		log.Fatal("set UAPI_EXAMPLE_AUTHORIZATION and UAPI_EXAMPLE_TARGET_BRANCH first")
	}
	city := "NYC"
	if len(os.Args) > 1 {
		city = os.Args[1]
	}

	client, err := sdk.New(sdk.WithLogger(logging.Noop()))
	if err != nil {
		log.Fatal(err)
	}
	defer client.Close()

	hotel, err := client.Hotel()
	if err != nil {
		log.Fatal(err)
	}

	// Requests are built from the generated types (pkg/generated/...);
	// dates use the XSD lexical form (YYYY-MM-DD) carried as strings.
	iata := common55.TypeIATACode(city)
	targetBranch := common55.TypeBranchCode(branch)
	req := &hotelxsd.HotelSearchAvailabilityReq{}
	req.TargetBranch = &targetBranch
	req.HotelSearchLocation = &hotelxsd.HotelSearchLocation{
		HotelLocation: &hotelxsd.HotelLocation{Location: &iata},
	}
	req.HotelStay = hotelxsd.HotelStay{
		CheckinDate:  hotelxsd.TypeDate(time.Now().AddDate(0, 0, 30).Format("2006-01-02")),
		CheckoutDate: hotelxsd.TypeDate(time.Now().AddDate(0, 0, 31).Format("2006-01-02")),
	}

	ctx := requestctx.WithAuthorization(context.Background(), auth)
	ctx = requestctx.WithRegion(ctx, "apac")

	rsp, err := hotel.SearchAvailability(ctx, req)
	if err != nil {
		log.Fatalf("search failed: %v", err)
	}

	for i, result := range rsp.HotelSearchResult {
		for _, prop := range result.HotelProperty {
			name := "<unnamed>"
			if prop.Name != nil {
				name = *prop.Name
			}
			fmt.Printf("%d: %s %s/%s — %d rates\n",
				i+1, name, prop.HotelChain, prop.HotelCode, len(result.RateInfo))
		}
	}
}

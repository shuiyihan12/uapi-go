// Command ping is the "hello world" of the uapi-go SDK: it builds a
// client, calls the System service's Ping against the real upstream and
// prints the response. Use it to verify credentials and connectivity.
//
// Usage:
//
//	export UAPI_EXAMPLE_AUTHORIZATION='Basic <credentials>'
//	export UAPI_EXAMPLE_REGION=apac   # optional; americas / apac / emea
//	go run ./examples/ping
package main

import (
	"context"
	"fmt"
	"log"
	"os"

	systemxsd "github.com/shuiyihan12/uapi-go/pkg/generated/system"
	"github.com/shuiyihan12/uapi-go/pkg/logging"
	"github.com/shuiyihan12/uapi-go/pkg/requestctx"
	"github.com/shuiyihan12/uapi-go/sdk"
)

func main() {
	auth := os.Getenv("UAPI_EXAMPLE_AUTHORIZATION")
	if auth == "" {
		log.Fatal("set UAPI_EXAMPLE_AUTHORIZATION='Basic <credentials>' first")
	}
	region := os.Getenv("UAPI_EXAMPLE_REGION")

	client, err := sdk.New(sdk.WithLogger(logging.Noop()))
	if err != nil {
		log.Fatal(err)
	}
	defer client.Close()

	system, err := client.System()
	if err != nil {
		log.Fatal(err)
	}

	// Credentials are request-level, exactly like the REST gateway.
	ctx := requestctx.WithAuthorization(context.Background(), auth)
	ctx = requestctx.WithRegion(ctx, region)

	rsp, err := system.Ping(ctx, &systemxsd.PingReq{})
	if err != nil {
		log.Fatalf("ping failed: %v", err)
	}

	fmt.Println("upstream ping OK")
	if rsp.TransactionId != nil {
		fmt.Println("transactionId:", *rsp.TransactionId)
	}
}

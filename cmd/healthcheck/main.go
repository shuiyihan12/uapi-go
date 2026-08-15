// Command healthcheck is a lightweight helper for container probes: it sends
// a GET request to the local daemon's operations endpoint and treats 2xx as
// healthy (exit code 0), otherwise it exits with code 1.
//
// It exists for two reasons:
//   - the production image is distroless-based (no shell / curl / wget), so
//     Docker HEALTHCHECK needs a statically compiled probe binary;
//   - Kubernetes httpGet probes cannot reference Secrets in headers, while
//     /health requires an Authorization header (passed through to the
//     upstream SystemPing). An exec probe with environment-variable
//     injection (secretRef) is the standard solution.
//
// Environment variables:
//
//	PORT                      target port (default 8080)
//	HEALTHCHECK_PATH          probe path (default /ready; /health requires
//	                          upstream reachability)
//	HEALTHCHECK_AUTHORIZATION optional; forwarded verbatim as the
//	                          Authorization header (used when probing /health)
//	HEALTHCHECK_TIMEOUT       timeout in seconds (default 3)
package main

import (
	"fmt"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"
)

func main() {
	port := envOr("PORT", "8080")
	path := envOr("HEALTHCHECK_PATH", "/ready")
	auth := strings.TrimSpace(os.Getenv("HEALTHCHECK_AUTHORIZATION"))
	timeout := 3 * time.Second
	if v, err := strconv.Atoi(envOr("HEALTHCHECK_TIMEOUT", "3")); err == nil && v > 0 {
		timeout = time.Duration(v) * time.Second
	}

	req, err := http.NewRequest(http.MethodGet, "http://127.0.0.1:"+port+path, nil)
	if err != nil {
		fmt.Fprintf(os.Stderr, "healthcheck: bad request: %v\n", err)
		os.Exit(1)
	}
	if auth != "" {
		req.Header.Set("Authorization", auth)
	}

	client := &http.Client{Timeout: timeout}
	resp, err := client.Do(req)
	if err != nil {
		fmt.Fprintf(os.Stderr, "healthcheck: %s failed: %v\n", path, err)
		os.Exit(1)
	}
	defer resp.Body.Close()

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		fmt.Fprintf(os.Stderr, "healthcheck: %s returned HTTP %s\n", path, resp.Status)
		os.Exit(1)
	}
	os.Exit(0)
}

// envOr reads an environment variable, returning fallback when unset or
// empty.
func envOr(key, fallback string) string {
	if v := strings.TrimSpace(os.Getenv(key)); v != "" {
		return v
	}
	return fallback
}

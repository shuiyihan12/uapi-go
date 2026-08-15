# Development Guide

**[English](./development.md)** | [简体中文](./development.zh-CN.md)

> Audience: engineers who want to **add endpoints, onboard new upstream domains, upgrade WSDL contracts or add business orchestration** on top of this project.
> Read [`architecture.md`](./architecture.md) first for the layering and design decisions, then come back here for the scenario playbooks.
> Route ↔ SOAP operation mapping: [`routing.md`](./routing.md); GDS terminology: [`glossary.md`](./glossary.md).

---

## 0. Understand the code structure in five minutes

```text
HTTP request
  └─ pkg/api/            route registration + JSON codec + error → HTTP status
      └─ pkg/usecase/    business use cases (facades): validation, orchestration; pure pass-through domains just forward
          └─ pkg/manager/  lazy factory for service instances (generic Get[T])
              └─ pkg/services/<domain>/  PortType methods → SOAP payloads
                  └─ pkg/client/  envelope assembly, auth pass-through, metrics, fault parsing
                      └─ Travelport UAPI (SOAP/XML)
```

Three "universal facilities" ship with the project — reuse them in new code instead of reinventing:

| Facility | Location | What it does |
|---|---|---|
| `registerPortHandler[TReq,TRsp]` | `pkg/api/handler.go` | registers a POST route in one line: strict JSON decoding (unknown fields rejected), request-level auth/region/trace injection, unified timeouts and error mapping |
| `CallPortType[T]` / `callPort[T]` | `pkg/client/enterprise.go` | makes the SOAP call and deserializes the response into the strongly typed `*T`, with metrics and leveled logging for free |
| `manager.Get[T]` | `pkg/manager/service_manager.go` | type-safely fetches (and caches) a domain's service instance |

---

## 1. Scenario A: expose a PortType operation the upstream already provides (most common)

Using the Air domain as the example — one operation touches four places from HTTP to SOAP. Reference implementation: `AirExchange` (`pkg/services/air/service.go`).

> **Generated-package import convention**: the services/usecase layers import both the business package (`pkg/services/air`) and the generated package (`pkg/generated/air`) whose package names clash, so generated packages are imported with an `xsd`-suffixed alias (`airxsd "github.com/shuiyihan12/uapi-go/pkg/generated/air"`; hotel uses the existing `hotelxsd`); the `commonNN` family has no clash and stays bare.

### Step 1: add the PortType method in the services layer

Add to both the `AirServicePort` interface and its implementation in `pkg/services/air/service.go`:

```go
// AirXxx corresponds to the Air service's AirXxxReq operation.
AirXxx(ctx context.Context, req *airxsd.AirXxxReq) (*airxsd.AirXxxRsp, error)

func (s *AirService) AirXxx(ctx context.Context, req *airxsd.AirXxxReq) (*airxsd.AirXxxRsp, error) {
	ctx = prepareReq(ctx, req) // inject TraceId (best effort: written when the request implements InjectInfrastructure)
	return callPort[airxsd.AirXxxRsp](s.client, ctx, "AirXxx", req)
}
```

Use the generated types from `pkg/generated/air` for request/response — **never hand-write request models**.

### Step 2: add the facade method in the usecase layer

`pkg/usecase/air_facade.go` (pure pass-through template):

```go
// AirXxx corresponds to AirServicePort.AirXxx.
func (f *AirFacade) AirXxx(ctx context.Context, req *airxsd.AirXxxReq) (*airxsd.AirXxxRsp, error) {
	svc, err := f.getService()
	if err != nil {
		return nil, err
	}
	ctx, _ = trace.Ensure(ctx)
	return svc.AirXxx(ctx, req)
}
```

> Pass-through facades are a deliberate "collapsible layer": no business logic today; when validation / orchestration is needed later it grows in place, without cross-layer changes.

### Step 3: register the route in the api layer

Add one line to `RegisterRoutes` in `pkg/api/air_handler.go`:

```go
registerPortHandler(mux, apiBasePath+"/air/air-xxx", f.AirXxx)
```

Naming rule: routes are kebab-case (`air-xxx`), mapping 1:1 to the method name (`AirXxx`).

### Step 4: docs and tests

- `docs/routing.md`: add a row to the domain's table (route | SOAP operation | upstream service | notes) and update the domain and total route counts
- For non-trivial field mappings or validation, add table-driven tests (see the Marshal assertions in `pkg/services/system/system_test.go`)
- Commit after `./scripts/build.sh all` is fully green

---

## 2. Scenario B: onboard a brand-new upstream domain

Onboarding a new domain `foo` (WSDL delivered) concentrates changes in six places:

1. **Archive the contract**: put the WSDL/XSD under `wsdl/foo_vNN_0/` and confirm `include`/`import` references are complete
2. **Generate structs**: update the input manifest in `tools/airxsdgen` (if new schema files are needed) → `./scripts/build.sh wsdl` → produces `pkg/generated/foo/` (single-version domains get versionless package names; only the common family keeps versions)
3. **Services layer**: create `pkg/services/foo/service.go` — define the `FooServicePort` interface (1:1 with the WSDL's *PortType) plus the `FooService` implementation (constructs a `client.EnterpriseSOAPClient`; each method is `prepareReq` + `callPort`), modeled on `pkg/services/rail/service.go` (small-domain template)
4. **Manager registration**: two spots in `pkg/manager/service_manager.go` — add a factory entry to `buildServiceFactories`, and add `case "foo": return "FooService", nil` to `serviceSuffix` (unknown keys error explicitly, so a missed entry surfaces on first call)
5. **Usecase + api layers**: create `pkg/usecase/foo_facade.go` and `pkg/api/foo_handler.go` (templates as in scenario A), and wire them in `cmd/daemon/main.go`: `NewFooFacade` → `NewFooHandler(...).RegisterRoutes(mux)`
6. **Docs**: the README capability table, a new domain section in `docs/routing.md`, and `docs/architecture.md` §2.1's layer table when relevant

---

## 3. Scenario C: upgrade the upstream contract (WSDL/XSD)

Full process in [`architecture.md` §6](./architecture.md). Essentials:

```bash
# After archiving the new contract into wsdl/:
./scripts/build.sh wsdl   # regenerate pkg/generated
./scripts/build.sh all    # build + test + lint
git diff --stat pkg/generated/   # manually review the generated diff (added/removed fields, type changes)
```

- `pkg/generated/` is a generated asset — **review the diff read-only; never hand-edit**
- Watch for recurrence of the two historical generator defects: dropped fields from cross-namespace `ref`s, and `xs:date` serialization anomalies (see `architecture.md` §3.4)
- Upstream field removals / type changes surface directly as compile errors — exactly the "contract moves forward to compile time" design intent; just fix the adapter layers accordingly

---

## 4. Scenario D: add business orchestration (validation, paging, aggregation)

When pass-through is not enough (multi-step SOAP calls, input validation, friendlier output), grow the logic in the **usecase layer**. Ready-made example: `pkg/usecase/hotel_facade.go`.

- **Input validation**: see `validateHotelStay` (date sanity) and `normalizeHotelPortReq` (required-field normalization); failed validation returns `*usecase.ValidationError`, which the API layer maps to 400 `INVALID_REQUEST` automatically
- **Multi-call / paging**: see the `NextResultReference` paging loop in `HotelDetails` — mind the two timeout classes: a single SOAP call is bounded by `UAPI_REQUEST_TIMEOUT`, while cumulative multi-call time is bounded by a business-level constant (e.g. `defaultHotelDetailsPageTimeout = 40s`, `defaultHotelDetailsMaxPages = 20`)
- **Output shape**: wrap orchestration results in a hand-written Output struct (e.g. `HotelSearchAvailabilityOutput`, holding a single `Response` field) to keep the top-level JSON shape stable — avoid `map[string]interface{}` pass-through
- **Note**: orchestrated endpoints can still use `registerPortHandler` — as long as the facade method's signature is `func(ctx, *Req) (*Rsp, error)` (all 8 hotel routes register this way), no hand-written handler is needed

---

## 5. Cross-cutting conventions (read before touching any code)

| Convention | Notes |
|---|---|
| Auth passes through everything | the gateway holds no Travelport credentials; `Authorization` / `X-UAPI-Region` are carried per request by callers and flow from HTTP headers to SOAP calls via `pkg/requestctx`. Same for `/health` (the prober carries credentials) |
| No retries | GDS write operations (ticketing/booking) are not idempotent; every failure surfaces directly and the caller decides compensation |
| Error model | upstream business/system errors surface as `*client.SOAPFaultError`; `pkg/api.writeError` maps them uniformly to 400/422/502/504 by `ErrorInfo/Type`. Never swallow errors or stringify them in intermediate layers |
| trace_id end-to-end | `trace.Ensure(ctx)` generates a UUID automatically and threads it through HTTP headers, logs and SOAP payloads; do not invent your own tracking IDs in new code |
| Generated code is untouchable | `pkg/generated/` regenerates only via `./scripts/build.sh wsdl` |
| Strict JSON | request decoding uses `DisallowUnknownFields`; new endpoint request fields are snake_case, aligned with the generated types' json tags |
| Comment style | English doc comments; exported symbols must document "what and why" |

---

## 6. Testing and debugging

### Testing

```bash
go test ./...                          # all unit tests (offline by default; CI-safe)
UAPI_INTEGRATION=1 \
UAPI_TEST_AUTHORIZATION="Basic xxx" \
go test ./test/... -run TestServiceManager -v   # real-upstream integration path
```

- Unit tests live next to the package under test (e.g. `pkg/services/system/system_test.go`); integration tests live in `test/`
- Cases needing real credentials are always gated by `UAPI_INTEGRATION`
- To fake upstream responses: `client.SOAPConfig.Transport` accepts an injected `http.RoundTripper` (see `pkg/client/soap_test.go`)

### Debugging

1. **Get the trace_id**: any response's `X-Trace-Id` header (or pass your own `X-Trace-Id`)
2. **Inspect raw payloads**: filter logs by trace_id for `[GDS REQUEST]` / `[GDS RESPONSE]` to check SOAP payloads field by field
3. **Check metrics**: `GET /metrics` (Prometheus format; `uapi_requests_total` by service/operation/status)
4. **Run locally**:

```bash
set -a && source .env && set +a
go run ./cmd/daemon          # default :8080; business shares the port with /health, /metrics

curl -X POST http://localhost:8080/api/system/ping \
  -H "Authorization: Basic xxxx" \
  -H "X-UAPI-Region: apac" \
  -H "Content-Type: application/json" -d '{}'
```

5. **Health-check troubleshooting**: `curl -H "Authorization: Basic xxxx" http://localhost:8080/health` — without credentials the upstream rejects and 503 is returned by design (the gateway itself holds no credentials to verify).

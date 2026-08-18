# Changelog

All notable changes to this project are documented in this file.

The version scheme follows `v0.WSDL.PATCH`:
- `WSDL` tracks the Travelport UAPI schema generation (e.g. `55` for `v55_0`).
- `PATCH` is bumped for SDK-level fixes and improvements that do not change the generated WSDL surface.

Entries are listed from earliest to latest.

---

## [v0.55.1] — 2026-08-16

Initial commit. First version of the Go SDK / REST gateway for Travelport UAPI.

- Project scaffold: `cmd/` entrypoints (`daemon` REST+Admin server, `generator`, `generate`), `pkg/` layers (`api`, `usecase`, `client`, `manager`, `services`, `requestctx`, `trace`), `pkg/generated` WSDL/XSD-generated code, `internal/` infrastructure, `tools/airxsdgen` generator, `wsdl/` upstream inputs, and `scripts/` build tooling.
- Container build assets (`Dockerfile`, `docker-compose.yml`, `deploy/k8s/`, `cmd/healthcheck/`) using a distroless nonroot runtime image.
- Documentation and contributor guidance (`AGENTS.md`, `README`, `docs/` architecture/routing/glossary).

## [v0.55.2] — 2026-08-16

Go SDK surface, drop stray binary, fix docker tags.

- **chore(version)**: adopt the `v0.WSDL.PATCH` module version scheme.
- **refactor(logging)**: promote the logger to the public `pkg/logging` package.
- **refactor(client)**: make SOAP metrics pluggable via `SOAPConfig.Metrics`.
- **feat(sdk)**: add the `sdk` entry package exposing a unified client surface.
- **ci**: guard the SDK import surface; **docs**: version discipline.
- **docs(sdk)**: add runnable examples (`ping`, `hotel-search`).
- **chore(sdk)**: polish options, CI guard, and version discipline.
- Drop a stray binary artifact and fix Docker image tags in the release pipeline.

## [v0.55.3] — 2026-08-19

SOAP request envelope cleanup and namespace hoisting.

### Changed
- **SOAP envelope prefix**: the envelope now uses the `soapenv` prefix (`xmlns:soapenv="http://schemas.xmlsoap.org/soap/envelope/"`) to match the official Travelport UAPI samples.
- **Empty Header dropped**: the empty `<soapenv:Header/>` is omitted. It carries no content (auth rides on the HTTP `Authorization` header and the trace ID on the body's `TraceId` attribute), so dropping it saves a few bytes with no behavioral change.
- **Namespace hoisting**: request-body namespaces are now hoisted to a single declaration on the request root element as readable prefixes — e.g. `xmlns:hotel="http://www.travelport.com/schema/hotel_v55_0"`, `xmlns:common="http://www.travelport.com/schema/common_v55_0"` — instead of `encoding/xml` repeatedly emitting a default `xmlns="..."` on every element. The output now mirrors the official sample's structure (one declaration, prefixed children) and is smaller.

### Notes
- The hoist transform is generic and URI-driven: it works off whatever namespace URIs appear in the body, so it applies unchanged to every UAPI service (air/rail/vehicle/hotel/universal/...) and every schema version. Namespace prefixes are derived by stripping the version suffix from the URI tail segment (`hotel_v55_0` → `hotel`, `common_v55_0` → `common`), so a future `v56_0` regen is handled automatically without code changes.
- No changes to the response parsing path, the generated code under `pkg/generated/`, or any per-service business logic.
- Covered by `pkg/client/soap_test.go` (`TestBuildEnvelopeNamespaceHoist`), which locks the envelope output shape.

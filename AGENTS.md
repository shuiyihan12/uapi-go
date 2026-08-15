# Repository Guidelines

## Project Structure & Module Organization
This repository is a Go SDK / REST gateway for Travelport UAPI.
- `cmd/`: runnable entrypoints (`daemon` is the REST+Admin server; `generator` scaffolds service packages; `generate` wraps gowsdl).
- `pkg/api/`: HTTP handler layer (generic `registerPortHandler` route registration, unified error responses).
- `pkg/usecase/`: business facade layer (one facade per domain; hotel/universal contain real orchestration).
- `pkg/client/`: SOAP transport (envelope build/parse, auth pass-through, observability client).
- `pkg/manager/`: service manager (typed lazy factory registry, health check).
- `pkg/services/*`: domain SOAP service adapters (air/rail/vehicle/hotel/universal/...).
- `pkg/requestctx/`, `pkg/trace/`: per-request auth/region context and trace_id propagation.
- `pkg/generated/`: WSDL/XSD-generated code. Treat as generated artifacts.
- `internal/`: internal-only infrastructure (`logging`, `metrics`).
- `tools/airxsdgen/`: XSD→Go struct generator (regenerates `pkg/generated`).
- `wsdl/`: upstream WSDL/XSD inputs.
- `scripts/`: build/generation tooling.
- `test/`: current integration-style tests (for service manager behavior).
- `docs/`: architecture, routing map, glossary.
- `Dockerfile`, `docker-compose.yml`, `deploy/k8s/`, `cmd/healthcheck/`: container build & deployment assets (distroless nonroot runtime image).

## Build, Test, and Development Commands
- `./scripts/build.sh deps`: tidy and download Go dependencies.
- `./scripts/build.sh wsdl`: regenerate all WSDL client code.
- `./scripts/build.sh build`: build binaries into `bin/`.
- `./scripts/build.sh test`: run test suite (`go test ./test/...`).
- `./scripts/build.sh lint`: run formatting check + `go vet`.
- `./scripts/build.sh all`: full pipeline (deps, generation, build, test, lint).
- `go test ./test/manager_test.go -v`: run focused tests.
- `go test ./test/manager_test.go -bench=.`: run benchmarks.

## Coding Style & Naming Conventions
- Follow Go defaults: tabs for indentation, `gofmt` formatting, `go vet` clean.
- Keep exported names in `PascalCase`, internal helpers in `camelCase`.
- Use package names in lower case and domain-aligned (`hotel`, `universal`, `util`, `system`).
- Do not hand-edit files under `pkg/generated/`; regenerate instead.

## Testing Guidelines
- Primary framework: Go `testing` package.
- Place tests in `*_test.go`; prefer table-driven tests for service logic.
- Name tests as `TestXxx` and benchmarks as `BenchmarkXxx`.
- Before opening a PR, run at least `./scripts/build.sh test` and `./scripts/build.sh lint`.

## Commit & Pull Request Guidelines
- Follow the existing conventional-commit style (e.g. `feat(trace): propagate caller X-Trace-Id header`); keep the subject short and imperative.
- Keep commits scoped to one concern (generation, service logic, docs, or tooling).
- PRs should include: purpose, key changes, test evidence (commands/output), and related issue/task.
- For behavior changes, include sample request/response snippets or logs when relevant.

## Security & Configuration Tips
- Never commit real credentials or production tokens in code, tests, or docs.
- Prefer environment variables or local config overrides for secrets.
- Validate daemon flags and `UAPI_*` environment variables (see `.env.example`) before running against non-test environments.

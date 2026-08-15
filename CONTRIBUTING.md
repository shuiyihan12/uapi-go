# Contributing Guide

**[English](./CONTRIBUTING.md)** | [简体中文](./CONTRIBUTING.zh-CN.md)

Thanks for considering a contribution to uapi-go! This is a Go project wrapping Travelport Universal API (SOAP/XML) as a REST/JSON gateway. This document covers the development environment, workflow and hard rules; for architecture and design decisions see [`docs/architecture.md`](./docs/architecture.md), and for step-by-step contributor tutorials see [`docs/development.md`](./docs/development.md).

## 1. Environment

- Go `1.25.0`+ (as pinned in `go.mod`)
- Public internet access (pulling Go modules; integration tests additionally need a reachable Travelport test environment)
- Docker not required (unless you want containerized deployment)

## 2. Development workflow

```bash
# 1. Clone and install dependencies
git clone <your-fork-url> && cd uapi-go
./scripts/build.sh deps

# 2. Iterate (as needed)
./scripts/build.sh build     # build
./scripts/build.sh test      # test
./scripts/build.sh lint      # gofmt check + go vet

# 3. Full verification before opening a PR (generation + build + test + lint)
./scripts/build.sh all
```

PR checklist:

- [ ] `./scripts/build.sh all` fully green
- [ ] Behavior changes come with test evidence (commands and output) or sample requests/responses
- [ ] New endpoints updated `docs/routing.md` (the route ↔ SOAP operation table)
- [ ] `pkg/generated/` contains only changes produced by `./scripts/build.sh wsdl` (never hand-edited)

## 3. Code style

- Follow Go defaults: tab indentation, `gofmt` formatting, `go vet` zero warnings
- Exported identifiers in `PascalCase`, package-local helpers in `camelCase`; package names lowercase and domain-aligned (`hotel`, `universal`, `util`, ...)
- Comments in English; exported symbols must carry doc comments stating "what and why", not restating the code
- One-way layering: `pkg/api → pkg/usecase → pkg/services → pkg/client`; no reverse imports; prefer reusing the existing generic facilities (`registerPortHandler`, `CallPortType[T]`, `manager.Get[T]`)
- Error handling: upstream business errors surface as `*client.SOAPFaultError` and `pkg/api.writeError` maps HTTP status codes uniformly; never swallow errors or alter their types in intermediate layers

## 4. Testing rules

- Framework: the standard library `testing`; tests live in `*_test.go`, table-driven preferred for service logic
- Naming: `TestXxx` / `BenchmarkXxx`
- Cases needing the real upstream must be gated by `UAPI_INTEGRATION=1` (credentials via `UAPI_TEST_AUTHORIZATION`); skipped in default CI
- Security red line: **no real credentials or production tokens in code, tests, log samples or docs — ever**

## 5. Commit messages

Short imperative lines that fit on one screen, e.g.:

```
Add hotel media-links port passthrough
Fix TLS verify default in service manager
Regenerate WSDL client code for v55 contracts
```

One commit does one thing (generation, service logic, docs and tooling are committed separately).

## 6. Contract (WSDL/XSD) upgrades

Upstream contract changes must go through build-time regeneration; see [`docs/architecture.md` §6](./docs/architecture.md) for the process: archive under `wsdl/` → update the generation manifest → `./scripts/build.sh wsdl` → `all` → manually review the `pkg/generated` diff.

## 7. License

- This project is open-sourced under the [Apache License 2.0](./LICENSE); **submitting a PR means you license your contribution under Apache-2.0** (no separate CLA)
- New source files are encouraged to carry the SPDX header: `// SPDX-License-Identifier: Apache-2.0`
- The upstream contract artifacts under `wsdl/` are copyrighted by Travelport — do not modify or redistribute them beyond the license scope

## 8. Code of conduct

- Keep discussions technical; critique ideas, not people
- Before asking, check [`docs/glossary.md`](./docs/glossary.md) (GDS glossary) and [`docs/architecture.md`](./docs/architecture.md)

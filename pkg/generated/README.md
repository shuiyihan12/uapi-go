# Generated Travelport UAPI XSD structs

This directory contains Go structs generated from the Travelport UAPI WSDL/XSD
schemas under `wsdl/`. They are **generated artifacts — do not edit by hand**.
Regenerate with:

```
./scripts/generate-air-xsd.sh
```

## Layout

Every XSD schema **namespace** becomes its own Go package (directly under this
directory, e.g. `air/`, `common55/`). Enumeration (`xs:enumeration`)
simpleTypes are emitted into a per-domain `enums` sub-package so the closed-set
value types do not bloat the struct packages.

**Package names carry no version suffix** for single-version domains
(`air_v55_0` → package `air`), so contract upgrades regenerate in place without
touching import paths. The only exception is the `commonNN` family: legacy
domains pin different common versions, so multiple commons genuinely coexist
and keep their version numbers (see docs/architecture.md ADR-004).

| Directory                  | XSD namespace (schema)              | Notes                                            |
|----------------------------|-------------------------------------|--------------------------------------------------|
| `common55/`                | `.../schema/common_v55_0`           | shared base types (v55 family; versioned by design) |
| `common32/` … `common38/`  | `.../schema/common_v3X_0`           | legacy: multiple common versions genuinely coexist |
| `air/`, `rail/`, `universal/` | `.../schema/air_v55_0` …       | merged into one package (mutually recursive namespaces) |
| `hotel/`, `vehicle/`, `system/`, … | per-domain namespaces       | one package per domain                           |

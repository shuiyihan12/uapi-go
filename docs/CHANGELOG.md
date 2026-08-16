# Changelog

**[English](./CHANGELOG.md)** | [简体中文](./CHANGELOG.zh-CN.md)

## Versioning scheme

- **Module major version is pinned at `0`** (Go only requires a `/vN` module path suffix when major ≥ 2, so the bare import path `github.com/shuiyihan12/uapi-go` stays valid). Tag format is `v0.WSDL.PATCH`, where `WSDL` is the Travelport WSDL contract version (the `wsdl/` folder suffix, e.g. `55` for the v55_0 contract) and `PATCH` is our own fix increment — e.g. `v0.55.1` = released 2026-08.
- **Only a git tag counts as a release**; untagged commits are development state (the daemon shows `dev`).
- **Patch versions** increment the `PATCH` segment: e.g. `v0.55.2` fixes anomalies in `v0.55.1`, with no new features. Patch releases carry backward-compatible fixes only; breaking changes always ride the WSDL minor (e.g. `v0.56.0`).
- The version is injected via `-ldflags "-X main.version=..."`; visible via `uapi-go-daemon --version` and in startup logs.

---

## v0.55.1 (2026-08)

Upstream contract **v54_0 → v55_0 major upgrade** plus three accompanying refactors. **This version breaks compatibility with the old one**; callers must follow the migration notes below.

### ⚠️ Breaking changes (read before upgrading)

1. **Universal hotel booking: `hotelProperty` changes from a JSON array to a single object**
   - Affected endpoint: `POST /api/universal/hotel-create-reservation` (and the `/api/hotel/book` alias — same engine).
   - Reason: the upstream v55 contract narrows `HotelProperty`/`HotelStay` cardinality from `0..99` to exactly 1, and `Guarantee` etc. from `0..99` to `0..1`.
   - Decision: the gateway **ships no compatibility shim**, to avoid maintaining two shapes (rationale in `docs/architecture.md` ADR-003). Callers change `"hotelProperty": [ {...} ]` to `"hotelProperty": { ... }`.
   - Field-level migration table (`HotelCreateReservationReq`, v54 → v55):

     | Field | v54 JSON shape | v55 JSON shape | Semantic change |
     |---|---|---|---|
     | `hotelProperty` | array (0..99, optional) | **required single object** (1..1) | one hotel booking per request; multiple hotels become multiple calls |
     | `hotelStay` | array (0..99, optional) | **required single object** (1..1) | pairs 1:1 with hotelProperty |
     | `guarantee` | array (0..99) | optional single object (0..1) | at most one guarantee |
     | `guestInformation` | array (0..99) | optional single object (0..1) | at most one guest information block |
     | `reservationName` | array (0..99) | optional single object (0..1) | at most one reservation name |
     | `hotelSpecialRequest` | array of objects (`[{key, hotelRateDetailRef}]`) | **optional plain-text string** (0..1, type simplified to `typeGeneralText`) | no longer a struct/array — pass text directly |

2. **Hotel `/search` and `/details` JSON contracts rebuilt on generated types**
   - Requests: the `transactionId` input **no longer exists** (not in the XSD contract; the old hand-written model invented it; tracing is injected server-side via TraceId).
   - Responses: the old "condensed" summaries become **complete WSDL-generated models** (different field set, richer information).
   - `nextResultReference` is **visible** in requests/responses (the details paging token; the server auto-pages, so callers normally need not send it).
   - Decision record: ADR-005 (full retirement of hand-written SOAP models; generated types as the single source of truth).

3. **`sharedUprofile` domain dependency chain fixed**: the `uprofileCommon_v30_0` contract is archived and registered with the generator; the domain's cross-package references now resolve by their real namespace instead of "fall back to the newest common package". If you integrated this domain before, re-run joint testing.

4. **Generated package names de-versioned** (contributors only): `air55`→`air`, `hotel55`→`hotel`, etc.; the `commonNN` family keeps versions because several coexist. From now on contract upgrades (v55→v56) change zero import paths (ADR-004).

### Added

- **`POST /api/sharedBooking/booking-vehicle-pnr-element`**: a v55 capability — **add / update / delete vehicle elements** on a shared-booking PNR (`addVehiclePnrElement` / `updateVehiclePnrElement` / `deleteVehiclePnrElement`), with `vehicleRateChangedInfo` in the response. Maps to the new upstream PortType `BookingVehiclePnrElementPortType`.
- Public endpoint count **180 → 181**.

### Changed

- All 11 core domains migrated wholesale to the v55_0 contract (air/common/cruise/gdsQueue/hotel/passive/rail/sharedBooking/universal/util/vehicle); contract-level diff conclusions and decisions in ADR-003.
- The hotel `/search` and `/details` hand-written SOAP models (~20 `Response*` types) and the demo-era dead code `Search`/`Book` are deleted — about 500 net lines removed; `/search` and `/details` now register through the generic `registerPortHandler`, with all 8 hotel operations sharing one validation path.
- The generator manifest switched to versionless package names; cross-namespace disambiguation type prefixes de-versioned in sync (`Rail55Characteristic`→`RailCharacteristic`).
- The v54 namespace URI bump for hand-written hotel models was resolved by their retirement; repo-wide test assertions moved to v55.

### Fixed

- Fixed read-only permissions on the `wsdl/` directory preventing deliveries from being written.
- Fixed the dependency break from the missing `sharedUprofile_v20_0` delivery (see breaking change 3).

### Upgrade and verification

- Generation scale: 2411 complexTypes / 22 packages; `build` / `vet` / `gofmt` / all tests pass; daemon smoke test (`/ready`, new-route validation chain) passes.
- Contract upgrade playbook: `docs/development.md` §6; full decision log: `docs/architecture.md` §9.

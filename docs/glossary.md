# Travel Industry / GDS Glossary

**[English](./glossary.md)** | [简体中文](./glossary.zh-CN.md)

> For developers integrating Travelport UAPI for the first time. Knowing these terms makes [`architecture.md`](./architecture.md) and [`routing.md`](./routing.md) much easier to read.
> Terms follow common industry usage; parentheses note where they appear in this project's code / routes.

---

## Basics

**GDS (Global Distribution System)**
 The hub booking network for travel content: hotels, air and more. The big three: Amadeus, Sabre, Travelport (formed by merging Galileo / Apollo / Worldspan). This project integrates **Travelport**.

**UAPI / Universal API**
 Umbrella term for Travelport's programmatic interfaces. It exposes air, hotel, rail, vehicle, profile and other capabilities uniformly as SOAP/XML; this project wraps it as REST/JSON.

**SOAP / XML**
 An XML-based remote-call protocol. UAPI requests are wrapped in `<SOAP:Envelope>` and every element carries a namespace URI (e.g. `http://www.travelport.com/schema/air_v55_0`). `pkg/client` assembles and parses these.

**WSDL / XSD**
 WSDL describes "which services and operations exist"; XSD describes "what requests / responses look like". Travelport delivers them by version (e.g. `air_v55_0`); this project archives them under `wsdl/` and generates Go code at build time.

**PortType / Operation**
 A specific callable action in UAPI, e.g. `AirSearch`, `HotelCreateReservation`. Every HTTP route in this project maps to one PortType.

**Provider**
 The party actually supplying inventory — typically the airline / hotel group behind the GDS. Requests often specify `ProviderCode` (e.g. `1G` = Galileo).

**Branch / TargetBranch**
 A branch under a Travelport account (one agency office or sub-account). Many operations require `TargetBranch` to state "as which identity" the request is sent. In this project the caller provides it in the request body.

---

## Booking-related

**PNR (Passenger Name Record)**
 The core record of one booking: passengers, segments, contacts, fares. Traditionally edited via terminal commands (like `1`/`NM`/`SS`); this project's `sharedBooking` domain offers element-level REST operations.

**Universal Record (UR)**
 Travelport's mechanism for aggregating **cross-product** bookings (air + hotel + rail + vehicle + passive segments) into one record. Create / cancel / unified modify all go through `UniversalRecordService`. That is why this project's `/air/book`, `/hotel/cancel` etc. "alias routes" land in `UniversalFacade` (see [`architecture.md` §4](./architecture.md)).

**Provider Reservation**
 The provider-side actual booking of one product segment inside a UR. This project exposes `provider-reservation-display-details`, `provider-reservation-divide` and friends.

**Saved Trip**
 An itinerary draft that is "not actually ticketed / temporarily stored" and can be retrieved later to continue. Maps to the `saved-trip-*` operation group.

**Passive Segment**
 A segment not actually booked in the GDS, kept purely for accounting/notation (e.g. a hotel booked outside GDS channels). It has no standalone WSDL and can only be created / cancelled via `UniversalRecordService`'s `Passive*` operations. This project proxies it through the `passive` domain.

**CreateReservation / Cancel**
 Each product's "commit the booking" action. Note that air/hotel/rail/vehicle create/cancel **live not in the product's own service but in UniversalRecordService** (e.g. `AirCreateReservation`, `HotelCancel`), hence this project's product alias routes.

---

## Air-specific

**Air Search / Low Fare Search / Availability**
 Granularities of flight search: `AirSearch` by specific flight, `LowFareSearch` for lowest fares, `AvailabilitySearch` for seat availability.

**Price / Reprice / Fare**
 Turning search results into concrete prices. `AirPrice` prices; `AirReprice` recomputes on an existing booking; `AirFareRules` fetches fare rules. A fare is an airline-published combination of price rules.

**Ticketing / Void / Refund**
 `AirTicketing` actually issues; `AirVoidDocument` voids unsettled tickets; `AirRefund` refunds. A document means a ticket or an EMD.

**Exchange**
 Rebooking / upgrading / reissuing after ticketing, involving fare differences: `AirExchangeQuote` quotes, `AirExchange` executes, `AirExchangeTicketing` issues the exchange.

**EMD (Electronic Miscellaneous Document)**
 A voucher for ancillary services beyond the ticket (baggage fees, seat fees). `emd-issuance` / `emd-retrieve`.

**Seat Map**
 A flight's seat layout and occupancy.

**Merchandising / Upsell / Ancillary**
 Paid services beyond the flight itself (priority boarding, extra baggage). `air-merchandising-*`, `air-upsell-search`, `air-pre-pay` (prepaid).

**Schedule Change**
 After an airline adjusts schedules, `AckScheduleChange` acknowledges the change notice.

---

## Profiles

**UProfile (Universal Profile)**
 Travelport's traveler / agency profile system storing traveler preferences, payment methods, loyalty numbers and the like. The `uprofile` domain is an agency's own profiles; `sharedUprofile` are profiles shareable with Travelport. Common actions: search / retrieve / create / modify / delete, plus sub-resource management for fields / tags / hierarchy / templates.

**Hierarchy Level**
 Levels in an agency's organizational structure (company → department → branch), used for permissions and data sharing.

---

## Utilities and miscellany

**MCO (Miscellaneous Charges Order)**
 A settleable charge document, often for agency fees or change fees. The `mco-*` operation group (create / issue / exchange / void / search / retrieve).

**MCT (Minimum Connection Time)**
 The minimum time needed to connect at one airport; queried via `mct-lookup` / `mct-count`.

**Reference Data**
 Code tables of all kinds (cities, airports, countries, currencies, fare types, ...), queried / updated via `reference-data-*`.

**GDS Queue**
 Background work queues that dispatch pending PNRs to agents. The `gdsQueue` domain offers list / count / place / remove / next and friends.

**Terminal**
 Sending native GDS commands directly (like `HELP`, `*A`). The `terminal` domain provides session create / send / end.

**Branded Fare**
 Airlines packaging fares into named products (e.g. "full-flex economy", "light"); queried / managed via `branded-fare-*`.

**Agency Service Fee**
 Service fees an agency charges its customers; created via `agency-service-fee-create`.

---

## Project-wide terms

**Facade (use-case layer)**
 One struct per domain in `pkg/usecase` (e.g. `AirFacade`, `UniversalFacade`) that maps REST inputs into SOAP requests and invokes the service. The "SOAP operation" column in the routing table lists facade method names.

**Alias Route**
 A product URL (e.g. `/api/air/book`) proxying to `UniversalFacade`'s cross-product operations — product-flavored URLs reusing the UniversalRecord engine. Nine in total.

**trace_id**
 The tracking ID (UUID v4) spanning the whole request chain, visible simultaneously in the `X-Trace-Id` HTTP header, logs, and the SOAP request's common attributes — used to correlate raw XML payloads during troubleshooting.

# UAPI Go Route ↔ Operation ↔ Service Mapping

**[English](./routing.md)** | [简体中文](./routing.zh-CN.md)

> This table lists every HTTP endpoint the gateway exposes (**181** in total). All endpoints are `POST`, mounted under the `/api` prefix, with JSON request bodies; responses are strongly typed structs of the upstream SOAP responses (snake_case JSON).
> Route naming rule: `POST /api/<domain>/<operation>`.
> Routes marked **alias → UniversalRecord** are proxied by the product handler to `UniversalFacade` (see [`architecture.md` §4](./architecture.md)) and hit `UniversalRecordService` underneath.

Legend:
- **Upstream service**: the Travelport SOAP service the operation ultimately calls (`AirService` / `HotelService` / `RailService` / `VehicleService` / `UniversalRecordService` / `UProfileService` / `SharedBookingService` / `GdsQueueService` / `UtilService` / `SystemService` / `TerminalService`).
- **SOAP operation**: the corresponding facade method name — the UAPI PortType operation name (minus the `Req`/`Rsp` suffix).

---

## Air (`/api/air`, 29 endpoints)

| Route | SOAP operation | Upstream service | Notes |
|---|---|---|---|
| `/api/air/availability-search` | `AvailabilitySearch` | AirService | availability search |
| `/api/air/low-fare-search` | `LowFareSearch` | AirService | lowest-fare search |
| `/api/air/schedule-search` | `ScheduleSearch` | AirService | schedule search |
| `/api/air/flight-time-table` | `FlightTimeTable` | AirService | flight timetable |
| `/api/air/flight-details` | `FlightDetails` | AirService | flight details |
| `/api/air/flight-information` | `FlightInformation` | AirService | flight information |
| `/api/air/air-price` | `AirPrice` | AirService | pricing |
| `/api/air/air-reprice` | `AirReprice` | AirService | repricing |
| `/api/air/air-fare-display` | `AirFareDisplay` | AirService | fare display |
| `/api/air/air-fare-rules` | `AirFareRules` | AirService | fare rules |
| `/api/air/air-ticketing` | `AirTicketing` | AirService | ticketing |
| `/api/air/air-void-document` | `AirVoidDocument` | AirService | void document |
| `/api/air/air-retrieve-document` | `AirRetrieveDocument` | AirService | retrieve document |
| `/api/air/air-refund` | `AirRefund` | AirService | refund |
| `/api/air/air-refund-quote` | `AirRefundQuote` | AirService | refund quote |
| `/api/air/air-exchange` | `AirExchange` | AirService | exchange |
| `/api/air/air-exchange-quote` | `AirExchangeQuote` | AirService | exchange quote |
| `/api/air/air-exchange-eligibility` | `AirExchangeEligibility` | AirService | exchange eligibility |
| `/api/air/air-exchange-multi-quote` | `AirExchangeMultiQuote` | AirService | multi-segment exchange quote |
| `/api/air/air-exchange-ticketing` | `AirExchangeTicketing` | AirService | exchange ticketing |
| `/api/air/air-merchandising-details` | `AirMerchandisingDetails` | AirService | ancillary details |
| `/api/air/air-merchandising-offer-availability` | `AirMerchandisingOfferAvailability` | AirService | ancillary offer availability |
| `/api/air/air-upsell-search` | `AirUpsellSearch` | AirService | upsell search |
| `/api/air/air-pre-pay` | `AirPrePay` | AirService | pre-pay |
| `/api/air/emd-issuance` | `EMDIssuance` | AirService | issue EMD |
| `/api/air/emd-retrieve` | `EMDRetrieve` | AirService | retrieve EMD |
| `/api/air/seat-map` | `SeatMap` | AirService | seat map |
| `/api/air/book` | `AirCreateReservation` | **UniversalRecordService (alias)** | create air booking |
| `/api/air/cancel` | `AirCancel` | **UniversalRecordService (alias)** | cancel air booking |

---

## Hotel (`/api/hotel`, 10 endpoints)

| Route | SOAP operation | Upstream service | Notes |
|---|---|---|---|
| `/api/hotel/search` | `HotelSearchAvailability` (generated types, direct) | HotelService | hotel availability search |
| `/api/hotel/details` | `HotelDetails` (generated types, direct + auto-paging) | HotelService | hotel details / rate rules |
| `/api/hotel/retrieve` | `Retrieve` | HotelService | retrieve hotel booking |
| `/api/hotel/rules` | `Rules` | HotelService | rate rules |
| `/api/hotel/media-links` | `MediaLinks` | HotelService | media links |
| `/api/hotel/keywords` | `Keywords` | HotelService | keywords |
| `/api/hotel/upsell` | `UpsellSearch` | HotelService | upsell search |
| `/api/hotel/super-shopper` | `SuperShopper` | HotelService | super shopper search |
| `/api/hotel/book` | `HotelCreateReservation` | **UniversalRecordService (alias)** | create hotel booking |
| `/api/hotel/cancel` | `HotelCancel` | **UniversalRecordService (alias)** | cancel hotel booking |

> Note: `/api/hotel/search` and `/api/hotel/details` use generated types directly (friendlier business fields); the rest are strongly typed pass-throughs. `hotel/book` and `hotel/cancel` proxy to UniversalRecord (hotel has no standalone create/cancel portType).

---

## Rail (`/api/rail`, 7 endpoints)

| Route | SOAP operation | Upstream service | Notes |
|---|---|---|---|
| `/api/rail/rail-availability-search` | `RailAvailabilitySearch` | RailService | rail availability search |
| `/api/rail/rail-seat-map` | `RailSeatMap` | RailService | rail seat map |
| `/api/rail/rail-exchange` | `RailExchange` | RailService | rail exchange |
| `/api/rail/rail-exchange-quote` | `RailExchangeQuote` | RailService | rail exchange quote |
| `/api/rail/rail-refund` | `RailRefund` | RailService | rail refund |
| `/api/rail/rail-refund-quote` | `RailRefundQuote` | RailService | rail refund quote |
| `/api/rail/book` | `RailCreateReservation` | **UniversalRecordService (alias)** | create rail booking (no `/rail/cancel`; cancellation goes through `/api/universal/universal-record-cancel`) |

---

## Vehicle (`/api/vehicle`, 10 endpoints)

| Route | SOAP operation | Upstream service | Notes |
|---|---|---|---|
| `/api/vehicle/vehicle-search-availability` | `VehicleSearchAvailability` | VehicleService | vehicle availability search |
| `/api/vehicle/vehicle-upsell-search-availability` | `VehicleUpsellSearchAvailability` | VehicleService | vehicle upsell search |
| `/api/vehicle/vehicle-location` | `VehicleLocation` | VehicleService | vehicle locations |
| `/api/vehicle/vehicle-location-detail` | `VehicleLocationDetail` | VehicleService | location details |
| `/api/vehicle/vehicle-rules` | `VehicleRules` | VehicleService | vehicle rules |
| `/api/vehicle/vehicle-retrieve` | `VehicleRetrieve` | VehicleService | retrieve vehicle booking |
| `/api/vehicle/vehicle-keyword` | `VehicleKeyword` | VehicleService | keywords |
| `/api/vehicle/vehicle-media-links` | `VehicleMediaLinks` | VehicleService | media links |
| `/api/vehicle/book` | `VehicleCreateReservation` | **UniversalRecordService (alias)** | create vehicle booking |
| `/api/vehicle/cancel` | `VehicleCancel` | **UniversalRecordService (alias)** | cancel vehicle booking |

---

## Universal Record (`/api/universal`, 23 endpoints)

Cross-product create / cancel / unified modify / saved trips / passive segments all live here.

| Route | SOAP operation | Upstream service | Notes |
|---|---|---|---|
| `/api/universal/air-create-reservation` | `AirCreateReservation` | UniversalRecordService | create air booking |
| `/api/universal/air-cancel` | `AirCancel` | UniversalRecordService | cancel air |
| `/api/universal/air-merchandising-fulfillment` | `AirMerchandisingFulfillment` | UniversalRecordService | ancillary fulfillment |
| `/api/universal/hotel-create-reservation` | `HotelCreateReservation` | UniversalRecordService | create hotel booking |
| `/api/universal/hotel-cancel` | `HotelCancel` | UniversalRecordService | cancel hotel |
| `/api/universal/rail-create-reservation` | `RailCreateReservation` | UniversalRecordService | create rail booking |
| `/api/universal/vehicle-create-reservation` | `VehicleCreateReservation` | UniversalRecordService | create vehicle booking |
| `/api/universal/vehicle-cancel` | `VehicleCancel` | UniversalRecordService | cancel vehicle |
| `/api/universal/passive-create-reservation` | `PassiveCreateReservation` | UniversalRecordService | create passive segment |
| `/api/universal/passive-cancel` | `PassiveCancel` | UniversalRecordService | cancel passive segment |
| `/api/universal/universal-record-cancel` | `UniversalRecordCancel` | UniversalRecordService | cancel universal record (rail cancellation goes here) |
| `/api/universal/universal-record-import` | `UniversalRecordImport` | UniversalRecordService | import universal record |
| `/api/universal/universal-record-modify` | `UniversalRecordModify` | UniversalRecordService | modify universal record |
| `/api/universal/universal-record-retrieve` | `UniversalRecordRetrieve` | UniversalRecordService | retrieve universal record |
| `/api/universal/universal-record-search` | `UniversalRecordSearch` | UniversalRecordService | search universal records |
| `/api/universal/provider-reservation-display-details` | `ProviderReservationDisplayDetails` | UniversalRecordService | provider reservation details |
| `/api/universal/provider-reservation-divide` | `ProviderReservationDivide` | UniversalRecordService | divide provider reservation |
| `/api/universal/ack-schedule-change` | `AckScheduleChange` | UniversalRecordService | acknowledge schedule change |
| `/api/universal/saved-trip-create` | `SavedTripCreate` | UniversalRecordService | create saved trip |
| `/api/universal/saved-trip-delete` | `SavedTripDelete` | UniversalRecordService | delete saved trip |
| `/api/universal/saved-trip-modify` | `SavedTripModify` | UniversalRecordService | modify saved trip |
| `/api/universal/saved-trip-retrieve` | `SavedTripRetrieve` | UniversalRecordService | retrieve saved trip |
| `/api/universal/saved-trip-search` | `SavedTripSearch` | UniversalRecordService | search saved trips |

---

## Passive (`/api/passive`, 2 endpoints)

Passive has no standalone WSDL / portType in UAPI and is only reachable via UniversalRecord operations; a dedicated handler proxies it.

| Route | SOAP operation | Upstream service | Notes |
|---|---|---|---|
| `/api/passive/book` | `PassiveCreateReservation` | **UniversalRecordService (alias)** | create passive segment |
| `/api/passive/cancel` | `PassiveCancel` | **UniversalRecordService (alias)** | cancel passive segment |

---

## Shared Booking (`/api/sharedBooking`, 15 endpoints)

PNR element-level booking operations (Travelport Shared Booking).

| Route | SOAP operation | Upstream service | Notes |
|---|---|---|---|
| `/api/sharedBooking/booking-start` | `BookingStart` | SharedBookingService | start booking |
| `/api/sharedBooking/booking-end` | `BookingEnd` | SharedBookingService | end booking |
| `/api/sharedBooking/booking-display` | `BookingDisplay` | SharedBookingService | display booking |
| `/api/sharedBooking/booking-traveler` | `BookingTraveler` | SharedBookingService | traveler |
| `/api/sharedBooking/booking-vehicle-pnr-element` | `BookingVehiclePnrElement` | SharedBookingService | vehicle PNR element (v55 addition: add/update/delete vehicle elements) |
| `/api/sharedBooking/booking-pnr-element` | `BookingPnrElement` | SharedBookingService | PNR element |
| `/api/sharedBooking/booking-air-segment` | `BookingAirSegment` | SharedBookingService | air segment |
| `/api/sharedBooking/booking-air-pnr-element` | `BookingAirPnrElement` | SharedBookingService | air PNR element |
| `/api/sharedBooking/booking-air-exchange` | `BookingAirExchange` | SharedBookingService | air exchange |
| `/api/sharedBooking/booking-air-exchange-quote` | `BookingAirExchangeQuote` | SharedBookingService | air exchange quote |
| `/api/sharedBooking/booking-hotel-segment` | `BookingHotelSegment` | SharedBookingService | hotel segment |
| `/api/sharedBooking/booking-hotel-pnr-element` | `BookingHotelPnrElement` | SharedBookingService | hotel PNR element |
| `/api/sharedBooking/booking-pricing` | `BookingPricing` | SharedBookingService | pricing |
| `/api/sharedBooking/booking-seat-assignment` | `BookingSeatAssignment` | SharedBookingService | seat assignment |
| `/api/sharedBooking/booking-retrieve-document` | `BookingRetrieveDocument` | SharedBookingService | retrieve document |
| `/api/sharedBooking/booking-terminal` | `BookingTerminal` | SharedBookingService | terminal command |

---

## UProfile (`/api/uprofile`, 25 endpoints)

| Route | SOAP operation | Upstream service | Notes |
|---|---|---|---|
| `/api/uprofile/profile-search` | `ProfileSearch` | UProfileService | profile search |
| `/api/uprofile/profile-retrieve` | `ProfileRetrieve` | UProfileService | retrieve profile |
| `/api/uprofile/profile-create` | `ProfileCreate` | UProfileService | create profile |
| `/api/uprofile/profile-modify` | `ProfileModify` | UProfileService | modify profile |
| `/api/uprofile/profile-delete` | `ProfileDelete` | UProfileService | delete profile |
| `/api/uprofile/profile-retrieve-history` | `ProfileRetrieveHistory` | UProfileService | profile history |
| `/api/uprofile/profile-child-search` | `ProfileChildSearch` | UProfileService | child profile search |
| `/api/uprofile/profile-create-field` | `ProfileCreateField` | UProfileService | create field |
| `/api/uprofile/profile-modify-field` | `ProfileModifyField` | UProfileService | modify field |
| `/api/uprofile/profile-search-field` | `ProfileSearchField` | UProfileService | search field |
| `/api/uprofile/profile-create-tags` | `ProfileCreateTags` | UProfileService | create tags |
| `/api/uprofile/profile-modify-tags` | `ProfileModifyTags` | UProfileService | modify tags |
| `/api/uprofile/profile-delete-tag` | `ProfileDeleteTag` | UProfileService | delete tag |
| `/api/uprofile/profile-search-tags` | `ProfileSearchTags` | UProfileService | search tags |
| `/api/uprofile/profile-create-hierarchy-level` | `ProfileCreateHierarchyLevel` | UProfileService | create hierarchy level |
| `/api/uprofile/profile-delete-hierarchy-level` | `ProfileDeleteHierarchyLevel` | UProfileService | delete hierarchy level |
| `/api/uprofile/profile-modify-hierarchy-level` | `ProfileModifyHierarchyLevel` | UProfileService | modify hierarchy level |
| `/api/uprofile/profile-retrieve-hierarchy` | `ProfileRetrieveHierarchy` | UProfileService | retrieve hierarchy |
| `/api/uprofile/profile-modify-bridge-branches` | `ProfileModifyBridgeBranches` | UProfileService | modify bridge branches |
| `/api/uprofile/profile-retrieve-bridge-branches` | `ProfileRetrieveBridgeBranches` | UProfileService | retrieve bridge branches |
| `/api/uprofile/profile-modify-template` | `ProfileModifyTemplate` | UProfileService | modify template |
| `/api/uprofile/profile-retrieve-template` | `ProfileRetrieveTemplate` | UProfileService | retrieve template |
| `/api/uprofile/profile-retrieve-action` | `ProfileRetrieveAction` | UProfileService | retrieve action |
| `/api/uprofile/profile-search-action` | `ProfileSearchAction` | UProfileService | search action |
| `/api/uprofile/single-profile-migration` | `SingleProfileMigration` | UProfileService | single profile migration |

---

## Shared UProfile (`/api/sharedUprofile`, 20 endpoints)

| Route | SOAP operation | Upstream service | Notes |
|---|---|---|---|
| `/api/sharedUprofile/profile-search` | `ProfileSearch` | SharedUProfileService | profile search |
| `/api/sharedUprofile/profile-retrieve` | `ProfileRetrieve` | SharedUProfileService | retrieve profile |
| `/api/sharedUprofile/profile-create` | `ProfileCreate` | SharedUProfileService | create profile |
| `/api/sharedUprofile/profile-modify` | `ProfileModify` | SharedUProfileService | modify profile |
| `/api/sharedUprofile/profile-delete` | `ProfileDelete` | SharedUProfileService | delete profile |
| `/api/sharedUprofile/profile-retrieve-history` | `ProfileRetrieveHistory` | SharedUProfileService | profile history |
| `/api/sharedUprofile/profile-retrieve-parent` | `ProfileRetrieveParent` | SharedUProfileService | retrieve parent profile |
| `/api/sharedUprofile/profile-child-search` | `ProfileChildSearch` | SharedUProfileService | child profile search |
| `/api/sharedUprofile/profile-create-field` | `ProfileCreateField` | SharedUProfileService | create field |
| `/api/sharedUprofile/profile-modify-field` | `ProfileModifyField` | SharedUProfileService | modify field |
| `/api/sharedUprofile/profile-search-field` | `ProfileSearchField` | SharedUProfileService | search field |
| `/api/sharedUprofile/profile-create-tags` | `ProfileCreateTags` | SharedUProfileService | create tags |
| `/api/sharedUprofile/profile-modify-tags` | `ProfileModifyTags` | SharedUProfileService | modify tags |
| `/api/sharedUprofile/profile-delete-tag` | `ProfileDeleteTag` | SharedUProfileService | delete tag |
| `/api/sharedUprofile/profile-search-tags` | `ProfileSearchTags` | SharedUProfileService | search tags |
| `/api/sharedUprofile/single-profile-migration` | `SingleProfileMigration` | SharedUProfileService | single profile migration |
| `/api/sharedUprofile/ui-meta-data-create` | `UIMetaDataCreate` | SharedUProfileService | create UI metadata |
| `/api/sharedUprofile/ui-meta-data-delete` | `UIMetaDataDelete` | SharedUProfileService | delete UI metadata |
| `/api/sharedUprofile/ui-meta-data-modify` | `UIMetaDataModify` | SharedUProfileService | modify UI metadata |
| `/api/sharedUprofile/ui-meta-data-retrieve` | `UIMetaDataRetrieve` | SharedUProfileService | retrieve UI metadata |

---

## GdsQueue (`/api/gdsQueue`, 8 endpoints)

| Route | SOAP operation | Upstream service | Notes |
|---|---|---|---|
| `/api/gdsQueue/gds-queue-list` | `GdsQueueList` | GdsQueueService | queue list |
| `/api/gdsQueue/gds-queue-count` | `GdsQueueCount` | GdsQueueService | queue count |
| `/api/gdsQueue/gds-queue-place` | `GdsQueuePlace` | GdsQueueService | place on queue |
| `/api/gdsQueue/gds-queue-remove` | `GdsQueueRemove` | GdsQueueService | remove from queue |
| `/api/gdsQueue/gds-enter-queue` | `GdsEnterQueue` | GdsQueueService | enter queue |
| `/api/gdsQueue/gds-exit-queue` | `GdsExitQueue` | GdsQueueService | exit queue |
| `/api/gdsQueue/gds-next-on-queue` | `GdsNextOnQueue` | GdsQueueService | next on queue |
| `/api/gdsQueue/gds-queue-agent-list` | `GdsQueueAgentList` | GdsQueueService | queue agent list |

---

## Terminal (`/api/terminal`, 3 endpoints)

| Route | SOAP operation | Upstream service | Notes |
|---|---|---|---|
| `/api/terminal/create-terminal-session` | `CreateTerminalSession` | TerminalService | create terminal session |
| `/api/terminal/terminal` | `Terminal` | TerminalService | send terminal command |
| `/api/terminal/end-terminal-session` | `EndTerminalSession` | TerminalService | end terminal session |

---

## System (`/api/system`, 4 endpoints)

| Route | SOAP operation | Upstream service | Notes |
|---|---|---|---|
| `/api/system/ping` | `Ping` | SystemService | health ping |
| `/api/system/info` | `Info` | SystemService | service info |
| `/api/system/time` | `Time` | SystemService | service time |
| `/api/system/cache` | `ExternalCacheAccess` | SystemService | external cache access |

---

## Util (`/api/util`, 24 endpoints)

| Route | SOAP operation | Upstream service | Notes |
|---|---|---|---|
| `/api/util/calculate-tax` | `CalculateTax` | UtilService | calculate tax |
| `/api/util/currency-conversion` | `CurrencyConversion` | UtilService | currency conversion |
| `/api/util/credit-card-auth` | `CreditCardAuth` | UtilService | credit card authorization |
| `/api/util/mct-count` | `MctCount` | UtilService | minimum connecting time count |
| `/api/util/mct-lookup` | `MctLookup` | UtilService | minimum connecting time lookup |
| `/api/util/mco-create` | `MCOCreate` | UtilService | create MCO |
| `/api/util/mco-exchange` | `MCOExchange` | UtilService | exchange MCO |
| `/api/util/mco-issue` | `MCOIssue` | UtilService | issue MCO |
| `/api/util/mco-retrieve` | `MCORetrieve` | UtilService | retrieve MCO |
| `/api/util/mco-search` | `McoSearch` | UtilService | search MCO |
| `/api/util/mco-void` | `McoVoid` | UtilService | void MCO |
| `/api/util/create-agency-fee-mco` | `CreateAgencyFeeMco` | UtilService | create agency-fee MCO |
| `/api/util/create-airline-fee-mco` | `CreateAirlineFeeMco` | UtilService | create airline-fee MCO |
| `/api/util/agency-service-fee-create` | `AgencyServiceFeeCreate` | UtilService | create agency service fee |
| `/api/util/branded-fare-search` | `BrandedFareSearch` | UtilService | branded fare search |
| `/api/util/branded-fare-admin` | `BrandedFareAdmin` | UtilService | branded fare admin |
| `/api/util/upsell-admin` | `UpsellAdmin` | UtilService | upsell admin |
| `/api/util/upsell-search` | `UpsellSearch` | UtilService | upsell search |
| `/api/util/reference-data-retrieve` | `ReferenceDataRetrieve` | UtilService | retrieve reference data |
| `/api/util/reference-data-search` | `ReferenceDataSearch` | UtilService | search reference data |
| `/api/util/reference-data-update` | `ReferenceDataUpdate` | UtilService | update reference data |
| `/api/util/content-provider-retrieve` | `ContentProviderRetrieve` | UtilService | retrieve content provider |
| `/api/util/mir-report-retrieve` | `MirReportRetrieve` | UtilService | retrieve MIR report |
| `/api/util/find-employees-on-flight` | `FindEmployeesOnFlight` | UtilService | find employees on flight |

---

## Statistics

| Domain | Routes | of which aliases (→UniversalRecord) |
|---|---|---|
| air | 29 | 2 |
| hotel | 10 (incl. 2 custom DTOs: search/details) | 2 |
| rail | 7 | 1 |
| vehicle | 10 | 2 |
| universal | 23 | — |
| passive | 2 | 2 |
| sharedBooking | 15 | — |
| uprofile | 25 | — |
| sharedUprofile | 20 | — |
| gdsQueue | 8 | — |
| system | 4 | — |
| terminal | 3 | — |
| util | 24 | — |
| **Total** | **181** | **9** |

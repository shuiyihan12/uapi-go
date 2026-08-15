# UAPI Go 路由 ↔ 操作 ↔ 服务 对照表

**[English](./routing.md)** | 简体中文

> 本表列出网关对外暴露的全部 HTTP 接口（共 **181** 个）。所有接口均为 `POST`，统一挂在 `/api` 前缀下，请求体为 JSON，响应为上游 SOAP 响应的强类型结构体（snake_case JSON）。
> 路由命名规则：`POST /api/<域>/<操作>`。
> 标记为 **别名 → UniversalRecord** 的路由，由产品 Handler 代理到 `UniversalFacade`（详见 [`architecture.md` §4](./architecture.zh-CN.md)），底层走 `UniversalRecordService`。

图例：
- **上游服务**：该操作最终调用的 Travelport SOAP Service（`AirService` / `HotelService` / `RailService` / `VehicleService` / `UniversalRecordService` / `UProfileService` / `SharedBookingService` / `GdsQueueService` / `UtilService` / `SystemService` / `TerminalService`）。
- **SOAP 操作**：对应的 Facade 方法名，即 UAPI PortType 操作名（去掉 `Req`/`Rsp`）。

---

## 机票 Air（`/api/air`，29 个）

| 路由 | SOAP 操作 | 上游服务 | 说明 |
|---|---|---|---|
| `/api/air/availability-search` | `AvailabilitySearch` | AirService |  availability 检索 |
| `/api/air/low-fare-search` | `LowFareSearch` | AirService | 最低价检索 |
| `/api/air/schedule-search` | `ScheduleSearch` | AirService | 班期检索 |
| `/api/air/flight-time-table` | `FlightTimeTable` | AirService | 航班时刻表 |
| `/api/air/flight-details` | `FlightDetails` | AirService | 航班详情 |
| `/api/air/flight-information` | `FlightInformation` | AirService | 航班信息 |
| `/api/air/air-price` | `AirPrice` | AirService | 计价 |
| `/api/air/air-reprice` | `AirReprice` | AirService | 重新计价 |
| `/api/air/air-fare-display` | `AirFareDisplay` | AirService | 票价展示 |
| `/api/air/air-fare-rules` | `AirFareRules` | AirService | 票价规则 |
| `/api/air/air-ticketing` | `AirTicketing` | AirService | 出票 |
| `/api/air/air-void-document` | `AirVoidDocument` | AirService | 作废票证 |
| `/api/air/air-retrieve-document` | `AirRetrieveDocument` | AirService | 取回票证 |
| `/api/air/air-refund` | `AirRefund` | AirService | 退票 |
| `/api/air/air-refund-quote` | `AirRefundQuote` | AirService | 退票报价 |
| `/api/air/air-exchange` | `AirExchange` | AirService | 换开 |
| `/api/air/air-exchange-quote` | `AirExchangeQuote` | AirService | 换开报价 |
| `/api/air/air-exchange-eligibility` | `AirExchangeEligibility` | AirService | 换开资格 |
| `/api/air/air-exchange-multi-quote` | `AirExchangeMultiQuote` | AirService | 多段换开报价 |
| `/api/air/air-exchange-ticketing` | `AirExchangeTicketing` | AirService | 换开出票 |
| `/api/air/air-merchandising-details` | `AirMerchandisingDetails` | AirService | 辅营详情 |
| `/api/air/air-merchandising-offer-availability` | `AirMerchandisingOfferAvailability` | AirService | 辅营报价 |
| `/api/air/air-upsell-search` | `AirUpsellSearch` | AirService | 升舱/加价搜索 |
| `/api/air/air-pre-pay` | `AirPrePay` | AirService | 预付 |
| `/api/air/emd-issuance` | `EMDIssuance` | AirService | 开具 EMD |
| `/api/air/emd-retrieve` | `EMDRetrieve` | AirService | 取回 EMD |
| `/api/air/seat-map` | `SeatMap` | AirService | 座位图 |
| `/api/air/book` | `AirCreateReservation` | **UniversalRecordService（别名）** | 创建机票预订 |
| `/api/air/cancel` | `AirCancel` | **UniversalRecordService（别名）** | 取消机票预订 |

---

## 酒店 Hotel（`/api/hotel`，10 个）

| 路由 | SOAP 操作 | 上游服务 | 说明 |
|---|---|---|---|
| `/api/hotel/search` | `HotelSearchAvailability`（生成类型直通） | HotelService | 酒店可用性搜索 |
| `/api/hotel/details` | `HotelDetails`（生成类型直通 + 自动翻页） | HotelService | 酒店详情 / 房价规则 |
| `/api/hotel/retrieve` | `Retrieve` | HotelService | 取回酒店预订 |
| `/api/hotel/rules` | `Rules` | HotelService | 房价规则 |
| `/api/hotel/media-links` | `MediaLinks` | HotelService | 媒体链接 |
| `/api/hotel/keywords` | `Keywords` | HotelService | 关键词 |
| `/api/hotel/upsell` | `UpsellSearch` | HotelService | 加价搜索 |
| `/api/hotel/super-shopper` | `SuperShopper` | HotelService | 超值搜索 |
| `/api/hotel/book` | `HotelCreateReservation` | **UniversalRecordService（别名）** | 创建酒店预订 |
| `/api/hotel/cancel` | `HotelCancel` | **UniversalRecordService（别名）** | 取消酒店预订 |

> 注：`/api/hotel/search` 与 `/api/hotel/details` 使用显式 DTO 映射（更友好的业务字段），其余为强类型透传。`hotel/book`、`hotel/cancel` 代理到 UniversalRecord（酒店无独立 create/cancel portType）。

---

## 火车 Rail（`/api/rail`，7 个）

| 路由 | SOAP 操作 | 上游服务 | 说明 |
|---|---|---|---|
| `/api/rail/rail-availability-search` | `RailAvailabilitySearch` | RailService | 火车余票检索 |
| `/api/rail/rail-seat-map` | `RailSeatMap` | RailService | 火车座位图 |
| `/api/rail/rail-exchange` | `RailExchange` | RailService | 火车换开 |
| `/api/rail/rail-exchange-quote` | `RailExchangeQuote` | RailService | 火车换开报价 |
| `/api/rail/rail-refund` | `RailRefund` | RailService | 火车退票 |
| `/api/rail/rail-refund-quote` | `RailRefundQuote` | RailService | 火车退票报价 |
| `/api/rail/book` | `RailCreateReservation` | **UniversalRecordService（别名）** | 创建火车预订（无 `/rail/cancel`，取消见 `/api/universal/universal-record-cancel`） |

---

## 租车 Vehicle（`/api/vehicle`，10 个）

| 路由 | SOAP 操作 | 上游服务 | 说明 |
|---|---|---|---|
| `/api/vehicle/vehicle-search-availability` | `VehicleSearchAvailability` | VehicleService | 租车可用性搜索 |
| `/api/vehicle/vehicle-upsell-search-availability` | `VehicleUpsellSearchAvailability` | VehicleService | 租车加价搜索 |
| `/api/vehicle/vehicle-location` | `VehicleLocation` | VehicleService | 租车门店 |
| `/api/vehicle/vehicle-location-detail` | `VehicleLocationDetail` | VehicleService | 门店详情 |
| `/api/vehicle/vehicle-rules` | `VehicleRules` | VehicleService | 租车规则 |
| `/api/vehicle/vehicle-retrieve` | `VehicleRetrieve` | VehicleService | 取回租车预订 |
| `/api/vehicle/vehicle-keyword` | `VehicleKeyword` | VehicleService | 关键词 |
| `/api/vehicle/vehicle-media-links` | `VehicleMediaLinks` | VehicleService | 媒体链接 |
| `/api/vehicle/book` | `VehicleCreateReservation` | **UniversalRecordService（别名）** | 创建租车预订 |
| `/api/vehicle/cancel` | `VehicleCancel` | **UniversalRecordService（别名）** | 取消租车预订 |

---

## 统一记录 Universal Record（`/api/universal`，23 个）

跨产品的创建 / 取消 / 统一修改 / 保存行程 / 被动段，全部在此。

| 路由 | SOAP 操作 | 上游服务 | 说明 |
|---|---|---|---|
| `/api/universal/air-create-reservation` | `AirCreateReservation` | UniversalRecordService | 机票创建预订 |
| `/api/universal/air-cancel` | `AirCancel` | UniversalRecordService | 机票取消 |
| `/api/universal/air-merchandising-fulfillment` | `AirMerchandisingFulfillment` | UniversalRecordService | 辅营履约 |
| `/api/universal/hotel-create-reservation` | `HotelCreateReservation` | UniversalRecordService | 酒店创建预订 |
| `/api/universal/hotel-cancel` | `HotelCancel` | UniversalRecordService | 酒店取消 |
| `/api/universal/rail-create-reservation` | `RailCreateReservation` | UniversalRecordService | 火车创建预订 |
| `/api/universal/vehicle-create-reservation` | `VehicleCreateReservation` | UniversalRecordService | 租车创建预订 |
| `/api/universal/vehicle-cancel` | `VehicleCancel` | UniversalRecordService | 租车取消 |
| `/api/universal/passive-create-reservation` | `PassiveCreateReservation` | UniversalRecordService | 被动段创建 |
| `/api/universal/passive-cancel` | `PassiveCancel` | UniversalRecordService | 被动段取消 |
| `/api/universal/universal-record-cancel` | `UniversalRecordCancel` | UniversalRecordService | 统一记录取消（火车取消走这里） |
| `/api/universal/universal-record-import` | `UniversalRecordImport` | UniversalRecordService | 导入统一记录 |
| `/api/universal/universal-record-modify` | `UniversalRecordModify` | UniversalRecordService | 统一记录修改 |
| `/api/universal/universal-record-retrieve` | `UniversalRecordRetrieve` | UniversalRecordService | 取回统一记录 |
| `/api/universal/universal-record-search` | `UniversalRecordSearch` | UniversalRecordService | 搜索统一记录 |
| `/api/universal/provider-reservation-display-details` | `ProviderReservationDisplayDetails` | UniversalRecordService | 供应商预订明细 |
| `/api/universal/provider-reservation-divide` | `ProviderReservationDivide` | UniversalRecordService | 拆分供应商预订 |
| `/api/universal/ack-schedule-change` | `AckScheduleChange` | UniversalRecordService | 确认航班变动 |
| `/api/universal/saved-trip-create` | `SavedTripCreate` | UniversalRecordService | 保存行程创建 |
| `/api/universal/saved-trip-delete` | `SavedTripDelete` | UniversalRecordService | 保存行程删除 |
| `/api/universal/saved-trip-modify` | `SavedTripModify` | UniversalRecordService | 保存行程修改 |
| `/api/universal/saved-trip-retrieve` | `SavedTripRetrieve` | UniversalRecordService | 保存行程取回 |
| `/api/universal/saved-trip-search` | `SavedTripSearch` | UniversalRecordService | 保存行程搜索 |

---

## 被动段 Passive（`/api/passive`，2 个）

Passive 在 UAPI 中无独立 WSDL / portType，只能经 UniversalRecord 操作，故单独建 Handler 代理。

| 路由 | SOAP 操作 | 上游服务 | 说明 |
|---|---|---|---|
| `/api/passive/book` | `PassiveCreateReservation` | **UniversalRecordService（别名）** | 被动段创建 |
| `/api/passive/cancel` | `PassiveCancel` | **UniversalRecordService（别名）** | 被动段取消 |

---

## 共享预订 Shared Booking（`/api/sharedBooking`，15 个）

PNR 元素级的预订操作（Travelport Shared Booking）。

| 路由 | SOAP 操作 | 上游服务 | 说明 |
|---|---|---|---|
| `/api/sharedBooking/booking-start` | `BookingStart` | SharedBookingService | 开始预订 |
| `/api/sharedBooking/booking-end` | `BookingEnd` | SharedBookingService | 结束预订 |
| `/api/sharedBooking/booking-display` | `BookingDisplay` | SharedBookingService | 显示预订 |
| `/api/sharedBooking/booking-traveler` | `BookingTraveler` | SharedBookingService | 旅客 |
| `/api/sharedBooking/booking-vehicle-pnr-element` | `BookingVehiclePnrElement` | SharedBookingService | 车辆 PNR 元素（v55 新增：新增/更新/删除车辆元素） |
| `/api/sharedBooking/booking-pnr-element` | `BookingPnrElement` | SharedBookingService | PNR 元素 |
| `/api/sharedBooking/booking-air-segment` | `BookingAirSegment` | SharedBookingService | 机票航段 |
| `/api/sharedBooking/booking-air-pnr-element` | `BookingAirPnrElement` | SharedBookingService | 机票 PNR 元素 |
| `/api/sharedBooking/booking-air-exchange` | `BookingAirExchange` | SharedBookingService | 机票换开 |
| `/api/sharedBooking/booking-air-exchange-quote` | `BookingAirExchangeQuote` | SharedBookingService | 机票换开报价 |
| `/api/sharedBooking/booking-hotel-segment` | `BookingHotelSegment` | SharedBookingService | 酒店航段 |
| `/api/sharedBooking/booking-hotel-pnr-element` | `BookingHotelPnrElement` | SharedBookingService | 酒店 PNR 元素 |
| `/api/sharedBooking/booking-pricing` | `BookingPricing` | SharedBookingService | 计价 |
| `/api/sharedBooking/booking-seat-assignment` | `BookingSeatAssignment` | SharedBookingService | 座位分配 |
| `/api/sharedBooking/booking-retrieve-document` | `BookingRetrieveDocument` | SharedBookingService | 取回票证 |
| `/api/sharedBooking/booking-terminal` | `BookingTerminal` | SharedBookingService | 终端指令 |

---

## 客户档案 UProfile（`/api/uprofile`，25 个）

| 路由 | SOAP 操作 | 上游服务 | 说明 |
|---|---|---|---|
| `/api/uprofile/profile-search` | `ProfileSearch` | UProfileService | 档案搜索 |
| `/api/uprofile/profile-retrieve` | `ProfileRetrieve` | UProfileService | 取回档案 |
| `/api/uprofile/profile-create` | `ProfileCreate` | UProfileService | 创建档案 |
| `/api/uprofile/profile-modify` | `ProfileModify` | UProfileService | 修改档案 |
| `/api/uprofile/profile-delete` | `ProfileDelete` | UProfileService | 删除档案 |
| `/api/uprofile/profile-retrieve-history` | `ProfileRetrieveHistory` | UProfileService | 档案历史 |
| `/api/uprofile/profile-child-search` | `ProfileChildSearch` | UProfileService | 子档案搜索 |
| `/api/uprofile/profile-create-field` | `ProfileCreateField` | UProfileService | 创建字段 |
| `/api/uprofile/profile-modify-field` | `ProfileModifyField` | UProfileService | 修改字段 |
| `/api/uprofile/profile-search-field` | `ProfileSearchField` | UProfileService | 搜索字段 |
| `/api/uprofile/profile-create-tags` | `ProfileCreateTags` | UProfileService | 创建标签 |
| `/api/uprofile/profile-modify-tags` | `ProfileModifyTags` | UProfileService | 修改标签 |
| `/api/uprofile/profile-delete-tag` | `ProfileDeleteTag` | UProfileService | 删除标签 |
| `/api/uprofile/profile-search-tags` | `ProfileSearchTags` | UProfileService | 搜索标签 |
| `/api/uprofile/profile-create-hierarchy-level` | `ProfileCreateHierarchyLevel` | UProfileService | 创建层级 |
| `/api/uprofile/profile-delete-hierarchy-level` | `ProfileDeleteHierarchyLevel` | UProfileService | 删除层级 |
| `/api/uprofile/profile-modify-hierarchy-level` | `ProfileModifyHierarchyLevel` | UProfileService | 修改层级 |
| `/api/uprofile/profile-retrieve-hierarchy` | `ProfileRetrieveHierarchy` | UProfileService | 取回层级 |
| `/api/uprofile/profile-modify-bridge-branches` | `ProfileModifyBridgeBranches` | UProfileService | 修改桥接分支 |
| `/api/uprofile/profile-retrieve-bridge-branches` | `ProfileRetrieveBridgeBranches` | UProfileService | 取回桥接分支 |
| `/api/uprofile/profile-modify-template` | `ProfileModifyTemplate` | UProfileService | 修改模板 |
| `/api/uprofile/profile-retrieve-template` | `ProfileRetrieveTemplate` | UProfileService | 取回模板 |
| `/api/uprofile/profile-retrieve-action` | `ProfileRetrieveAction` | UProfileService | 取回动作 |
| `/api/uprofile/profile-search-action` | `ProfileSearchAction` | UProfileService | 搜索动作 |
| `/api/uprofile/single-profile-migration` | `SingleProfileMigration` | UProfileService | 单档案迁移 |

---

## 共享客户档案 Shared UProfile（`/api/sharedUprofile`，20 个）

| 路由 | SOAP 操作 | 上游服务 | 说明 |
|---|---|---|---|
| `/api/sharedUprofile/profile-search` | `ProfileSearch` | SharedUProfileService | 档案搜索 |
| `/api/sharedUprofile/profile-retrieve` | `ProfileRetrieve` | SharedUProfileService | 取回档案 |
| `/api/sharedUprofile/profile-create` | `ProfileCreate` | SharedUProfileService | 创建档案 |
| `/api/sharedUprofile/profile-modify` | `ProfileModify` | SharedUProfileService | 修改档案 |
| `/api/sharedUprofile/profile-delete` | `ProfileDelete` | SharedUProfileService | 删除档案 |
| `/api/sharedUprofile/profile-retrieve-history` | `ProfileRetrieveHistory` | SharedUProfileService | 档案历史 |
| `/api/sharedUprofile/profile-retrieve-parent` | `ProfileRetrieveParent` | SharedUProfileService | 取回父档案 |
| `/api/sharedUprofile/profile-child-search` | `ProfileChildSearch` | SharedUProfileService | 子档案搜索 |
| `/api/sharedUprofile/profile-create-field` | `ProfileCreateField` | SharedUProfileService | 创建字段 |
| `/api/sharedUprofile/profile-modify-field` | `ProfileModifyField` | SharedUProfileService | 修改字段 |
| `/api/sharedUprofile/profile-search-field` | `ProfileSearchField` | SharedUProfileService | 搜索字段 |
| `/api/sharedUprofile/profile-create-tags` | `ProfileCreateTags` | SharedUProfileService | 创建标签 |
| `/api/sharedUprofile/profile-modify-tags` | `ProfileModifyTags` | SharedUProfileService | 修改标签 |
| `/api/sharedUprofile/profile-delete-tag` | `ProfileDeleteTag` | SharedUProfileService | 删除标签 |
| `/api/sharedUprofile/profile-search-tags` | `ProfileSearchTags` | SharedUProfileService | 搜索标签 |
| `/api/sharedUprofile/single-profile-migration` | `SingleProfileMigration` | SharedUProfileService | 单档案迁移 |
| `/api/sharedUprofile/ui-meta-data-create` | `UIMetaDataCreate` | SharedUProfileService | 创建 UI 元数据 |
| `/api/sharedUprofile/ui-meta-data-delete` | `UIMetaDataDelete` | SharedUProfileService | 删除 UI 元数据 |
| `/api/sharedUprofile/ui-meta-data-modify` | `UIMetaDataModify` | SharedUProfileService | 修改 UI 元数据 |
| `/api/sharedUprofile/ui-meta-data-retrieve` | `UIMetaDataRetrieve` | SharedUProfileService | 取回 UI 元数据 |

---

## GDS 队列 GdsQueue（`/api/gdsQueue`，8 个）

| 路由 | SOAP 操作 | 上游服务 | 说明 |
|---|---|---|---|
| `/api/gdsQueue/gds-queue-list` | `GdsQueueList` | GdsQueueService | 队列列表 |
| `/api/gdsQueue/gds-queue-count` | `GdsQueueCount` | GdsQueueService | 队列计数 |
| `/api/gdsQueue/gds-queue-place` | `GdsQueuePlace` | GdsQueueService | 入队 |
| `/api/gdsQueue/gds-queue-remove` | `GdsQueueRemove` | GdsQueueService | 出队移除 |
| `/api/gdsQueue/gds-enter-queue` | `GdsEnterQueue` | GdsQueueService | 进入队列 |
| `/api/gdsQueue/gds-exit-queue` | `GdsExitQueue` | GdsQueueService | 退出队列 |
| `/api/gdsQueue/gds-next-on-queue` | `GdsNextOnQueue` | GdsQueueService | 取队列下一条 |
| `/api/gdsQueue/gds-queue-agent-list` | `GdsQueueAgentList` | GdsQueueService | 队列坐席列表 |

---

## 终端 Terminal（`/api/terminal`，3 个）

| 路由 | SOAP 操作 | 上游服务 | 说明 |
|---|---|---|---|
| `/api/terminal/create-terminal-session` | `CreateTerminalSession` | TerminalService | 创建终端会话 |
| `/api/terminal/terminal` | `Terminal` | TerminalService | 发送终端指令 |
| `/api/terminal/end-terminal-session` | `EndTerminalSession` | TerminalService | 结束终端会话 |

---

## 系统 System（`/api/system`，4 个）

| 路由 | SOAP 操作 | 上游服务 | 说明 |
|---|---|---|---|
| `/api/system/ping` | `Ping` | SystemService | 探活 |
| `/api/system/info` | `Info` | SystemService | 服务信息 |
| `/api/system/time` | `Time` | SystemService | 服务时间 |
| `/api/system/cache` | `ExternalCacheAccess` | SystemService | 外部缓存访问 |

---

## 工具 Util（`/api/util`，24 个）

| 路由 | SOAP 操作 | 上游服务 | 说明 |
|---|---|---|---|
| `/api/util/calculate-tax` | `CalculateTax` | UtilService | 计算税费 |
| `/api/util/currency-conversion` | `CurrencyConversion` | UtilService | 币种换算 |
| `/api/util/credit-card-auth` | `CreditCardAuth` | UtilService | 信用卡授权 |
| `/api/util/mct-count` | `MctCount` | UtilService | 最短衔接时间计数 |
| `/api/util/mct-lookup` | `MctLookup` | UtilService | 最短衔接时间查询 |
| `/api/util/mco-create` | `MCOCreate` | UtilService | 创建 MCO |
| `/api/util/mco-exchange` | `MCOExchange` | UtilService | MCO 换开 |
| `/api/util/mco-issue` | `MCOIssue` | UtilService | MCO 开具 |
| `/api/util/mco-retrieve` | `MCORetrieve` | UtilService | 取回 MCO |
| `/api/util/mco-search` | `McoSearch` | UtilService | 搜索 MCO |
| `/api/util/mco-void` | `McoVoid` | UtilService | MCO 作废 |
| `/api/util/create-agency-fee-mco` | `CreateAgencyFeeMco` | UtilService | 创建代理费 MCO |
| `/api/util/create-airline-fee-mco` | `CreateAirlineFeeMco` | UtilService | 创建航司费 MCO |
| `/api/util/agency-service-fee-create` | `AgencyServiceFeeCreate` | UtilService | 创建代理服务费 |
| `/api/util/branded-fare-search` | `BrandedFareSearch` | UtilService | 品牌运价搜索 |
| `/api/util/branded-fare-admin` | `BrandedFareAdmin` | UtilService | 品牌运价管理 |
| `/api/util/upsell-admin` | `UpsellAdmin` | UtilService | 加价管理 |
| `/api/util/upsell-search` | `UpsellSearch` | UtilService | 加价搜索 |
| `/api/util/reference-data-retrieve` | `ReferenceDataRetrieve` | UtilService | 参考数取回 |
| `/api/util/reference-data-search` | `ReferenceDataSearch` | UtilService | 参考数搜索 |
| `/api/util/reference-data-update` | `ReferenceDataUpdate` | UtilService | 参考数更新 |
| `/api/util/content-provider-retrieve` | `ContentProviderRetrieve` | UtilService | 内容供应商取回 |
| `/api/util/mir-report-retrieve` | `MirReportRetrieve` | UtilService | MIR 报告取回 |
| `/api/util/find-employees-on-flight` | `FindEmployeesOnFlight` | UtilService | 航班员工查询 |

---

## 统计

| 域 | 路由数 | 其中别名（→UniversalRecord） |
|---|---|---|
| air | 29 | 2 |
| hotel | 10（含 2 自定义 DTO：search/details） | 2 |
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
| **合计** | **181** | **9** |

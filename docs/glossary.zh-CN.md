# 旅行行业 / GDS 术语表

**[English](./glossary.md)** | 简体中文

> 给第一次对接 Travelport UAPI 的开发者。看懂这些词，再看 [`architecture.md`](./architecture.zh-CN.md) 和 [`routing.md`](./routing.zh-CN.md) 会轻松很多。
> 术语以行业通行含义为准，括号里是本项目代码 / 路由里的对应位置。

---

## 基础概念

**GDS（Global Distribution System，全球分销系统）**
 travels / 航空 / 酒店等内容的中枢预订网络。主流三家：Amadeus、Sabre、Travelport（Galileo / Apollo / Worldspan 合并而来）。本项目对接的是 **Travelport**。

**UAPI / Universal API**
 Travelport 对外提供的编程接口总称。它把机票、酒店、火车、租车、客户档案等能力统一以 SOAP/XML 暴露，本项目把它包成 REST/JSON。

**SOAP / XML**
 一种基于 XML 的远程调用协议。UAPI 的请求要包在 `<SOAP:Envelope>` 里，且每个元素都带命名空间 URI（如 `http://www.travelport.com/schema/air_v54_0`）。本项目在 `pkg/client` 负责组装与解析。

**WSDL / XSD**
 WSDL 描述"有哪些 Service 和操作"，XSD 描述"请求 / 响应长什么样"。Travelport 按版本交付（`air_v54_0` 等），本项目归档在 `wsdl/`，并在构建期生成 Go 代码。

**PortType / 操作（Operation）**
 UAPI 里一个具体的可调用的动作，如 `AirSearch`、`HotelCreateReservation`。本项目的每个 HTTP 路由都对应一个 PortType。

**Provider（供应商 / 承运方）**
 实际提供库存的一方，通常指 GDS 背后的航空公司 / 酒店集团。`ProviderCode`（如 `1G` = Galileo）常在请求里指定。

**Branch / TargetBranch**
 Travelport 账号体系下的分支（对应一家代理点或子账号）。很多操作需要带 `TargetBranch` 指明"以哪个身份"发请求。本项目中它由调用方在请求体里提供。

---

## 预订相关

**PNR（Passenger Name Record，旅客订座记录）**
 一次预订的核心记录，包含旅客、航段、联系方式、票价等。传统上靠终端指令（如 `1`/`NM`/`SS`）编辑；本项目的 `sharedBooking` 域提供元素级 REST 操作。

**Universal Record（统一记录，UR）**
 Travelport 用来把**跨产品**（机票 + 酒店 + 火车 + 租车 + 被动段）预订聚合到同一条记录里的机制。创建 / 取消 / 统一修改都走 `UniversalRecordService`。这也是为什么本项目的 `/air/book`、`/hotel/cancel` 等"别名路由"最终落到 `UniversalFacade`（详见 [`architecture.md` §4](./architecture.zh-CN.md)）。

**Provider Reservation（供应商预订）**
 供应商侧对一条 UR 里某段产品的实际预订。本项目有 `provider-reservation-display-details`、`provider-reservation-divide` 等操作。

**Saved Trip（保存的行程）**
 一种"未真正出票 / 仅暂存"的行程草稿，可稍后取回继续处理。对应 `saved-trip-*` 一组操作。

**Passive Segment（被动段）**
 不在 GDS 里真正预订、只是"记账/备注"的一段（例如非 GDS 渠道订的酒店）。它没有独立 WSDL，只能经 `UniversalRecordService` 的 `Passive*` 操作创建 / 取消。本项目用 `passive` 域代理。

**CreateReservation / Cancel（创建预订 / 取消）**
 各产品的"落地预订"动作。注意机票/酒店/火车/租车的 create/cancel **不在各自产品 Service，而在 UniversalRecordService**（如 `AirCreateReservation`、`HotelCancel`），所以本项目把它们做成产品别名路由。

---

## 机票专用

**Air Search / Low Fare Search / Availability**
 航班检索的几种粒度：`AirSearch` 按指定航班，`LowFareSearch` 找最低价，`AvailabilitySearch` 查余座。

**Price / Reprice / Fare（计价 / 重新计价 / 运价）**
 把检索结果算成具体价格。`AirPrice` 计价，`AirReprice` 在已有预订上重算，`AirFareRules` 取运价规则。运价（Fare）是航司公布的价规组合。

**Ticketing / Void / Refund（出票 / 作废 / 退票）**
 `AirTicketing` 真正开票；`AirVoidDocument` 作废尚未结算的票；`AirRefund` 退票。票证（Document）指机票或 EMD。

**Exchange（换开）**
 已出票后的改期 / 升舱 / 重新开票，涉及差价计算：`AirExchangeQuote` 报价、`AirExchange` 执行、`AirExchangeTicketing` 换开出票。

**EMD（Electronic Miscellaneous Document，杂费电子票）**
 除机票外的附加服务凭证（如行李费、选座费）。`emd-issuance` / `emd-retrieve`。

**Seat Map（座位图）**
 取航班座位布局与占用情况。

**Merchandising / Upsell / Ancillary（辅营 / 加价 / 附加服务）**
 机票之外的付费服务（优先登机、额外行李等）。`air-merchandising-*`、`air-upsell-search`、`air-pre-pay`（预付）。

**Schedule Change（航班变动）**
 航司调整时刻后，`AckScheduleChange` 用于确认接收变动通知。

---

## 客户档案

**UProfile（Universal Profile，统一客户档案）**
 Travelport 的旅客 / 代理档案系统，存旅客偏好、支付方式、常旅客号等。`uprofile` 域是代理自有档案，`sharedUprofile` 是可与 Travelport 共享的档案。常见动作：search / retrieve / create / modify / delete，以及 field / tags / hierarchy / template 等子资源管理。

**Hierarchy Level（层级）**
 代理组织架构里的层级（公司 → 部门 → 分公司），用于权限与数据共享。

---

## 工具与杂项

**MCO（Miscellaneous Charges Order，杂费订单）**
 一种可结算的费单凭证，常用于代理费、改期费。`mco-*` 一组操作（create / issue / exchange / void / search / retrieve）。

**MCT（Minimum Connection Time，最短衔接时间）**
 同机场中转所需的最少时间，`mct-lookup` / `mct-count` 查询。

**Reference Data（参考数）**
 各种编码表（城市、机场、国家、币种、运价类型等），`reference-data-*` 查询 / 更新。

**GDS Queue（GDS 队列）**
 后台工作队列，用于把待处理 PNR 分派给坐席。`gdsQueue` 域提供 list / count / place / remove / next 等操作。

**Terminal（终端仿真）**
 直接发 GDS 原生指令（如 `HELP`、`*A`）。`terminal` 域提供会话的创建 / 发送 / 结束。

**Branded Fare（品牌运价）**
 航司把运价打包成带名称的产品（如"经济全价""轻享"），`branded-fare-*` 查询 / 管理。

**Agency Service Fee（代理服务费）**
 代理向客户收的服务费，`agency-service-fee-create` 创建。

---

## 本项目通用词

**Facade（用例层）**
 `pkg/usecase` 里每个域一个的结构体（如 `AirFacade`、`UniversalFacade`），把 REST 入参映射成 SOAP 请求并调用 Service。路由表里"SOAP 操作"一列写的就是 Facade 方法名。

**别名路由（Alias Route）**
 产品 URL（如 `/api/air/book`）代理到 `UniversalFacade` 的跨产品操作，既保留产品化直观 URL，又复用 UniversalRecord 的统一引擎。共 9 条。

**trace_id**
 贯穿整条请求链路的追踪 ID（UUID v4），同时出现在 HTTP 头 `X-Trace-Id`、日志、以及 SOAP 请求公共属性，用于排障时关联原始 XML 报文。

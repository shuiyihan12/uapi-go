# 更新日志（Changelog）

**[English](./CHANGELOG.md)** | 简体中文

## 版本规范

- **模块主版本固定为 `0`**（Go 仅在 major ≥ 2 时要求 `/vN` 路径后缀，故裸路径 `github.com/shuiyihan12/uapi-go` 始终合法）。tag 格式为 `v0.WSDL.PATCH`，其中 `WSDL` 为 Travelport WSDL 契约版本号（即 `wsdl/` 目录版本后缀，如 v55_0 契约对应 `55`），`PATCH` 为我方自增补丁号——例如 `v0.55.1` = 2026-08 发布。
- **只有打 git tag 才算发布新版本**；未打 tag 的提交属于开发态（daemon 显示 `dev`）。
- **补丁版本**递增 `PATCH` 段：如 `v0.55.2` 紧随 `v0.55.1`。补丁版本包含向后兼容的修复与新增；破坏性变更一律随 WSDL 位（minor）发布（如 `v0.56.0`）。
- 版本号经 `-ldflags "-X main.version=..."` 注入二进制，`uapi-go-daemon --version` / 启动日志可见。

---

## v0.55.1（2026-08）

上游契约 **v54_0 → v55_0 主版本升级** + 三项配套重构。**本版本不兼容旧版**，调用方必须按下方迁移说明改造。

### ⚠️ 破坏性变更（调用方必读）

1. **universal 酒店预订：`hotelProperty` 由 JSON 数组变为单个对象**
   - 影响接口：`POST /api/universal/hotel-create-reservation`（及 `/api/hotel/book` 别名，二者同引擎）。
   - 原因：上游 v55 契约将 `HotelProperty`/`HotelStay` 基数由 `0..99` 收窄为必填 1 个、`Guarantee` 等由 `0..99` 变 `0..1`。
   - 决策：网关**不做兼容垫层**，避免维护两套形态（理由见 `docs/architecture.md` ADR-003）。调用方将 `"hotelProperty": [ {...} ]` 改为 `"hotelProperty": { ... }` 即可。
   - 字段级迁移对照（`HotelCreateReservationReq`，v54 → v55）：

     | 字段 | v54 JSON 形态 | v55 JSON 形态 | 语义变化 |
     |---|---|---|---|
     | `hotelProperty` | 数组（0..99，可空） | **必填单对象**（1..1） | 一次请求创建一个酒店预订；多酒店改为多次调用 |
     | `hotelStay` | 数组（0..99，可空） | **必填单对象**（1..1） | 与 hotelProperty 一一对应 |
     | `guarantee` | 数组（0..99） | 可选单对象（0..1） | 至多一个担保 |
     | `guestInformation` | 数组（0..99） | 可选单对象（0..1） | 至多一组客人信息 |
     | `reservationName` | 数组（0..99） | 可选单对象（0..1） | 至多一个预订名 |
     | `hotelSpecialRequest` | 对象数组（`[{key, hotelRateDetailRef}]`） | **可选纯文本字符串**（0..1，类型简化为 `typeGeneralText`） | 不再是结构体/数组，直接传文本 |

2. **酒店 `/search`、`/details` JSON 契约随生成类型重建**
   - 请求：`transactionId` 入参**不再存在**（XSD 契约无此字段，旧手写模型自造；追踪由服务端注入 TraceId）。
   - 响应：由旧「精简版」摘要变为**完整 WSDL 生成模型**（字段集不同、信息更全）。
   - `nextResultReference` 在请求/响应中**可见**（详情翻页令牌，服务端已自动翻页，一般无需传入）。
   - 决策记录：ADR-005（手写 SOAP 模型全量退役，生成类型成为唯一真相源）。

3. **`sharedUprofile` 域依赖链修复**：补充归档 `uprofileCommon_v30_0` 契约并注册生成器，该域跨包引用由「兜底解析到最新 common 包」改为按真实命名空间解析。若此前对接过该域，建议重新联调验证。

4. **生成包名去版本化**（仅影响二次开发者）：`air55`→`air`、`hotel55`→`hotel` 等；`commonNN` 系因多版本并存保留版本号。此后契约升级（v55→v56）import 路径零改动（ADR-004）。

### 新增

- **`POST /api/sharedBooking/booking-vehicle-pnr-element`**：v55 新能力——在共享预订 PNR 上**新增 / 更新 / 删除车辆元素**（`addVehiclePnrElement` / `updateVehiclePnrElement` / `deleteVehiclePnrElement`），响应含 `vehicleRateChangedInfo` 车辆价变信息。对应上游新 PortType `BookingVehiclePnrElementPortType`。
- 对外接口总数 **180 → 181**。

### 变更

- 11 个核心域全量迁移至 v55_0 契约（air/common/cruise/gdsQueue/hotel/passive/rail/sharedBooking/universal/util/vehicle）；契约级 diff 结论与决策见 ADR-003。
- 酒店 `/search`、`/details` 手写 SOAP 模型（约 20 个 `Response*` 类型）与演示期死代码 `Search`/`Book` 全部删除，净删约 500 行；`/search`、`/details` 改经泛型 `registerPortHandler` 注册，8 个酒店操作共享同一条校验路径。
- 生成器清单改无版本包名；跨命名空间消歧类型前缀同步去版本（`Rail55Characteristic`→`RailCharacteristic`）。
- 手写酒店模型的 v54 命名空间 URI 升版问题随退役一并解决；全仓测试断言同步 v55。

### 修复

- 修复 `wsdl/` 目录只读权限导致交付无法落盘的问题。
- 修复 `sharedUprofile_v20_0` 交付缺失导致的依赖断裂（见破坏性变更 3）。

### 升级与验证

- 生成规模：2411 complexType / 22 包；`build` / `vet` / `gofmt` / 全部测试通过；daemon 冒烟（`/ready`、新路由校验链路）通过。
- 契约升级操作手册：`docs/development.md` §6；架构决策全记录：`docs/architecture.md` §9。

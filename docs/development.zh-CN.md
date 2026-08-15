# 二次开发指南（Development Guide）

**[English](./development.md)** | 简体中文

> 面向读者：想在本项目基础上**新增接口、新增上游域、升级 WSDL 契约或添加业务编排**的工程师。
> 先读 [`architecture.md`](./architecture.zh-CN.md) 理解分层与设计决策，再回到这里按场景操作。
> 路由 ↔ SOAP 操作对照表见 [`routing.md`](./routing.zh-CN.md)；GDS 术语见 [`glossary.md`](./glossary.zh-CN.md)。

---

## 0. 五分钟理解代码结构

```text
HTTP 请求
  └─ pkg/api/            路由注册 + JSON 编解码 + 错误 → HTTP 状态码
      └─ pkg/usecase/    业务用例（Facade）：校验、编排；纯透传域只做转发
          └─ pkg/manager/  服务实例的懒加载工厂（泛型 Get[T]）
              └─ pkg/services/<domain>/  PortType 方法 → SOAP 报文
                  └─ pkg/client/  Envelope 组装、鉴权透传、指标、Fault 解析
                      └─ Travelport UAPI (SOAP/XML)
```

配套的三件"万能设施"，新代码优先复用而不是另起炉灶：

| 设施 | 位置 | 作用 |
|---|---|---|
| `registerPortHandler[TReq,TRsp]` | `pkg/api/handler.go` | 一行注册一个 POST 路由：严格 JSON 解码（拒绝未知字段）、请求级鉴权/区域/trace 注入、统一超时与错误映射 |
| `CallPortType[T]` / `callPort[T]` | `pkg/client/enterprise.go` | 发起 SOAP 调用并把响应反序列化为强类型 `*T`，自动带指标与分级日志 |
| `manager.Get[T]` | `pkg/manager/service_manager.go` | 类型安全地取（并缓存）某个域的服务实例 |

---

## 1. 场景 A：暴露一个上游已提供的 PortType 操作（最常见）

以 Air 域为例，一个操作从 HTTP 到 SOAP 共四处改动。参考样板：`AirExchange`（`pkg/services/air/service.go`）。

> **生成包导入约定**：services/usecase 层同时 import 同名业务包（`pkg/services/air`）与生成包（`pkg/generated/air`），二者包名相同会冲突，因此生成包统一以 `xsd` 后缀别名导入（`airxsd "github.com/shuiyihan12/uapi-go/pkg/generated/air"`，hotel 域为既有的 `hotelxsd`）；`commonNN` 系无冲突，保持裸导入。

### 第 1 步：services 层加 PortType 方法

在 `pkg/services/air/service.go` 的 `AirServicePort` 接口与实现里各加一段：

```go
// AirXxx 对应 Air 服务的 AirXxxReq 操作。
AirXxx(ctx context.Context, req *airxsd.AirXxxReq) (*airxsd.AirXxxRsp, error)

func (s *AirService) AirXxx(ctx context.Context, req *airxsd.AirXxxReq) (*airxsd.AirXxxRsp, error) {
	ctx = prepareReq(ctx, req) // 注入 TraceId（尽力而为：请求实现 InjectInfrastructure 即写入）
	return callPort[airxsd.AirXxxRsp](s.client, ctx, "AirXxx", req)
}
```

请求/响应类型直接用 `pkg/generated/air` 里的生成类型——**不要手写请求模型**。

### 第 2 步：usecase 层加 Facade 方法

`pkg/usecase/air_facade.go`（纯透传样板）：

```go
// AirXxx 对应 AirServicePort.AirXxx。
func (f *AirFacade) AirXxx(ctx context.Context, req *airxsd.AirXxxReq) (*airxsd.AirXxxRsp, error) {
	svc, err := f.getService()
	if err != nil {
		return nil, err
	}
	ctx, _ = trace.Ensure(ctx)
	return svc.AirXxx(ctx, req)
}
```

> 透传 Facade 是刻意的「可折叠层」：现在没有业务逻辑，将来需要校验 / 编排时就地生长，不跨层改动。

### 第 3 步：api 层注册路由

`pkg/api/air_handler.go` 的 `RegisterRoutes` 里加一行：

```go
registerPortHandler(mux, apiBasePath+"/air/air-xxx", f.AirXxx)
```

命名规范：路由用 kebab-case（`air-xxx`），与方法名（`AirXxx`）一一对应。

### 第 4 步：文档与测试

- `docs/routing.md`：对应域的表格加一行（路由 | SOAP 操作 | 上游服务 | 说明），并更新该域与合计的路由数
- 如涉及非平凡的字段映射或校验，补表驱动测试（参考 `pkg/services/system/system_test.go` 的 Marshal 断言写法）
- `./scripts/build.sh all` 全绿后提交

---

## 2. 场景 B：接入一个全新的上游域

以接入新域 `foo`（WSDL 已交付）为例，改动集中在六处：

1. **归档契约**：WSDL/XSD 放入 `wsdl/foo_vNN_0/`，确认 `include`/`import` 引用完整
2. **生成结构体**：更新 `tools/airxsdgen` 的输入清单（若需新增 schema 文件）→ `./scripts/build.sh wsdl` → 产出 `pkg/generated/foo/`（单版本域用无版本包名；仅 common 系多版本并存保留版本号）
3. **services 层**：新建 `pkg/services/foo/service.go`——定义 `FooServicePort` 接口（与 WSDL 的 *PortType 一一对应）+ `FooService` 实现（构造 `client.EnterpriseSOAPClient`，每个方法 `prepareReq` + `callPort`），参考 `pkg/services/rail/service.go`（小域样板）
4. **manager 注册**：`pkg/manager/service_manager.go` 两处——`buildServiceFactories` 加一条工厂项、`serviceSuffix` 加 `case "foo": return "FooService", nil`（未知键会显式报错，漏配会在首次调用时暴露）
5. **usecase + api 层**：新建 `pkg/usecase/foo_facade.go` 与 `pkg/api/foo_handler.go`（样板同场景 A），并在 `cmd/daemon/main.go` 装配：`NewFooFacade` → `NewFooHandler(...).RegisterRoutes(mux)`
6. **文档**：README 能力表、`docs/routing.md` 新增域章节、必要时 `docs/architecture.md` §2.1 分层表

---

## 3. 场景 C：上游契约（WSDL/XSD）升级

完整流程见 [`architecture.md` §6](./architecture.zh-CN.md)。要点：

```bash
# 新契约归档进 wsdl/ 后：
./scripts/build.sh wsdl   # 重新生成 pkg/generated
./scripts/build.sh all    # 构建 + 测试 + lint
git diff --stat pkg/generated/   # 人工审查生成代码 diff（字段增删、类型变化）
```

- `pkg/generated/` 是生成资产，**diff 只读不手改**
- 关注生成器两类历史缺陷的复发迹象：跨命名空间 `ref` 丢字段、`xs:date` 序列化异常（详见 `architecture.md` §3.4）
- 上游删字段 / 改类型会直接导致编译错误——这正是"契约前置到编译期"的设计意图，顺势修适配层即可

---

## 4. 场景 D：添加业务编排（校验、翻页、聚合）

当透传不够用（需要多步 SOAP 调用、输入校验或友好输出）时，在 **usecase 层**生长逻辑。现成范本：`pkg/usecase/hotel_facade.go`。

- **输入校验**：参考 `validateHotelStay`（日期合法性）与 `normalizeHotelPortReq`（必填字段归一化），校验失败返回 `*usecase.ValidationError`，API 层会自动映射为 400 `INVALID_REQUEST`
- **多次调用 / 翻页**：参考 `HotelDetails` 的 `NextResultReference` 翻页循环——注意区分两类超时：单次 SOAP 调用受 `UAPI_REQUEST_TIMEOUT` 约束，多次调用的累计耗时用业务级常量兜底（如 `defaultHotelDetailsPageTimeout = 40s`、`defaultHotelDetailsMaxPages = 20`）
- **输出形态**：编排结果建议包一层手写 Output 结构（如 `HotelSearchAvailabilityOutput`，仅含一个 `Response` 字段），保持对外 JSON 顶层形状稳定，避免 `map[string]interface{}` 透传
- **注意**：编排型接口同样可以走 `registerPortHandler`——只要 facade 方法签名形如 `func(ctx, *Req) (*Rsp, error)`（酒店 8 个路由全部如此注册），无需手写 handler

---

## 5. 横切约定（改任何代码前先读）

| 约定 | 说明 |
|---|---|
| 鉴权一切透传 | 网关不持有 Travelport 凭据；`Authorization` / `X-UAPI-Region` 由调用方逐请求携带，经 `pkg/requestctx` 从 HTTP 头透传到 SOAP 调用。`/health` 同理（探活方带凭据） |
| 不做重试 | GDS 写操作（出票/预订）不幂等，任何失败直接透出，由调用方决定补偿 |
| 错误模型 | 上游业务/系统错误经 `*client.SOAPFaultError` 透出；`pkg/api.writeError` 按 `ErrorInfo/Type` 统一映射 400/422/502/504。中间层不要吞错或转成字符串 |
| trace_id 全链路 | `trace.Ensure(ctx)` 自动生成 UUID 并贯穿 HTTP 头、日志、SOAP 报文；新代码不要自造追踪 ID |
| 生成代码禁改 | `pkg/generated/` 只能经 `./scripts/build.sh wsdl` 再生 |
| 严格 JSON | 请求解码 `DisallowUnknownFields`；新接口的请求字段命名 snake_case，与生成类型的 json tag 对齐 |
| 注释风格 | 中文文档注释，导出符号必须注释"做什么、为什么" |

---

## 6. 测试与调试

### 测试

```bash
go test ./...                          # 全部单元测试（默认不触网，CI 可跑）
UAPI_INTEGRATION=1 \
UAPI_TEST_AUTHORIZATION="Basic xxx" \
go test ./test/... -run TestServiceManager -v   # 真实上游集成路径
```

- 单元测试贴近被测包（如 `pkg/services/system/system_test.go`）；集成测试放 `test/`
- 需要真实凭据的用例一律以 `UAPI_INTEGRATION` 门控
- 想伪造上游响应：`client.SOAPConfig.Transport` 支持注入 `http.RoundTripper`（用法见 `pkg/client/soap_test.go`）

### 调试

1. **拿 trace_id**：任一响应头 `X-Trace-Id`（或自己传入 `X-Trace-Id`）
2. **看原始报文**：日志中按 trace_id 过滤 `[GDS REQUEST]` / `[GDS RESPONSE]`，可逐字段核对 SOAP 报文
3. **看指标**：`GET /metrics`（Prometheus 格式，`uapi_requests_total` 按 service/operation/status 维度）
4. **本地起服务**：

```bash
set -a && source .env && set +a
go run ./cmd/daemon          # 默认 :8080，业务与 /health、/metrics 同端口

curl -X POST http://localhost:8080/api/system/ping \
  -H "Authorization: Basic xxxx" \
  -H "X-UAPI-Region: apac" \
  -H "Content-Type: application/json" -d '{}'
```

5. **健康检查排障**：`curl -H "Authorization: Basic xxxx" http://localhost:8080/health`——不带凭据时上游拒绝、返回 503 属预期行为（网关本身无凭据可验证）

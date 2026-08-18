# 更新日志

本文件记录项目的所有重要变更。

版本号采用 `v0.WSDL.PATCH` 规则：
- `WSDL` 对应 Travelport UAPI 的模式生成版本（例如 `55` 对应 `v55_0`）。
- `PATCH` 用于不改动生成 WSDL 面的 SDK 级修复与改进。

条目按时间从早到晚排列。

---

## [v0.55.1] — 2026-08-16

初始提交。Travelport UAPI 的 Go SDK / REST 网关首个版本。

- 项目骨架：`cmd/` 入口（`daemon` REST+管理端服务、`generator`、`generate`），`pkg/` 各层（`api`、`usecase`、`client`、`manager`、`services`、`requestctx`、`trace`），`pkg/generated` WSDL/XSD 生成代码，`internal/` 基础设施，`tools/airxsdgen` 生成器，`wsdl/` 上游输入文件，以及 `scripts/` 构建脚本。
- 容器构建资源（`Dockerfile`、`docker-compose.yml`、`deploy/k8s/`、`cmd/healthcheck/`），采用 distroless 非 root 运行时镜像。
- 文档与贡献指南（`AGENTS.md`、`README`、`docs/` 架构/路由/术语）。

## [v0.55.2] — 2026-08-16

Go SDK 对外接口、移除游离二进制、修复 Docker 镜像标签。

- **chore(version)**：采用 `v0.WSDL.PATCH` 模块版本规则。
- **refactor(logging)**：将日志器提升为公开的 `pkg/logging` 包。
- **refactor(client)**：通过 `SOAPConfig.Metrics` 使 SOAP 指标可插拔。
- **feat(sdk)**：新增 `sdk` 入口包，提供统一的客户端对外接口。
- **ci**：守护 SDK 导入面；**docs**：版本纪律。
- **docs(sdk)**：新增可运行示例（`ping`、`hotel-search`）。
- **chore(sdk)**：打磨选项、CI 守护与版本纪律。
- 移除游离的二进制产物，并修复发布流水线中的 Docker 镜像标签。

## [v0.55.3] — 2026-08-19

SOAP 请求信封清理与命名空间提升。

### 变更
- **SOAP 信封前缀**：信封现使用 `soapenv` 前缀（`xmlns:soapenv="http://schemas.xmlsoap.org/soap/envelope/"`），以对齐官方 Travelport UAPI 示例。
- **移除空 Header**：省略空的 `<soapenv:Header/>`。它不含任何内容（鉴权走 HTTP `Authorization` 头、链路追踪 ID 走请求体的 `TraceId` 属性），移除后可节省少量字节且行为不变。
- **命名空间提升**：请求体的命名空间现被提升到请求根元素上一次性声明为可读前缀——例如 `xmlns:hotel="http://www.travelport.com/schema/hotel_v55_0"`、`xmlns:common="http://www.travelport.com/schema/common_v55_0"`——取代 `encoding/xml` 原本在每个元素上重复输出的默认 `xmlns="..."`。输出现与官方示例结构一致（一次声明、子节点使用前缀），且体积更小。

### 说明
- 提升变换是通用且基于 URI 驱动的：它依据请求体中实际出现的命名空间 URI 工作，因此对所有 UAPI 服务（air/rail/vehicle/hotel/universal/…）及所有模式版本均无需改动即可生效。命名空间前缀由 URI 尾段去除版本后缀得到（`hotel_v55_0` → `hotel`、`common_v55_0` → `common`），未来升级到 `v56_0` 重新生成时会自动处理，无需改代码。
- 响应解析路径、`pkg/generated/` 下的生成代码，以及各 service 的业务逻辑均未改动。
- 由 `pkg/client/soap_test.go` 中的 `TestBuildEnvelopeNamespaceHoist` 测试覆盖，锁定信封输出形态。

# 贡献指南（Contributing Guide）

**[English](./CONTRIBUTING.md)** | 简体中文

感谢你考虑为 uapi-go 做贡献！这是一个把 Travelport Universal API（SOAP/XML）封装成 REST/JSON 网关的 Go 项目。本文说明开发环境、工作流与红线；架构与设计决策见 [`docs/architecture.md`](./docs/architecture.md)，二次开发的分步教程见 [`docs/development.md`](./docs/development.md)。

## 1. 环境要求

- Go `1.25.0`+（以 `go.mod` 为准）
- 可访问公网（拉取 Go modules；运行集成测试还需要可访问 Travelport 测试环境）
- 无需 Docker（除非你想容器化部署）

## 2. 开发工作流

```bash
# 1. 克隆并安装依赖
git clone <your-fork-url> && cd uapi-go
./scripts/build.sh deps

# 2. 开发迭代（按需）
./scripts/build.sh build     # 构建
./scripts/build.sh test      # 测试
./scripts/build.sh lint      # gofmt 检查 + go vet

# 3. 提交 PR 前的完整校验（生成 + 构建 + 测试 + lint）
./scripts/build.sh all
```

PR 检查清单：

- [ ] `./scripts/build.sh all` 全绿
- [ ] 行为变更附有测试证据（命令与输出）或示例请求/响应
- [ ] 新增接口同步更新了 `docs/routing.md`（路由 ↔ SOAP 操作对照表）
- [ ] `pkg/generated/` 只有经 `./scripts/build.sh wsdl` 产生的变更（绝不手改）

## 3. 代码规范

- 遵循 Go 默认风格：tab 缩进、`gofmt` 格式化、`go vet` 零告警
- 导出标识符用 `PascalCase`，包内辅助函数用 `camelCase`，包名小写且与域对齐（`hotel`、`universal`、`util`…）
- 注释使用中文，导出符号须有文档注释（说明"做什么、为什么"，而非复述代码）
- 分层单向依赖：`pkg/api → pkg/usecase → pkg/services → pkg/client`，禁止反向 import；新增能力优先复用既有泛型设施（`registerPortHandler`、`CallPortType[T]`、`manager.Get[T]`）
- 错误处理：上游业务错误经 `*client.SOAPFaultError` 透出，由 `pkg/api.writeError` 统一映射 HTTP 状态码；不要在中间层吞错或改动错误类型

## 4. 测试规范

- 框架：标准库 `testing`；测试放 `*_test.go`，服务逻辑优先表驱动
- 命名：`TestXxx` / `BenchmarkXxx`
- 需要真实上游的用例必须以 `UAPI_INTEGRATION=1` 门控（凭据经 `UAPI_TEST_AUTHORIZATION` 提供），默认 CI 不跑
- 安全红线：**任何真实凭据、生产 token 都不得出现在代码、测试、日志样例或文档中**

## 5. 提交信息

简短祈使句，一屏内说明意图，例如：

```
Add hotel media-links port passthrough
Fix TLS verify default in service manager
Regenerate WSDL client code for v54 contracts
```

一个提交只做一件事（生成、服务逻辑、文档、工具分开提交）。

## 6. 契约（WSDL/XSD）升级

上游契约变更必须走构建期重新生成，流程见 [`docs/architecture.md` §6](./docs/architecture.md)：归档 `wsdl/` → 更新生成清单 → `./scripts/build.sh wsdl` → `all` → 人工审查 `pkg/generated` diff。

## 7. 许可

- 本项目以 [Apache License 2.0](./LICENSE) 开源；**提交 PR 即表示你同意以 Apache-2.0 许可你的贡献**（无需额外 CLA）
- 新增源文件建议在文件头附 SPDX 标识：`// SPDX-License-Identifier: Apache-2.0`
- `wsdl/` 下的上游契约工件版权归 Travelport，不要修改或再分发到许可范围之外

## 8. 行为准则

- 保持讨论聚焦于技术；对评审意见对事不对人
- 提问前先查 [`docs/glossary.md`](./docs/glossary.md)（GDS 术语表）与 [`docs/architecture.md`](./docs/architecture.md)

# Go 后台服务模板

一个可直接运行的 **Go 后台服务骨架**，附带完整的 AI 辅助开发配置（`.codebuddy/`）。
用它作为新项目的起点：clone → 改模块名 → 写业务代码。

## 特性

- 极少依赖（仅 `gopkg.in/yaml.v3`），`go run` 立即可跑
- 标准 Go 项目布局（`cmd/` + `internal/` + `pkg/`）
- HTTP 服务骨架：路由注册、超时配置、优雅关闭、健康检查
- 配置支持 YAML 文件 + 环境变量覆盖（env 优先）
- 内置 AI 开发工作流：需求澄清 → 系统设计 → 代码生成 → 测试生成 → 代码评审 → 代码提交

## 快速开始

```bash
# 修改模块名（把 your-org 替换为你的组织）
go mod edit -module github.com/<your-org>/<your-repo>

# 运行
go run ./cmd/server

# 构建
make build

# 测试
make test
```

服务默认监听 `:8080`：

| 路由 | 说明 |
|------|------|
| `GET /healthz` | 健康检查，返回 `{"status":"ok"}` |
| `GET /` | 服务信息 |

## 配置

**唯一配置来源是 YAML 文件**，启动时用 `-config` 指定路径（默认 `configs/config.yaml`）:

```bash
go run ./cmd/server -config configs/config.yaml
go run ./cmd/server -config /etc/myapp/config.yaml
```

```yaml
# configs/config.yaml
server:
  addr: ":8080"
  read_timeout: 5s
  write_timeout: 10s
  shutdown_timeout: 10s
log:
  level: info    # debug | info | warn | error
  format: json   # json | text
```

| 配置项 | 默认值 | 说明 |
|--------|--------|------|
| `server.addr` | `:8080` | 监听地址 |
| `server.read_timeout` | `5s` | 读超时 |
| `server.write_timeout` | `10s` | 写超时 |
| `server.shutdown_timeout` | `10s` | 优雅关闭等待时间 |
| `log.level` | `info` | 日志级别 |
| `log.format` | `json` | 日志格式 |

要点：

- 配置文件**缺失或不可读会直接启动失败**，不会静默用默认值跑起来
- 文件里省略的字段回落默认值；**未知字段会直接报错**，避免拼写错误被忽略
- 时间字段写 `5s` / `300ms` 这类 Go duration 字符串
- 不同环境用不同文件（`config.dev.yaml` / `config.prod.yaml`），而不是环境变量

## 目录结构

| 目录 | 说明 |
|------|------|
| `cmd/server` | 程序入口，`main` 包 |
| `internal/config` | 配置加载（YAML + 环境变量） |
| `internal/server` | HTTP 服务与路由 |
| `pkg/` | 可被外部项目复用的公共库 |
| `configs/` | 配置文件（运行时读取） |
| `scripts/` | 构建与运维脚本 |
| `docs/design/<feature>/` | 设计文档（`spec.md`、`tasks.md`） |
| `test/` | 跨包集成测试与测试数据 |
| `.codebuddy/` | AI 协作配置（rules、skills、agents、commands） |

## 开发约定

- 代码规范见 `.codebuddy/skills/std-company-go/`
- 新增业务代码放在 `internal/`，只有确定要对外复用时才放 `pkg/`
- 每个需求先写 `docs/design/<feature>/spec.md`，再进入编码

## AI 工作流

`.codebuddy/` 中预置了完整的 AI 协作流程，按需触发：

| 阶段 | 入口 | 产物 |
|------|------|------|
| 需求澄清 | `workflow-requirements-clarification` | `spec.md` 背景/目标/需求 |
| 系统设计 | `workflow-system-design` | `spec.md` 设计章节 |
| 代码生成 | `workflow-code-generation` | 代码 + `tasks.md` |
| 测试生成 | `workflow-test-generation` | `_test.go` |
| 代码评审 | `workflow-code-review` | 评审报告 |
| 代码提交 | `workflow-code-submission` | commit |

## Make 目标

| 命令 | 说明 |
|------|------|
| `make build` | 编译到 `bin/server` |
| `make run` | 直接运行 |
| `make test` | 运行测试 |
| `make check` | `go vet` + `go test`（提交前必跑） |
| `make fmt` | 格式化代码 |

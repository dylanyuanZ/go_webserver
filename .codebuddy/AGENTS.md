# go_webserver

## 项目概述

go_webserver 是一个 **Go 语言 AI 辅助开发模板项目**。它本身不含业务代码，提供的是一套可复用的工程骨架 + AI 协作配置（`.codebuddy/`），用于快速启动新的 Go 后台服务项目。

在此模板上开发时，请按本文的目录约定组织代码，并遵循 `.codebuddy/` 中定义的规则与工作流。

## 目录结构

| 目录 | 说明 |
|------|------|
| `cmd/` | 程序入口，每个子目录一个可执行程序（如 `cmd/server`） |
| `internal/` | 私有业务代码，外部项目不可导入 |
| `pkg/` | 可被外部项目复用的公共库 |
| `configs/` | 配置文件模板与默认配置 |
| `scripts/` | 构建、部署、运维脚本 |
| `docs/` | 项目文档；设计文档放 `docs/design/<feature>/`（`spec.md`、`tasks.md`） |
| `test/` | 跨包的集成测试与测试数据 |
| `web/` | 静态资源、模板等前端产物（按需启用） |

## 技术栈

- **语言**: Go
- **构建**: `go build` / `make build`
- **测试**: `go test`（标准库 + testify）
- **文档**: Markdown（中文）

## 编码约定

- 语言版本、格式化、命名、错误处理等规范见 `.codebuddy/skills/std-company-go/`
- 通用编码最佳实践见 `.codebuddy/skills/bp-coding-best-practices/`

## AI 工作流

本项目配置了完整的 AI 辅助开发工作流：

```
需求澄清 → 系统设计 → 代码生成 → 测试生成 → 代码评审 → 代码提交
```

详见 `.codebuddy/skills/` 和 `.codebuddy/rules/` 目录。

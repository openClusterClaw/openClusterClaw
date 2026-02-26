# 测试规范

本目录包含 Open Cluster Claw 项目的测试规范文档。

## 文档列表

| 文档 | 说明 |
| ------ | ------ |
| [后端测试规范](./后端测试规范.md) | 后端测试规范，包括 Go testing、testify、mock |

## 快速开始

### 运行所有测试

```bash
go test -v ./...
```

### 运行指定包测试

```bash
go test -v ./internal/service/...
```

### 生成覆盖率报告

```bash
go test -coverprofile=coverage.out ./...
go tool cover -html=coverage.out
```

## 测试覆盖目标

| 类型 | 目标覆盖率 |
| ---- | ---------- |
| Service 层 | 80% |
| Repository 层 | 70% |
| API 层 | 60% |
| Domain 层 | 85% |

## AI 生成测试

使用测试规范中的 AI Prompt 模板，可以快速为新功能生成测试代码。

详见 [后端测试规范.md](./后端测试规范.md)
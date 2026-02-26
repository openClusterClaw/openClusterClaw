# Open Cluster Claw

多实例、多类型 Claw 统一控制平面（Control Plane）

## 快速开始

### 后端

```bash
# 安装依赖
go mod download

# 运行服务
go run cmd/controlplane/main.go
```

### 前端

```bash
cd frontend
pnpm install
pnpm dev
```

## 项目结构

```
openClusterClaw/
├── cmd/                 # 应用入口
│   └── controlplane/    # 控制平面启动入口
├── internal/            # 内部包
│   ├── api/            # API 路由和 Handler
│   ├── service/        # 业务服务层
│   ├── domain/         # 领域模型
│   ├── repository/     # 数据访问层
│   ├── model/          # 数据库模型
│   └── middleware/     # 中间件
├── config/             # 配置管理
├── migrations/         # 数据库迁移
├── docs/              # 文档
└── pkg/               # 通用工具
```

## 技术栈

### 后端
- Go 1.21+
- Gin (Web 框架)
- PostgreSQL (数据库)
- Kubernetes (运行时)
- pgx/v5 (数据库驱动)

### 前端
- React 18
- TypeScript
- Vite
- Ant Design 5
- pnpm

## API 文档

查看 `docs/API.md` 获取详细的 API 文档

## 开发文档

- [需求文档](requirements.md)
- [架构文档](architecture.md)
- [开发规范](docs/开发规范/)

## License

MIT

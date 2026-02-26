# Open Cluster Claw 开发计划

> 基于架构设计文档 (`architecture.md`) 与当前代码实现状态的对比分析
>
> 更新日期: 2026-02-26

---

## 一、实现状态概览

### 整体进度: 约 20% 完成

| 架构层级 | 完成度 | 说明 |
|---------|--------|------|
| API Gateway Layer | 0% | 未实现 |
| Control Plane Core | 30% | 实例管理部分完成 |
| Config Manager | 10% | 适配器完成，模板管理未实现 |
| Resource Orchestrator | 40% | K8S Pod 完成，其他运行时未实现 |
| Storage Layer | 20% | SQLite 完成，对象存储未实现 |
| Runtime Layer (K8S) | 50% | Pod/ConfigMap 完成，PVC 未实现 |
| Skill & Plugin Manager | 0% | 未实现 |
| User Portal (Frontend) | 25% | 基础 UI 完成，缺少管理界面 |
| 多租户模型 | 10% | 表结构完成，逻辑未实现 |

---

## 二、已实现功能清单

### 2.1 控制面核心 (Control Plane Core)

#### InstanceManager (实例管理器)
| 功能 | 状态 | 文件位置 |
|------|------|---------|
| 实例创建 (Create) | ✅ | `internal/service/instance.go` |
| 实例查询 (Get) | ✅ | `internal/service/instance.go` |
| 实例列表 (List) | ✅ | `internal/service/instance.go` |
| 实例更新 (Update) | ✅ | `internal/service/instance.go` |
| 实例删除 (Delete) | ✅ | `internal/service/instance.go` |
| 实例启动 (Start) | ✅ | `internal/service/instance.go` |
| 实例停止 (Stop) | ✅ | `internal/service/instance.go` |
| 实例重启 (Restart) | ✅ | `internal/service/instance.go` |
| 状态同步 (Sync) | ✅ | `internal/service/instance.go:syncPodStatus()` |

#### 数据库层 (Repository)
| 功能 | 状态 | 文件位置 |
|------|------|---------|
| InstanceRepository | ✅ | `internal/repository/instance.go` |
| 数据库迁移 (Migrations) | ✅ | `cmd/controlplane/main.go:runMigrations()` |
| SQLite 支持 | ✅ | `cmd/controlplane/main.go` |

#### API 层
| 功能 | 状态 | 文件位置 |
|------|------|---------|
| RESTful API Router | ✅ | `internal/api/router.go` |
| 实例管理 Handler | ✅ | `internal/api/handler.go` |
| 健康检查 | ✅ | `internal/api/router.go` |

---

### 2.2 适配器引擎 (Adapter Engine)

| 功能 | 状态 | 文件位置 |
|------|------|---------|
| ClawAdapter 接口定义 | ✅ | `internal/adapter/adapter.go` |
| UnifiedConfig 模型 | ✅ | `internal/adapter/adapter.go` |
| OpenClawAdapter | ✅ | `internal/adapter/openclaw.go` |
| Adapter Factory | ✅ | `internal/adapter/factory.go` |

---

### 2.3 Kubernetes 运行时集成

| 功能 | 状态 | 文件位置 |
|------|------|---------|
| K8S 客户端初始化 | ✅ | `internal/runtime/k8s/client.go` |
| Pod Manager | ✅ | `internal/runtime/k8s/pod.go` |
| ConfigMap Manager | ✅ | `internal/runtime/k8s/configmap.go` |
| Pod 创建/删除/查询 | ✅ | `internal/runtime/k8s/pod.go` |
| Pod 状态查询 | ✅ | `internal/runtime/k8s/pod.go` |
| Pod 日志获取 | ✅ | `internal/runtime/k8s/pod.go` |
| Pod 事件查询 | ✅ | `internal/runtime/k8s/pod.go` |
| Pod 就绪等待 | ✅ | `internal/runtime/k8s/pod.go` |

---

### 2.4 前端 (Frontend)

| 功能 | 状态 | 文件位置 |
|------|------|---------|
| 项目初始化 (Vite + React + TS) | ✅ | `frontend/` |
| Ant Design UI 框架 | ✅ | `frontend/package.json` |
| 实例列表页面 | ✅ | `frontend/src/pages/InstanceList.tsx` |
| 创建实例模态框 | ✅ | `frontend/src/components/instances/CreateInstanceModal.tsx` |
| 实例详情页面 | ✅ | `frontend/src/pages/InstanceDetail.tsx` |
| API 客户端 | ✅ | `frontend/src/api/instance.ts` |
| 状态管理 (Zustand) | ✅ | `frontend/src/store/instance.ts` |
| 路由配置 (React Router) | ✅ | `frontend/src/main.tsx` |
| 静态文件嵌入 | ✅ | `internal/embed/frontend.go` |

---

### 2.5 配置管理

| 功能 | 状态 | 文件位置 |
|------|------|---------|
| 配置加载 (Viper) | ✅ | `config/config.go` |
| 配置文件 (YAML) | ✅ | `config/config.yaml` |
| 数据库配置 | ✅ | `config/config.go` |
| K8S 配置 | ✅ | `config/config.go` |
| 日志配置 | ✅ | `config/config.go` |

---

### 2.6 数据库表结构

| 表 | 状态 | 说明 |
|---|------|------|
| tenants | ✅ | 租户表结构已创建 |
| projects | ✅ | 项目表结构已创建 |
| config_templates | ✅ | 配置模板表结构已创建 |
| claw_instances | ✅ | 实例表结构已创建 |

---

## 三、未实现功能清单

### 3.1 API Gateway Layer (0%)

| 功能 | 优先级 | 说明 |
|------|--------|------|
| 统一认证 (OAuth2 / JWT) | P1 | 请求认证与授权 |
| RBAC 权限控制 | P1 | 基于角色的访问控制 |
| 租户上下文注入 | P1 | 请求自动携带租户信息 |
| API 限流保护 | P2 | 防止 API 滥用 |
| 请求审计日志 | P2 | 操作审计与追溯 |
| WebSocket 支持 | P2 | 流式输出和实时通信 |
| 移动端 API Gateway | P3 | Mobile App 接入支持 |

---

### 3.2 控制面核心 - 未完成项

#### ConfigManager (配置管理器)
| 功能 | 优先级 | 说明 |
|------|--------|------|
| 配置模板 CRUD | P1 | 模板的增删改查 |
| 配置版本控制 | P1 | 配置历史版本管理 |
| 批量配置下发 | P1 | 多实例批量更新配置 |
| 配置一致性校验 | P2 | 验证配置正确性 |
| 灰度发布 | P2 | 分阶段发布配置 |
| 配置回滚 | P2 | 回退到历史版本 |
| ConfigTemplateRepository 实现 | P1 | 数据访问层 |

#### UsageMonitor (使用监控器)
| 功能 | 优先级 | 说明 |
|------|--------|------|
| 实例级监控 (CPU/Memory/重启次数) | P1 | 实例健康监控 |
| 使用级监控 (Token/调用次数/错误率) | P1 | 使用统计与计费 |
| 指标聚合与存储 | P1 | 时序数据存储 |
| 监控 API 端点 | P1 | 查询监控数据 |

#### PolicyEngine (策略引擎)
| 功能 | 优先级 | 说明 |
|------|--------|------|
| 租户资源配额控制 | P1 | CPU/Memory/实例数限制 |
| 访问策略执行 | P2 | 策略规则引擎 |
| 自动扩缩容策略 | P3 | 基于使用量的自动调整 |

---

### 3.3 Skill & Plugin Manager (0%)

| 功能 | 优先级 | 说明 |
|------|--------|------|
| Skill CRUD API | P1 | 技能的增删改查 |
| Plugin CRUD API | P1 | 插件的增删改查 |
| Skill 版本管理 | P2 | 技能版本控制 |
| Plugin 版本管理 | P2 | 插件版本控制 |
| 批量升级 | P2 | 多实例批量升级 |
| 启用/禁用控制 | P1 | 开关控制 |
| 内网 Skill Registry | P1 | 技能仓库 |
| Skill Distributor | P2 | 统一推送/拉取机制 |

---

### 3.4 Resource Orchestrator Layer - 未完成项

| 功能 | 优先级 | 说明 |
|------|--------|------|
| Docker 运行时 | P2 | 直接 Docker 支持 |
| Podman 运行时 | P3 | Podman 支持 |
| 物理机运行时 | P3 | 裸机部署支持 |
| 虚拟机运行时 | P3 | 虚拟机部署支持 |
| 云厂商 API 集成 | P3 | AWS/阿里云等云服务 |
| 自动扩缩容 | P2 | 水平扩缩容逻辑 |

---

### 3.5 Storage Layer - 未完成项

| 功能 | 优先级 | 说明 |
|------|--------|------|
| 对象存储集成 (MinIO/S3) | P1 | 文件存储与共享 |
| PVC 管理 | P1 | Kubernetes 持久卷 |
| 文件上传/下载 API | P1 | 文件管理接口 |
| 在线文件编辑 | P2 | Web 编辑器 |
| 文件共享 (多实例) | P2 | 云盘模式 |
| 数据与实例解耦 | P1 | 数据持久化独立于生命周期 |

---

### 3.6 前端 - 未完成页面

| 页面 | 优先级 | 说明 |
|------|--------|------|
| 租户管理页面 | P1 | 租户列表/创建/编辑 |
| 项目管理页面 | P1 | 项目列表/创建/编辑 |
| 配置模板管理页面 | P1 | 配置模板 CRUD |
| 配置版本历史页面 | P2 | 版本对比与回滚 |
| 技能市场页面 | P1 | 技能浏览与安装 |
| 插件管理页面 | P1 | 插件安装与配置 |
| 监控仪表板 | P1 | 实例指标可视化 |
| 使用统计页面 | P2 | 资源使用统计 |
| 日志查看页面 | P1 | 实例日志实时查看 |
| 事件查看页面 | P2 | 实例事件查看 |
| 文件管理页面 | P1 | 实例文件管理 |
| 用户认证页面 | P1 | 登录/登出 |
| 权限管理页面 | P2 | 用户与角色管理 |

---

### 3.7 多租户模型 - 未完成项

| 功能 | 优先级 | 说明 |
|------|--------|------|
| 租户隔离逻辑 | P1 | Namespace 级隔离 |
| 租户配额检查 | P1 | 创建实例前检查资源配额 |
| TenantRepository 实现 | P1 | 租户数据访问层 |
| ProjectRepository 实现 | P1 | 项目数据访问层 |
| 租户默认配置覆盖 | P2 | 全局配置 + 租户覆盖 |

---

### 3.8 API 缺失端点

#### 实例管理 API
| 端点 | 优先级 | 说明 |
|------|--------|------|
| `POST /instances/:id/kill` | P2 | 强制终止 |
| `POST /instances/:id/console` | P3 | Console 登录 |
| `GET /instances/:id/logs` | P1 | 日志查询 |
| `WS /instances/:id/logs/stream` | P1 | 实时日志流 |
| `GET /instances/:id/metrics` | P1 | 实例指标 |
| `WS /instances/:id/metrics/stream` | P2 | 实时指标流 |
| `GET /instances/:id/events` | P2 | 实例事件 |

#### 配置管理 API
| 端点 | 优先级 | 说明 |
|------|--------|------|
| `GET /configs` | P1 | 配置模板列表 |
| `POST /configs` | P1 | 创建配置模板 |
| `GET /configs/:id` | P1 | 配置模板详情 |
| `PUT /configs/:id` | P1 | 更新配置模板 |
| `DELETE /configs/:id` | P1 | 删除配置模板 |
| `POST /configs/:id/publish` | P1 | 发布配置 |
| `POST /configs/:id/rollback` | P2 | 回滚配置 |
| `GET /configs/:id/versions` | P2 | 版本历史 |
| `POST /configs/:id/validate` | P1 | 验证配置 |

#### 租户管理 API
| 端点 | 优先级 | 说明 |
|------|--------|------|
| `GET /tenants` | P1 | 租户列表 |
| `POST /tenants` | P1 | 创建租户 |
| `GET /tenants/:id` | P1 | 租户详情 |
| `PUT /tenants/:id` | P1 | 更新租户 |
| `DELETE /tenants/:id` | P1 | 删除租户 |
| `GET /tenants/:id/quota` | P1 | 查询配额 |
| `PUT /tenants/:id/quota` | P1 | 更新配额 |

#### 项目管理 API
| 端点 | 优先级 | 说明 |
|------|--------|------|
| `GET /projects` | P1 | 项目列表 |
| `GET /projects/:id` | P1 | 项目详情 |
| `POST /projects` | P1 | 创建项目 |
| `PUT /projects/:id` | P1 | 更新项目 |
| `DELETE /projects/:id` | P1 | 删除项目 |

#### 监控 API
| 端点 | 优先级 | 说明 |
|------|--------|------|
| `GET /metrics/usage` | P1 | 使用统计 |
| `GET /metrics/overview` | P1 | 概览指标 |

#### Skill & Plugin API
| 端点 | 优先级 | 说明 |
|------|--------|------|
| `GET /skills` | P1 | Skill 列表 |
| `POST /skills` | P1 | 上传 Skill |
| `GET /skills/:id` | P1 | Skill 详情 |
| `DELETE /skills/:id` | P1 | 删除 Skill |
| `POST /skills/:id/publish` | P1 | 发布 Skill |
| `GET /plugins` | P1 | Plugin 列表 |
| `POST /plugins` | P1 | 安装 Plugin |
| `GET /plugins/:id` | P1 | Plugin 详情 |
| `DELETE /plugins/:id` | P1 | 卸载 Plugin |
| `POST /plugins/:id/enable` | P1 | 启用 Plugin |
| `POST /plugins/:id/disable` | P1 | 禁用 Plugin |

#### 文件管理 API
| 端点 | 优先级 | 说明 |
|------|--------|------|
| `GET /instances/:id/files` | P1 | 文件列表 |
| `POST /instances/:id/files/upload` | P1 | 上传文件 |
| `GET /instances/:id/files/:path` | P1 | 下载文件 |
| `PUT /instances/:id/files/:path` | P1 | 更新文件 |
| `DELETE /instances/:id/files/:path` | P1 | 删除文件 |

---

## 四、开发阶段路线图

### Phase 1: MVP 完善 (预计 2-3 周)

**目标：** 完成架构中 Phase 1 的核心功能

#### Sprint 1.1: 认证与授权
- [ ] JWT 认证实现
- [ ] 用户登录/登出 API
- [ ] 中间件：租户上下文注入
- [ ] 前端登录页面

#### Sprint 1.2: 配置模板管理
- [ ] ConfigTemplate Repository 实现
- [ ] 配置模板 CRUD API
- [ ] 前端配置模板管理页面
- [ ] 配置模板关联到实例创建

#### Sprint 1.3: 租户与项目管理
- [ ] Tenant Repository 实现
- [ ] Project Repository 实现
- [ ] 租户/项目 CRUD API
- [ ] 前端租户/项目管理页面
- [ ] 租户隔离逻辑实现

#### Sprint 1.4: 实例监控与日志
- [ ] 监控指标采集 (CPU/Memory)
- [ ] 日志流 WebSocket
- [ ] 前端日志查看页面
- [ ] 前端监控仪表板

---

### Phase 2: 配置与适配完善 (预计 2-3 周)

**目标：** 完善配置管理和多类型适配

#### Sprint 2.1: 配置版本控制
- [ ] 配置版本历史表设计
- [ ] 版本创建逻辑
- [ ] 版本对比与回滚
- [ ] 前端版本历史页面

#### Sprint 2.2: NanoClaw Adapter
- [ ] NanoClaw Adapter 实现
- [ ] NanoClaw 配置格式定义
- [ ] Adapter 工厂注册

#### Sprint 2.3: 批量操作
- [ ] 批量配置下发
- [ ] 批量实例操作
- [ ] 批量升级机制

---

### Phase 3: 企业能力 (预计 3-4 周)

**目标：** 实现企业级治理能力

#### Sprint 3.1: 配额与策略
- [ ] 配额检查中间件
- [ ] 资源配额 API
- [ ] 策略引擎框架
- [ ] 配额超限处理

#### Sprint 3.2: 文件统一管理
- [ ] PVC 创建与管理
- [ ] 对象存储集成 (MinIO)
- [ ] 文件上传/下载 API
- [ ] 前端文件管理页面

#### Sprint 3.3: Skill 市场
- [ ] Skill Registry 服务
- [ ] Skill 上传与分发
- [ ] 实例 Skill 拉取机制
- [ ] 前端 Skill 市场页面

#### Sprint 3.4: 插件系统
- [ ] Plugin 加载机制
- [ ] Plugin 生命周期管理
- [ ] Plugin API 端点
- [ ] 前端插件管理页面

---

### Phase 4: 运维与扩展 (预计 2-3 周)

**目标：** 运维支持和扩展能力

#### Sprint 4.1: RBAC 与审计
- [ ] 用户与角色管理
- [ ] RBAC 权限中间件
- [ ] 操作审计日志
- [ ] 前端权限管理页面

#### Sprint 4.2: 多运行时支持
- [ ] Docker 运行时适配器
- [ ] 运行时抽象层
- [ ] 运行时切换配置

#### Sprint 4.3: 高可用与扩展
- [ ] Leader Election
- [ ] 任务队列 (NATS)
- [ ] 多实例部署支持
- [ ] 状态同步机制

#### Sprint 4.4: 可观测性增强
- [ ] Prometheus 指标导出
- [ ] Grafana 仪表板
- [ ] 日志聚合 (Loki)
- [ ] 告警规则配置

---

### Phase 5: 移动端与高级特性 (预计 2-3 周)

**目标：** 移动端支持和高级功能

#### Sprint 5.1: 移动端 API
- [ ] Mobile API Gateway
- [ ] 流式输出优化
- [ ] 离线同步机制

#### Sprint 5.2: 自动化运维
- [ ] 自动扩缩容
- [ ] 自愈机制
- [ ] 备份与恢复

#### Sprint 5.3: 多集群支持
- [ ] 多集群配置
- [ ] 跨集群调度
- [ ] 集群状态聚合

---

## 五、技术债务与改进项

### 5.1 代码质量

| 问题 | 优先级 | 说明 |
|------|--------|------|
| TODO 注释清理 | P2 | 代码中有 TODO 标记未完成 |
| 测试覆盖 | P1 | 缺少单元测试和集成测试 |
| 文档完善 | P2 | API 文档缺失 |
| 错误处理 | P1 | 部分错误被忽略 |

### 5.2 架构改进

| 问题 | 优先级 | 说明 |
|------|--------|------|
| 状态同步机制 | P1 | 目前为异步 goroutine，缺少错误处理 |
| 配置下发优化 | P2 | 当前每次都更新 ConfigMap，可优化 |
| 事件驱动 | P2 | 缺少事件总线 |
| 重试机制 | P1 | K8S 操作缺少重试 |

### 5.3 性能优化

| 问题 | 优先级 | 说明 |
|------|--------|------|
| 数据库连接池 | P1 | SQLite 连接管理优化 |
| API 响应缓存 | P2 | 静态数据缓存 |
| 分页查询优化 | P1 | 目前未实现总数查询 |
| 资源限制 | P2 | 防止资源耗尽 |

---

## 六、依赖与集成点

### 6.1 外部依赖

| 组件 | 用途 | 集成状态 |
|------|------|---------|
| Kubernetes | 容器编排 | ✅ 已集成 (kom 库) |
| SQLite | 元数据存储 | ✅ 已使用 |
| MinIO/S3 | 对象存储 | ❌ 未集成 |
| Prometheus | 指标采集 | ❌ 未集成 |
| Grafana | 可视化 | ❌ 未集成 |
| NATS/Kafka | 事件队列 | ❌ 未集成 |

### 6.2 Claw 依赖

| 组件 | 说明 | 状态 |
|------|------|------|
| OpenClaw | 目标类型 1 | ✅ Adapter 已实现 |
| NanoClaw | 目标类型 2 | ❌ Adapter 未实现 |
| 自定义 Claw | 扩展支持 | ✅ 接口已定义 |

---

## 七、风险评估

| 风险 | 级别 | 缓解措施 |
|------|--------|----------|
| K8S 操作失败影响状态一致性 | 中 | 增加重试机制 + 状态修复任务 |
| 多租户隔离不彻底 | 高 | 加强中间件 + 数据库查询隔离 |
| 配置下发失败导致实例启动异常 | 中 | 配置验证 + 回滚机制 |
| 高并发下资源竞争 | 中 | 使用分布式锁 + 事件队列 |
| 数据库迁移风险 | 中 | 版本化迁移脚本 + 回滚方案 |

---

## 八、附录

### 8.1 文件索引

| 模块 | 关键文件 |
|------|----------|
| 主程序 | `cmd/controlplane/main.go` |
| API 层 | `internal/api/`, `internal/router.go` |
| 服务层 | `internal/service/` |
| 数据访问层 | `internal/repository/`, `internal/model/` |
| 适配器 | `internal/adapter/` |
| K8S 运行时 | `internal/runtime/k8s/` |
| 配置 | `config/` |
| 前端 | `frontend/`, `internal/embed/frontend.go` |

### 8.2 数据库表详情

**tenants 表**
```sql
id, name, max_instances, max_cpu, max_memory, max_storage,
created_at, updated_at
```

**projects 表**
```sql
id, tenant_id (FK), name, created_at, updated_at
```

**config_templates 表**
```sql
id, name, description, variables (BLOB), adapter_type, version,
created_at, updated_at
```

**claw_instances 表**
```sql
id, name, tenant_id (FK), project_id (FK), type, version,
status, config (BLOB), cpu, memory, config_dir, data_dir,
storage_size, created_at, updated_at
```

---

**文档版本:** v1.0
**最后更新:** 2026-02-26
**维护者:** Open Cluster Claw 开发团队

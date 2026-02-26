# Open Cluster Claw 总体架构设计

## 一、产品定位

**Open Cluster Claw = 多实例、多类型 Claw 统一控制平面（Control Plane）**

本质上类似于 Kubernetes 的控制平面，但面向的是 "Claw 实例" 而不是容器。

### 核心目标

| 目标 | 说明 |
|------|------|
| 规模化管理 | 统一管理多个 Claw 实例 |
| 生命周期管理 | 创建、启动、停止、销毁等全生命周期控制 |
| 集中化配置 | 统一配置模板、版本控制、批量下发 |
| 多类型适配 | 支持 OpenClaw、NanoClaw 等多种类型 |
| 数据与文件统一 | 用户数据与实例生命周期解耦 |
| 可观测性 | 实例监控、使用统计、日志聚合 |

### 核心设计原则

1. **控制面 / 数据面分离**
2. **多租户隔离优先**
3. **可水平扩展**
4. **声明式驱动（Declarative Model）**
5. **强可观测性**
6. **可演进为分布式架构**

---

## 二、整体架构分层

```
┌─────────────────────────────────────────────────────────────┐
│                       User Portal                           │
│  (Web Console / CLI / Mobile App / Open API)                │
└────────────────────────────┬────────────────────────────────┘
                             │
┌────────────────────────────▼────────────────────────────────┐
│                      API Gateway Layer                      │
│  Auth / RBAC / Rate Limit / Multi-Tenant Context           │
│  WebSocket Support / Mobile Support / Audit Logging        │
└────────────────────────────┬────────────────────────────────┘
                             │
        ┌────────────────────┼────────────────────┐
        │                    │                    │
┌───────▼───────┐   ┌────────▼────────┐   ┌──────▼──────┐
│ Control Plane │   │ Config Center   │   │  Skill &    │
│ (Core Engine) │   │                 │   │  Plugin     │
│               │   │                 │   │  Manager    │
└───────┬───────┘   └────────┬────────┘   └──────┬──────┘
        │                    │                    │
        └────────────────────┼────────────────────┘
                             │
┌────────────────────────────▼────────────────────────────────┐
│              Resource Orchestrator Layer                    │
│  (Cluster Adapter / Infra Integration / Scheduler)          │
└────────────────────────────┬────────────────────────────────┘
                             │
┌────────────────────────────▼────────────────────────────────┐
│                    Storage Layer                            │
│  PostgreSQL / Object Storage / PVC / Vector DB (Future)    │
└────────────────────────────┬────────────────────────────────┘
                             │
┌────────────────────────────▼────────────────────────────────┐
│                  Runtime Layer (K8S)                        │
│  Claw Instances Pool (Distributed Runtime / Nodes / Pods)  │
└─────────────────────────────────────────────────────────────┘
```

---

## 三、核心组件详解

### 1. API Gateway Layer

**职责：**

| 功能 | 说明 |
|------|------|
| 统一入口 | 所有请求的统一接入点 |
| 认证授权 | OAuth2 / JWT / SSO |
| 租户识别 | 上下文注入与租户隔离 |
| RBAC 权限 | 基于角色的访问控制 |
| 限流保护 | API 限流与熔断 |
| 审计日志 | 操作审计与追溯 |
| WebSocket | 流式输出支持 |
| 移动端接入 | Mobile App API 支持 |

**移动端访问模型：**

```
Mobile App
      ↓
  Gateway (Token 验证、会话隔离、流式转发)
      ↓
Claw Instance
```

---

### 2. Control Plane Core

**核心调度与治理引擎，包含以下子模块：**

#### 2.1 Instance Manager（实例管理器）

**职责：**
- 实例生命周期管理（创建、启动、停止、重启、销毁、强制终止）
- 实例调度与扩缩容
- 状态同步（期望状态 vs 实际状态）

**生命周期状态机：**

```
      Creating
          ↓
       Running ←→ Stopped
          ↓
       Failed
          ↓
      Destroyed
```

**支持的操作：**

| 操作 | API | CLI |
|------|-----|-----|
| 创建实例 | POST /instances | claw create |
| 启动 | POST /instances/:id/start | claw start |
| 停止 | POST /instances/:id/stop | claw stop |
| 重启 | POST /instances/:id/restart | claw restart |
| 销毁 | DELETE /instances/:id | claw destroy |
| 强制终止 | POST /instances/:id/kill | claw kill --force |
| Console 登录 | POST /instances/:id/console | claw console |
| 本地 CLI 连接 | - | claw connect |

#### 2.2 Config Manager（配置管理器）

**职责：**
- 配置模板管理
- 配置版本控制
- 批量配置下发
- 灰度发布与回滚
- 配置一致性校验

**配置分类：**

| 级别 | 说明 | 示例 |
|------|------|------|
| 全局级 | 公共变量 | 大模型 API Key、飞书 Bot Key、代理配置 |
| 租户级 | 租户覆盖配置 | 企业内部 API、租户特定策略 |
| 实例级 | 实例特定配置 | 模型参数、Memory 配置、Skill 开关 |

**配置流程：**

```
Unified Config (JSON)
        ↓
    Adapter Engine
        ↓
Target Claw Config File
        ↓
    Distribute to Instance
```

#### 2.3 Adapter Engine（适配器引擎）

**职责：**
- 多类型 Claw 适配
- 配置格式转换
- 运行时注入

**ClawAdapter 接口：**

```go
type ClawAdapter interface {
    // 解析统一配置
    ParseConfig(unifiedConfig UnifiedConfig) error

    // 生成目标配置文件
    GenerateConfig() (targetConfig string, err error)

    // 配置验证
    Validate() error

    // 运行时注入
    InjectRuntime(runtimeInfo RuntimeInfo) error
}
```

**支持的 Claw 类型：**

| 类型 | Adapter |
|------|---------|
| OpenClaw | OpenClawAdapter |
| NanoClaw | NanoClawAdapter |
| 自定义 | 自定义 Adapter |

#### 2.4 Usage Monitor（使用监控器）

**职责：**
- 实例级监控（CPU、Memory、重启次数）
- 使用级监控（Token 使用量、调用次数、错误率、活跃时间）
- 指标聚合与存储

#### 2.5 Policy Engine（策略引擎）

**职责：**
- 租户资源配额控制
- 访问策略执行
- 自动扩缩容策略

---

### 3. Skill & Plugin Manager

**职责：**

| 功能 | Skill | Plugin |
|------|-------|--------|
| 含义 | AI 能力扩展 | 系统级扩展能力 |
| 版本管理 | ✓ | ✓ |
| 批量升级 | ✓ | ✓ |
| 启用/禁用 | ✓ | ✓ |

**内网 Skill 市场架构：**

```
┌─────────────────┐
│ Skill Registry  │
│  (集中存储)     │
└────────┬────────┘
         │
┌────────▼────────┐
│ Skill           │
│ Distributor     │
│ (统一推送)      │
└────────┬────────┘
         │
┌────────▼────────┐
│ Claw Instance   │
│ (主动拉取)      │
└─────────────────┘
```

**支持模式：**
1. 内网 Skill Market（实例主动拉取）
2. 统一推送（管理平台直接下发）

---

### 4. Resource Orchestrator Layer

**职责：**
- 与底层基础设施对接
- 实例调度与资源分配
- 自动扩缩容

**支持的底层运行时：**

| 运行时 | 状态 |
|--------|------|
| Kubernetes | ✓ |
| Docker | 计划中 |
| Podman | 计划中 |
| 物理机 | 计划中 |
| 虚拟机 | 计划中 |
| 云厂商 API | 计划中 |

---

### 5. Storage Layer

**职责：**
- 元数据存储（PostgreSQL）
- 文件存储（对象存储 / PVC）
- 日志存储
- 向量数据库（未来）

**文件系统模型：**

每个实例挂载两个目录：

| 目录 | 说明 | 实现方式 |
|------|------|----------|
| ConfigDir | 配置目录（如 /root/.config） | PVC / 对象存储 |
| DataDir | 用户数据目录 | PVC / 对象存储 |

**文件管理功能：**
- 文件上传/下载
- 在线编辑
- 文件共享（多实例共享同一目录）
- 云盘模式

**核心原则：** 用户数据与实例生命周期解耦

---

### 6. Runtime Layer (K8S)

**职责：**
- Claw 实例运行
- 资源隔离
- 健康检查

**实现方式：**
- 每个 Claw 实例 = 一个独立 Pod
- 使用 CRD 进行声明式管理
- 支持多集群部署

---

## 四、多租户模型

### 模型结构

```
Tenant
  ├── Projects
  │      ├── Claw Instances
  │      ├── Quota
  │      └── Config Override
  └── Policies
```

### 隔离方式

| 类型 | 说明 |
|------|------|
| 逻辑隔离 | Namespace 级隔离 |
| 资源配额 | CPU、Memory、实例数量限制 |
| 网络隔离 | 可选的网络策略隔离 |

---

## 五、可观测性体系

### 监控维度

#### 实例级监控
- Running 状态
- CPU 使用率
- Memory 使用率
- 重启次数

#### 使用级监控
- Token 使用量
- 调用次数
- 错误率
- 活跃时间

### 监控架构

```
┌─────────────────────────────────────────────────────────────┐
│                    Metrics Collection                       │
│  (Prometheus / Custom Collector)                            │
└────────────────────────────┬────────────────────────────────┘
                             │
┌────────────────────────────▼────────────────────────────────┐
│                    Metrics Storage                          │
│  (Time Series DB)                                          │
└────────────────────────────┬────────────────────────────────┘
                             │
┌────────────────────────────▼────────────────────────────────┐
│                    Visualization                           │
│  (Grafana / Web UI Dashboard)                              │
└─────────────────────────────────────────────────────────────┘
```

---

## 六、数据模型设计

### ClawInstance 资源

```yaml
apiVersion: clusterclaw.io/v1
kind: ClawInstance
metadata:
  name: claw-instance-01
  namespace: tenant-a
  labels:
    app: claw
    type: openclaw
spec:
  # Claw 类型
  type: OpenClaw
  # 版本
  version: "1.2.0"
  # 副本数
  replicas: 1
  # 资源限制
  resources:
    cpu: "2"
    memory: "4Gi"
  # 配置引用
  config:
    template: default-config
    overrides:
      - key: model.name
        value: claude-3-opus
      - key: memory.limit
        value: 200000
  # 存储配置
  storage:
    configDir: /root/.config
    dataDir: /data
    size: 10Gi
  # 网络配置
  network:
    expose: true
    port: 8080
status:
  # 生命周期阶段
  phase: Running
  # 就绪副本数
  readyReplicas: 1
  # 最后状态变更时间
  lastTransitionTime: "2026-02-25T10:00:00Z"
  # 实例状态
  conditions:
    - type: Ready
      status: "True"
      reason: PodReady
      message: Instance is ready
```

### ConfigTemplate 资源

```yaml
apiVersion: clusterclaw.io/v1
kind: ConfigTemplate
metadata:
  name: default-config
  namespace: default
spec:
  # 配置描述
  description: Default OpenClaw configuration
  # 变量定义
  variables:
    - name: model.name
      type: string
      default: "claude-3-haiku"
      required: true
      description: Model name to use
    - name: model.apiKey
      type: string
      required: true
      secret: true
      description: Model API key
    - name: memory.limit
      type: number
      default: 100000
      required: false
      description: Memory token limit
  # 适配器类型
  adapter: OpenClawAdapter
  # 版本
  version: "1.0.0"
```

### Tenant 资源

```yaml
apiVersion: clusterclaw.io/v1
kind: Tenant
metadata:
  name: tenant-a
spec:
  # 租户配额
  quota:
    instances: 10
    cpu: "20"
    memory: "40Gi"
    storage: "100Gi"
  # 租户策略
  policies:
    - name: max-replicas
      value: "3"
  # 租户管理员
    admins:
      - user1@example.com
```

---

## 七、API 设计

### RESTful API

#### 实例管理 API

```
GET    /api/v1/instances                    # 列表（支持分页、过滤、排序）
POST   /api/v1/instances                    # 创建
GET    /api/v1/instances/:id                # 详情
PUT    /api/v1/instances/:id                # 更新
DELETE /api/v1/instances/:id                # 删除
POST   /api/v1/instances/:id/start          # 启动
POST   /api/v1/instances/:id/stop           # 停止
POST   /api/v1/instances/:id/restart        # 重启
POST   /api/v1/instances/:id/kill           # 强制终止
POST   /api/v1/instances/:id/console        # Console 连接
GET    /api/v1/instances/:id/logs           # 日志流
```

#### 配置管理 API

```
GET    /api/v1/configs                      # 列表
POST   /api/v1/configs                      # 创建
GET    /api/v1/configs/:id                  # 详情
PUT    /api/v1/configs/:id                  # 更新
DELETE /api/v1/configs/:id                  # 删除
POST   /api/v1/configs/:id/publish          # 发布
POST   /api/v1/configs/:id/rollback         # 回滚
GET    /api/v1/configs/:id/versions         # 版本历史
POST   /api/v1/configs/:id/validate         # 验证配置
```

#### 租户管理 API

```
GET    /api/v1/tenants                      # 列表
POST   /api/v1/tenants                      # 创建
GET    /api/v1/tenants/:id                  # 详情
PUT    /api/v1/tenants/:id                  # 更新
DELETE /api/v1/tenants/:id                  # 删除
GET    /api/v1/tenants/:id/quota            # 查询配额
PUT    /api/v1/tenants/:id/quota            # 更新配额
```

#### 监控 API

```
GET    /api/v1/instances/:id/metrics        # 实例指标
GET    /api/v1/instances/:id/logs           # 实例日志
GET    /api/v1/instances/:id/events         # 实例事件
GET    /api/v1/metrics/usage                # 使用统计
GET    /api/v1/metrics/overview             # 概览指标
```

#### Skill & Plugin API

```
GET    /api/v1/skills                       # Skill 列表
POST   /api/v1/skills                       # 上传 Skill
GET    /api/v1/skills/:id                   # Skill 详情
DELETE /api/v1/skills/:id                   # 删除 Skill
POST   /api/v1/skills/:id/publish           # 发布 Skill
GET    /api/v1/plugins                      # Plugin 列表
POST   /api/v1/plugins                      | 安装 Plugin
GET    /api/v1/plugins/:id                  | Plugin 详情
DELETE /api/v1/plugins/:id                  | 卸载 Plugin
POST   /api/v1/plugins/:id/enable           | 启用 Plugin
POST   /api/v1/plugins/:id/disable          | 禁用 Plugin
```

#### 文件管理 API

```
GET    /api/v1/instances/:id/files          | 文件列表
POST   /api/v1/instances/:id/files/upload   | 上传文件
GET    /api/v1/instances/:id/files/:path    | 下载文件
PUT    /api/v1/instances/:id/files/:path    | 更新文件
DELETE /api/v1/instances/:id/files/:path    | 删除文件
```

### WebSocket API

#### 流式输出

```
WS /api/v1/instances/:id/stream
```

#### 实时日志

```
WS /api/v1/instances/:id/logs/stream
```

#### 实时指标

```
WS /api/v1/instances/:id/metrics/stream
```

---

## 八、关键技术架构模式

### 1. 声明式资源模型

采用 CRD 风格的声明式资源模型，控制面负责：

1. 监听资源变化
2. 执行 Reconcile Loop
3. 确保运行状态与 Spec 一致

**Reconcile Loop 流程：**

```
Desired State (DB)
        ↓
Reconcile Engine
        ↓
Infra Adapter
        ↓
Actual State
        ↓
Status Update (DB)
```

### 2. 事件驱动架构

```
┌─────────────────────────────────────────────────────────────┐
│                    Event Bus (NATS/Kafka)                   │
└────────────────────────────┬────────────────────────────────┘
                             │
        ┌────────────────────┼────────────────────┐
        │                    │                    │
┌───────▼───────┐   ┌────────▼────────┐   ┌──────▼──────┐
│ Instance     │   │ Config          │   │ Usage       │
│ Events       │   │ Events          │   │ Events      │
└───────────────┘   └─────────────────┘   └─────────────┘
```

### 3. 横向扩展设计

控制面支持多实例部署：

| 能力 | 实现方式 |
|------|----------|
| Leader Election | etcd / 分布式锁 |
| 任务队列 | NATS / Kafka |
| 状态同步 | 共享数据库 |

---

## 九、非功能性需求

### 性能指标

| 指标 | 目标值 |
|------|--------|
| 并发实例数 | ≥ 100 |
| API 响应时间（P95） | < 200ms |
| 配置下发延迟 | < 5s |
| 实例启动时间 | < 30s |

### 可用性

- 控制面高可用（多副本部署）
- 数据持久化（不依赖控制面状态）
- 自动故障转移

### 安全性

| 能力 | 说明 |
|------|------|
| 租户隔离 | 逻辑隔离 + 资源配额 |
| RBAC | 基于角色的访问控制 |
| 审计日志 | 操作审计与追溯 |
| 敏感配置加密 | 敏感信息加密存储 |

### 可扩展性

- 水平扩展控制面
- 支持多集群部署
- 支持多种底层运行时

---

## 十、推荐技术栈

| 组件 | 技术选型 |
|------|----------|
| 控制面 | Golang |
| API Gateway | Golang / Nginx |
| 存储 | PostgreSQL |
| 对象存储 | MinIO / S3 |
| 事件系统 | NATS / Kafka |
| 任务调度 | 内置 worker + queue |
| Web UI | React / Vue |
| 运行时 | Kubernetes |
| 监控 | Prometheus + Grafana |
| 日志 | Loki / ELK |

---

## 十一、部署架构

### 小规模企业

```
┌─────────────────────────────────────────────────────────────┐
│                      Single Node                            │
│  ┌─────────────┐  ┌─────────────┐  ┌─────────────┐         │
│  │ Control     │  │ PostgreSQL  │  │ K8S Cluster │         │
│  │ Plane       │  │             │  │             │         │
│  └─────────────┘  └─────────────┘  └─────────────┘         │
└─────────────────────────────────────────────────────────────┘
```

### 中大型企业

```
┌─────────────────────────────────────────────────────────────┐
│                      HA Deployment                          │
│  ┌─────────────┐  ┌─────────────┐  ┌─────────────┐         │
│  │ Control     │  │ PostgreSQL  │  │ K8S Cluster │         │
│  │ Plane (x3)  │  │   (HA)      │  │ (Multi)     │         │
│  └─────────────┘  └─────────────┘  └─────────────┘         │
└─────────────────────────────────────────────────────────────┘
```

---

## 十二、开发阶段路线图

### Phase 1: MVP（核心控制面）

**目标：** 基本的多实例生命周期管理

- [x] 架构设计
- [ ] 实例 CRUD API
- [ ] 基础配置模板
- [ ] K8S Pod 部署集成
- [ ] 实例状态监控
- [ ] 简单 Web UI

### Phase 2: 配置与适配

**目标：** 完善配置管理和多类型适配

- [ ] 配置版本控制
- [ ] 批量配置下发
- [ ] Adapter 机制
- [ ] OpenClaw Adapter
- [ ] NanoClaw Adapter

### Phase 3: 企业能力

**目标：** 企业级治理能力

- [ ] 多租户隔离
- [ ] 文件统一管理
- [ ] Skill 市场（内网）
- [ ] 插件系统
- [ ] RBAC 权限

### Phase 4: 运维与扩展

**目标：** 运维支持和扩展能力

- [ ] 移动端 API Gateway
- [ ] 应用部署能力
- [ ] 多集群支持
- [ ] 自动扩缩容

---

## 十三、总结

openClusterClaw 本质上是：

> 一个专为 OpenClaw 设计的企业级控制平面（Control Plane）

三大核心：

1. **生命周期统一管理** - 统一控制多个实例的完整生命周期
2. **配置集中化** - 通过 Adapter 机制实现多类型配置统一管理
3. **多类型适配** - 支持不同类型 Claw 的统一接入

它将：

- 分散实例 → 统一纳管
- 手工部署 → 声明式管理
- 工具运行 → 基础设施服务

最终目标：

**让 OpenClaw 成为企业内部的标准化服务基础设施层。**
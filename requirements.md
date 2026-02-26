# Open Cluster Claw 功能需求规范

## 产品定位

**Open Cluster Claw = 多实例、多类型 Claw 统一控制平面（Control Plane）**

核心目标：

- 规模化管理多个 Claw 实例
- 统一生命周期管理
- 集中化配置管理
- 多类型 Claw 适配
- 数据与文件统一
- 可观测性与运维控制

本质上类似于 Kubernetes 的控制平面，但面向的是 "Claw 实例" 而不是容器。

---

## 一、平台级能力（基础架构）

### 1. API Gateway 层

统一入口层，负责：

- 认证与授权（OAuth2 / JWT / SSO）
- 租户识别与上下文注入
- RBAC 权限控制
- 限流与审计
- 移动端接入支持

### 2. 多租户模型

```
Tenant
  ├── Projects
  │      ├── OpenClaw Instances
  │      ├── Quota
  │      └── Config Override
  └── Policies
```

隔离方式：

- 逻辑隔离（Namespace 级）
- 资源配额限制
- 网络隔离（可选）

### 3. 可观测性体系

统一监控 OpenClaw 集群状态：

- 实例健康状态
- CPU / Memory / IO
- 请求吞吐量
- 错误率
- 日志聚合
- 审计日志

---

## 二、核心能力（必须优先实现）

### 1️⃣ 多实例统一管理（Control Plane）

#### 1.1 实例生命周期管理

支持操作：

| 操作 | 说明 |
|------|------|
| 创建实例 | 初始化新实例 |
| 启动 | 启动已停止实例 |
| 停止 | 优雅停止实例 |
| 重启 | 重启实例 |
| 销毁 | 销毁实例及资源 |
| 强制终止 | kill 强制终止 |
| Console 登录 | 命令行连接 |
| Web Console | Web 终端访问 |
| 本地 CLI 连接 | 本地 CLI 连接实例 |

#### 1.2 生命周期状态

```
Creating → Running → Stopped → Destroyed
                   ↓
                 Failed
```

#### 1.3 实现方式

- 底层运行基于 Kubernetes Pod
- 每个 Claw 实例 = 一个独立 Pod
- 使用 CRD（未来可扩展）

---

### 2️⃣ 统一配置管理（Configuration Center）

#### 2.1 配置能力

支持：

- 可视化配置界面
- 生成标准配置文件
- 配置模板管理
- 配置版本控制
- 配置与实例绑定
- 批量下发
- 批量修改
- 灰度发布
- 配置回滚
- 一致性校验

#### 2.2 配置分类

##### 公共变量（全局级）

- 大模型 API Key
- 飞书 Bot Key
- 企业内部 API
- 代理配置

##### 实例级配置

- 模型参数
- Memory 配置
- Skill 开关
- 权限控制

##### 类型适配配置

通过 Adapter 生成对应配置文件：

```
Unified Config → Adapter → Target Claw Config File
```

---

### 3️⃣ 多类型 Claw 适配（Adapter 架构）

支持 Claw 类型：

- OpenClaw
- NanoClaw
- 其他 Claw 类型（可扩展）

#### ClawAdapter 接口定义

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

统一配置中心输出标准 JSON，各 Adapter 负责转换为对应格式。

---

### 4️⃣ 数据与文件统一管理

#### 4.1 文件系统模型

每个实例挂载两个目录：

| 类型 | 说明 |
|------|------|
| ConfigDir | 配置目录（如 /root） |
| DataDir | 用户数据目录 |

底层实现：

- 持久化卷（PVC）
- 或对象存储挂载

#### 4.2 支持功能

- 文件上传
- 文件下载
- 在线编辑
- 文件共享（多实例共享同一目录）
- 云盘模式

**核心原则：** 用户数据与实例生命周期解耦

---

### 5️⃣ Skill 市场（内网可用）

企业环境无法访问公网时必须支持两种模式：

#### 模式一：内网 Skill Market

- 集中存储
- 实例主动拉取

#### 模式二：统一推送

- 管理平台直接下发到实例

架构设计：

```
Skill Registry
     ↓
Skill Distributor
     ↓
Claw Instance
```

---

### 6️⃣ 插件统一管理

插件与 Skill 分离：

| 类型 | 含义 |
|------|------|
| Skill | AI 能力扩展 |
| Plugin | 系统级扩展能力 |

支持功能：

- 插件版本管理
- 批量升级
- 启用/禁用

---

### 7️⃣ 运行状态监控（Observability）

#### 7.1 实例级监控

- Running 状态
- CPU 使用率
- Memory 使用率
- 重启次数

#### 7.2 使用级监控

- Token 使用量
- 调用次数
- 错误率
- 活跃时间

#### 7.3 命令式管理

```bash
claw status
claw stop --force
claw usage
```

---

### 8️⃣ 统一数据模型

目标：所有 Claw 使用同一个用户数据空间

实现方式：

- 统一账号体系
- 统一存储
- 统一向量数据库（未来）

---

### 9️⃣ 移动端支持（API Gateway）

手机 App 访问模型：

```
Mobile App
      ↓
  Gateway
      ↓
Claw Instance
```

能力：

- 转发流式输出
- WebSocket 支持
- Token 验证
- 会话隔离

---

## 三、扩展能力（第二阶段）

### 10️⃣ 应用部署能力（Runtime + Deploy）

涉及：

- 容器构建
- 镜像生成
- 自动部署
- 子域名分配

底层可能涉及：

- Docker
- Podman

---

## 四、系统架构分层

```
┌────────────────────────────┐
│        Web Console         │
└────────────┬───────────────┘
             ↓
┌────────────────────────────┐
│      API Gateway Layer     │
└────────────┬───────────────┘
             ↓
┌────────────────────────────┐
│   Control Plane Core       │
│ - Instance Manager         │
│ - Config Manager           │
│ - Adapter Engine           │
│ - Skill Manager            │
│ - Plugin Manager           │
│ - Usage Monitor            │
└────────────┬───────────────┘
             ↓
┌────────────────────────────┐
│    Runtime Layer (K8S)     │
└────────────────────────────┘
```

---

## 五、功能优先级排序

### MVP（必做）

1. 多实例生命周期管理
2. 统一配置管理
3. Adapter 机制
4. 批量管理能力
5. 基础运行状态监控

### 第二阶段

6. 文件统一管理
7. Skill 市场
8. 插件系统

### 第三阶段

9. 移动端支持
10. 自动部署能力

---

## 六、核心本质总结

Open Cluster Claw 本质是：

> 一个面向 Claw 的企业级控制平面系统

三大核心：

1. 生命周期统一管理
2. 配置集中化
3. 多类型适配

其它功能都围绕这三点展开。

---

## 七、数据模型设计（建议）

### Instance 模型

```yaml
apiVersion: clusterclaw.io/v1
kind: ClawInstance
metadata:
  name: claw-instance-01
  namespace: tenant-a
spec:
  type: OpenClaw
  version: "1.2.0"
  replicas: 1
  resources:
    cpu: "2"
    memory: "4Gi"
  config:
    template: default-config
    overrides:
      - key: model.name
        value: claude-3-opus
  storage:
    configDir: /root/.config
    dataDir: /data
status:
  phase: Running
  readyReplicas: 1
  lastTransitionTime: "2026-02-25T10:00:00Z"
```

### ConfigTemplate 模型

```yaml
apiVersion: clusterclaw.io/v1
kind: ConfigTemplate
metadata:
  name: default-config
spec:
  variables:
    - name: model.name
      type: string
      default: "claude-3-haiku"
      required: true
    - name: memory.limit
      type: number
      default: 100000
      required: false
```

---

## 八、API 设计建议

### 实例管理 API

```
GET    /api/v1/instances                    # 列表
POST   /api/v1/instances                    # 创建
GET    /api/v1/instances/:id                # 详情
PUT    /api/v1/instances/:id                # 更新
DELETE /api/v1/instances/:id                # 删除
POST   /api/v1/instances/:id/start          # 启动
POST   /api/v1/instances/:id/stop           # 停止
POST   /api/v1/instances/:id/restart        # 重启
POST   /api/v1/instances/:id/console        # Console
```

### 配置管理 API

```
GET    /api/v1/configs                      # 列表
POST   /api/v1/configs                      # 创建
GET    /api/v1/configs/:id                  # 详情
PUT    /api/v1/configs/:id                  # 更新
DELETE /api/v1/configs/:id                  # 删除
POST   /api/v1/configs/:id/publish          # 发布
POST   /api/v1/configs/:id/rollback         # 回滚
```

### 监控 API

```
GET    /api/v1/instances/:id/metrics        # 实例指标
GET    /api/v1/instances/:id/logs           # 实例日志
GET    /api/v1/instances/:id/events         # 实例事件
GET    /api/v1/metrics/usage                # 使用统计
```

---

## 九、非功能性需求

### 性能

- 支持至少 100+ 实例并发管理
- API 响应时间 < 200ms（P95）
- 配置下发延迟 < 5s

### 可用性

- 控制面高可用（多副本部署）
- 数据持久化（不依赖控制面状态）

### 安全性

- 租户间数据隔离
- RBAC 权限控制
- 操作审计日志
- 敏感配置加密存储

### 可扩展性

- 水平扩展控制面
- 支持多集群部署
- 支持多种底层运行时

---

## 十、开发阶段路线图

### Phase 1: MVP（核心控制面）

**目标：** 基本的多实例生命周期管理

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
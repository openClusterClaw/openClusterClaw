# Open Cluster Claw 开发计划

> 变更记录表

| 版本 | 日期 | 作者 | 变更说明 |
|------|------|------|----------|
| v1.3 | 2026-02-27 | AI | 全面更新实现状态，基于代码实际审查更新进度和完成项 |
| v1.2 | 2026-02-27 | AI | 更新实现状态，增加认证授权、OTP功能完成项 |
| v1.1 | 2026-02-27 | AI | 更新实现状态，增加认证授权完成项 |
| v1.0 | 2026-02-26 | AI | 初始版本 |

> 基于架构设计文档 (`architecture.md`) 与当前代码实现状态的对比分析
>
> 更新日期: 2026-02-27

---

## 一、实现状态概览

### 整体进度: 约 35% 完成

| 架构层级 | 完成度 | 说明 |
|---------|--------|------|
| API Gateway Layer | 65% | JWT认证、OTP完成，RBAC基础完成，缺少限流、审计、WebSocket |
| Control Plane Core | 40% | 实例管理+认证完成，配置管理(模板CRUD)未实现 |
| Config Manager | 20% | 适配器完成，模板管理Model定义但未实现Service/Repository |
| Resource Orchestrator | 45% | K8S Pod/ConfigMap完成，其他运行时未实现 |
| Storage Layer | 25% | SQLite完成，GORM AutoMigrate完成，对象存储未实现 |
| Runtime Layer (K8S) | 55% | Pod/ConfigMap完成，PVC未实现 |
| Skill & Plugin Manager | 0% | 未实现 |
| User Portal (Frontend) | 40% | 基础UI+登录+OTP设置完成 |
| 多租户模型 | 20% | 表结构+租户上下文+默认租户完成，隔离逻辑未完成 |
---

## 二、已实现功能清单

### 2.1 控制面核心 (Control Plane Core)

#### InstanceManager (实例管理器)
| 功能 | 状态 | 文件位置 |
|------|------|---------|
| 实例创建 (Create) | ✅ | `internal/service/instance.go:CreateInstance()` |
| 实例查询 (Get) | ✅ | `internal/service/instance.go:GetInstance()` |
| 实例列表 (List) | ✅ | `internal/service/instance.go:ListInstances()` |
| 实例更新 (Update) | ✅ | `internal/service/instance.go:UpdateInstance()` |
| 实例删除 (Delete) | ✅ | `internal/service/instance.go:DeleteInstance()` |
| 实例启动 (Start) | ✅ | `internal/service/instance.go:StartInstance()` |
| 实例停止 (Stop) | ✅ | `internal/service/instance.go:StopInstance()` |
| 实例重启 (Restart) | ✅ | `internal/service/instance.go:RestartInstance()` |
| 状态同步 (Sync) | ✅ | `internal/service/instance.go:syncPodStatus()` |
| ConfigMap 关联 | ✅ | `internal/service/instance.go:CreateInstance()` |
| Pod 等待就绪 | ✅ | `internal/runtime/k8s/pod.go:WaitForPodReady()` |

#### 数据库层 (Repository)
| 功能 | 状态 | 文件位置 |
|------|------|---------|
| InstanceRepository | ✅ | `internal/repository/instance.go` |
| UserRepository | ✅ | `internal/repository/user.go` |
| 数据库迁移 (Migrations) | ✅ | `cmd/controlplane/main.go:runMigrations()` |
| SQLite 支持 | ✅ | `cmd/controlplane/main.go` |
| GORM AutoMigrate | ✅ | `cmd/controlplane/main.go:runMigrations()` |
| 默认租户创建 | ✅ | `cmd/controlplane/main.go:initializeDefaultData()` |
| 默认管理员用户创建 | ✅ | `cmd/controlplane/main.go:initializeDefaultData()` |

#### API 层
| 功能 | 状态 | 文件位置 |
|------|------|---------|
| RESTful API Router | ✅ | `internal/api/router.go` |
| 实例管理 Handler | ✅ | `internal/api/handler.go` (Create/Get/List/Start/Stop/Restart/Delete) |
| 认证 Handler | ✅ | `internal/api/auth.go` (Login/Refresh/Logout/Me/CreateUser) |
| OTP Handler | ✅ | `internal/api/otp.go` (Generate/Enable/Disable/Verify/Status/Backup/LoginWithOTP) |
| 健康检查 | ✅ | `internal/api/router.go:/health` |

---

### 2.2 认证与授权 (Auth & Authorization)

#### JWT 认证
| 功能 | 状态 | 文件位置 |
|------|------|---------|
| JWT 服务 | ✅ | `internal/pkg/jwt/jwt.go` |
| Access Token 生成 | ✅ | `internal/pkg/jwt/jwt.go:GenerateAccessToken()` |
| Refresh Token 生成 | ✅ | `internal/pkg/jwt/jwt.go:GenerateRefreshToken()` |
| Token 验证 | ✅ | `internal/pkg/jwt/jwt.go:ValidateToken()` |
| Claims 模型 | ✅ | `internal/pkg/jwt/jwt.go:Claims` |

#### 认证服务
| 功能 | 状态 | 文件位置 |
|------|------|---------|
| 用户登录 | ✅ | `internal/service/auth.go:Login()` |
| Token 刷新 | ✅ | `internal/service/auth.go:RefreshToken()` |
| 用户创建 | ✅ | `internal/service/auth.go:CreateUser()` |
| 密码哈希 (bcrypt) | ✅ | `internal/service/auth.go:HashPassword()` |
| OTP 集成登录 | ✅ | `internal/service/otp.go:LoginWithOTP()` |
| 用户查询 | ✅ | `internal/service/auth.go:GetUserByID()` |

#### 中间件
| 功能 | 状态 | 文件位置 |
|------|------|---------|
| AuthMiddleware | ✅ | `internal/middleware/auth.go:AuthMiddleware()` |
| RequireAdmin | ✅ | `internal/middleware/auth.go:RequireAdmin()` |
| 租户上下文注入 | ✅ | `internal/middleware/auth.go` (Set ContextUserIDKey/UsernameKey/TenantIDKey/RoleKey) |
| 上下文获取辅助函数 | ✅ | `internal/middleware/auth.go` (GetUserID/GetTenantID/GetUsername/GetUserRole) |
| OptionalAuthMiddleware | ✅ | `internal/middleware/auth.go:OptionalAuthMiddleware()` |

#### OTP 双因素认证
| 功能 | 状态 | 文件位置 |
|------|------|---------|
| TOTP 服务 | ✅ | `internal/pkg/otp/totp.go` |
| OTP Secret 生成 | ✅ | `internal/pkg/otp/totp.go:GenerateSecret()` |
| OTP 验证 | ✅ | `internal/pkg/otp/totp.go:ValidateCode()` |
| OTP 验证 (带时间窗口) | ✅ | `internal/pkg/otp/totp.go:ValidateCodeWithWindow()` |
| 备用码生成 | ✅ | `internal/pkg/otp/totp.go:GenerateBackupCodes()` |
| Secret 加密 (AES-256-GCM) | ✅ | `internal/pkg/otp/totp.go:EncryptSecret()` |
| Secret 解密 | ✅ | `internal/pkg/otp/totp.go:DecryptSecret()` |
| QR 码生成 | ✅ | `internal/pkg/otp/totp.go:GenerateQRCode()` |
| OTP 登录流程 | ✅ | `internal/api/otp.go:LoginWithOTP()` |
| 临时 Token 机制 | ✅ | `internal/service/otp.go` (TempOTPToken + ExpiresAt) |
| 备用码验证 | ✅ | `internal/service/otp.go:validateBackupCode()` |

#### 用户模型
| 功能 | 状态 | 文件位置 |
|------|------|---------|
| User Model | ✅ | `internal/model/user.go` |
| UserRepository | ✅ | `internal/repository/user.go` |
| 用户角色 (admin/user) | ✅ | `internal/model/user.go:UserRole` |
| OTP 字段支持 | ✅ | `internal/model/user.go` (OTPSecret/OTPEnabled/OTPBackupCodes/TempOTPToken) |
| UserResponse (脱敏) | ✅ | `internal/model/user.go:ToResponse()` |

---

### 2.3 适配器引擎 (Adapter Engine)

| 功能 | 状态 | 文件位置 |
|------|------|---------|
| ClawAdapter 接口定义 | ✅ | `internal/adapter/adapter.go:ClawAdapter` |
| UnifiedConfig 模型 | ✅ | `internal/adapter/adapter.go` (Model/Memory/Server/Logging/Plugins Config) |
| VolumeMount 模型 | ✅ | `internal/adapter/adapter.go` |
| OpenClawAdapter | ✅ | `internal/adapter/openclaw.go` |
| OpenClaw 配置生成 | ✅ | `internal/adapter/openclaw.go:GenerateConfig()` |
| OpenClaw 配置验证 | ✅ | `internal/adapter/openclaw.go:Validate()` |
| OpenClaw 默认配置 | ✅ | `internal/adapter/openclaw.go:GetDefaultOpenClawConfig()` |
| Adapter Factory | ✅ | `internal/adapter/factory.go` |
| Factory 注册机制 | ✅ | `internal/adapter/factory.go:Register()` |
| 工厂支持类型检查 | ✅ | `internal/adapter/factory.go:IsSupported()` |

---

### 2.4 Kubernetes 运行时集成

| 功能 | 状态 | 文件位置 |
|------|------|---------|
| K8S 客户端初始化 | ✅ | `internal/runtime/k8s/client.go` |
| Pod Manager | ✅ | `internal/runtime/k8s/pod.go` |
| ConfigMap Manager | ✅ | `internal/runtime/k8s/configmap.go` |
| Pod 创建/删除/查询 | ✅ | `internal/runtime/k8s/pod.go` (CreatePod/DeletePod/GetPod) |
| Pod 状态查询 | ✅ | `internal/runtime/k8s/pod.go:GetPodStatus()` |
| Pod 日志获取 | ✅ | `internal/runtime/k8s/pod.go:GetPodLogs()` |
| Pod 事件查询 | ✅ | `internal/runtime/k8s/pod.go:GetPodStatus()` (Events included) |
| Pod 就绪等待 | ✅ | `internal/runtime/k8s/pod.go:WaitForPodReady()` |
| 实例ID查找Pod | ✅ | `internal/runtime/k8s/pod.go:GetPodByInstanceID()` |
| ConfigMap 创建/更新/删除 | ✅ | `internal/runtime/k8s/configmap.go` |
| ConfigMap 生成名称 | ✅ | `internal/runtime/k8s/pod.go:GenerateConfigMapName()` |
| Pod 生成名称 | ✅ | `internal/runtime/k8s/pod.go:GeneratePodName()` |

---

### 2.5 前端 (Frontend)

| 功能 | 状态 | 文件位置 |
|------|------|---------|
| 项目初始化 (Vite + React + TS) | ✅ | `frontend/package.json` |
| Ant Design UI 框架 | ✅ | `frontend/package.json` |
| 实例列表页面 | ✅ | `frontend/src/pages/InstanceList.tsx` |
| 创建实例模态框 | ✅ | `frontend/src/components/instances/CreateInstanceModal.tsx` |
| 实例详情页面 | ✅ | `frontend/src/pages/InstanceDetail.tsx` |
| 登录页面 | ✅ | `frontend/src/pages/Login.tsx` |
| OTP 设置页面 | ✅ | `frontend/src/pages/OTPSettings.tsx` |
| 认证 API 客户端 | ✅ | `frontend/src/api/auth.ts` (登录/登出/刷新) |
| OTP API 客户端 | ✅ | `frontend/src/api/otp.ts` |
| 实例 API 客户端 | ✅ | `frontend/src/api/instance.ts` |
| HTTP 客户端封装 | ✅ | `frontend/src/api/client.ts` |
| 状态管理 (Zustand) | ✅ | `frontend/src/store/instance.ts` |
| 路由配置 (React Router) | ✅ | `frontend/src/main.tsx` |
| 静态文件嵌入 | ✅ | `internal/embed/frontend.go` |
| Layout 组件 | ✅ | `frontend/src/components/common/Layout.tsx` |
| 类型定义 | ✅ | `frontend/src/types/index.ts` |
| 格式化工具 | ✅ | `frontend/src/utils/format.ts` |

---

### 2.6 配置管理

| 功能 | 状态 | 文件位置 |
|------|------|---------|
| 配置加载 (Viper) | ✅ | `config/config.go` |
| 配置文件 (YAML) | ✅ | `config/config.yaml` |
| 数据库配置 | ✅ | `config/config.go` |
| K8S 配置 | ✅ | `config/config.go` |
| JWT 配置 | ✅ | `config/config.go` |
| OTP 配置 | ✅ | `config/config.go` |
| 日志配置 | ✅ | `config/config.go` |
| Server 配置 | ✅ | `config/config.go` |

---

### 2.7 数据库表结构 (Model 层)

| 表 | 状态 | 说明 |
|---|------|------|
| users | ✅ | `internal/model/user.go` - 包含 OTP 字段 |
| tenants | ✅ | `internal/model/instance.go` - 包含配额字段 |
| projects | ✅ | `internal/model/instance.go` - 租户关联 |
| config_templates | ✅ | `internal/model/instance.go` - Model定义，Repository接口存在但未实现 |
| claw_instances | ✅ | `internal/model/instance.go` - 包含完整字段 |

---

### 2.8 实例配置管理 (Service 层)

| 功能 | 状态 | 文件位置 |
|------|------|---------|
| 配置适配器调用 | ✅ | `internal/service/instance.go:generateInstanceConfig()` |
| 统一配置解析 | ✅ | `internal/service/instance.go` (OpenClaw) |
| 配置验证 | ✅ | `internal/service/instance.go` (调用 adapter.Validate()) |
| 实例镜像获取 | ✅ | `internal/service/instance.go:getImageForInstance()` |

---

## 三、未实现功能清单

### 3.1 API Gateway Layer (65% 已完成)

| 功能 | 优先级 | 状态 | 说明 |
|------|--------|------|------|
| JWT 认证 | P1 | ✅ 已完成 | 请求认证与授权 |
| OTP 双因素认证 | P1 | ✅ 已完成 | TOTP认证，含备用码和QR码 |
| RBAC 权限控制 | P1 | 🔄 部分完成 | RequireAdmin 中间件完成，缺少细粒度权限 |
| 租户上下文注入 | P1 | ✅ 已完成 | 请求自动携带租户信息 (user_id/tenant_id/role) |
| OptionalAuthMiddleware | P1 | ✅ 已完成 | 可选认证中间件 |
| API 限流保护 | P2 | ❌ 未实现 | 防止 API 滥用 |
| 请求审计日志 | P2 | ❌ 未实现 | 操作审计与追溯 |
| WebSocket 支持 | P2 | ❌ 未实现 | 流式输出和实时通信 |
| 移动端 API Gateway | P3 | ❌ 未实现 | Mobile App 接入支持 |
---

### 3.2 控制面核心 - 未完成项

#### ConfigManager (配置管理器)
| 功能 | 优先级 | 状态 | 说明 |
|------|--------|------|------|
| 配置模板 Model | P1 | ✅ 已完成 | `internal/model/instance.go` 中的 ConfigTemplate |
| ConfigTemplateRepository 接口 | P1 | ✅ 已完成 | `internal/repository/instance.go` 中定义 |
| ConfigTemplateRepository 实现 | P1 | ❌ 未实现 | 数据访问层实现缺失 |
| 配置模板 CRUD Service | P1 | ❌ 未实现 | 业务逻辑层缺失 |
| 配置版本控制 | P2 | ❌ 未实现 | 配置历史版本管理 |
| 批量配置下发 | P1 | ❌ 未实现 | 多实例批量更新配置 |
| 配置一致性校验 | P2 | ❌ 未实现 | 验证配置正确性 |
| 灰度发布 | P2 | ❌ 未实现 | 分阶段发布配置 |
| 配置回滚 | P2 | ❌ 未实现 | 回退到历史版本 |
| 配置模板 API 端点 | P1 | ❌ 未实现 | POST/GET/PUT/DELETE /configs |

#### UsageMonitor (使用监控器)
| 功能 | 优先级 | 状态 | 说明 |
|------|--------|------|------|
| 实例级监控 (CPU/Memory/重启次数) | P1 | 🔄 部分完成 | Pod 状态获取完成，指标聚合未实现 |
| Pod 事件查询 | P1 | ✅ 已完成 | `internal/runtime/k8s/pod.go:GetPodStatus()` |
| Pod 日志查询 | P1 | ✅ 已完成 | `internal/runtime/k8s/pod.go:GetPodLogs()` |
| 使用级监控 (Token/调用次数/错误率) | P1 | ❌ 未实现 | 使用统计与计费 |
| 指标聚合与存储 | P1 | ❌ 未实现 | 时序数据存储 |
| 监控 API 端点 | P1 | ❌ 未实现 | 查询监控数据 |

#### PolicyEngine (策略引擎)
| 功能 | 优先级 | 状态 | 说明 |
|------|--------|------|------|
| 租户资源配额 Model | P1 | ✅ 已完成 | `internal/model/instance.go` Tenant 包含 MaxInstances/MaxCPU/MaxMemory |
| 租户配额检查逻辑 | P1 | ❌ 未实现 | 创建实例前检查资源配额 |
| 访问策略执行 | P2 | ❌ 未实现 | 策略规则引擎 |
| 自动扩缩容策略 | P3 | ❌ 未实现 | 基于使用量的自动调整 |

#### 多租户管理
| 功能 | 优先级 | 状态 | 说明 |
|------|--------|------|------|
| Tenant Model | P1 | ✅ 已完成 | `internal/model/instance.go` |
| Project Model | P1 | ✅ 已完成 | `internal/model/instance.go` |
| TenantRepository 接口 | P1 | ✅ 已完成 | `internal/repository/instance.go` |
| ProjectRepository 接口 | P1 | ✅ 已完成 | `internal/repository/instance.go` |
| TenantRepository 实现 | P1 | ❌ 未实现 | 数据访问层 |
| ProjectRepository 实现 | P1 | ❌ 未实现 | 数据访问层 |
| 租户隔离逻辑 | P1 | ❌ 未实现 | Namespace 级隔离 |
| 租户默认配置覆盖 | P2 | ❌ 未实现 | 全局配置 + 租户覆盖 |

---

### 3.3 Skill & Plugin Manager (0%)

| 功能 | 优先级 | 状态 | 说明 |
|------|--------|------|------|
| Skill CRUD API | P1 | ❌ 未实现 | 技能的增删改查 |
| Plugin CRUD API | P1 | ❌ 未实现 | 插件的增删改查 |
| Skill 版本管理 | P2 | ❌ 未实现 | 技能版本控制 |
| Plugin 版本管理 | P2 | ❌ 未实现 | 插件版本控制 |
| 批量升级 | P2 | ❌ 未实现 | 多实例批量升级 |
| 启用/禁用控制 | P1 | ❌ 未实现 | 开关控制 |
| 内网 Skill Registry | P1 | ❌ 未实现 | 技能仓库 |
| Skill Distributor | P2 | ❌ 未实现 | 统一推送/拉取机制 |

---

### 3.4 Resource Orchestrator Layer - 未完成项

| 功能 | 优先级 | 状态 | 说明 |
|------|--------|------|------|
| Kubernetes 运行时 | P1 | ✅ 已完成 | 使用 kom 库实现 K8S 集成 |
| K8S Pod 管理 | P1 | ✅ 已完成 | 创建/删除/查询/等待就绪 |
| K8S ConfigMap 管理 | P1 | ✅ 已完成 | 创建/更新/删除 |
| K8S 日志查询 | P1 | ✅ 已完成 | Pod 日志获取 |
| K8S 事件查询 | P1 | ✅ 已完成 | Pod 事件获取 |
| PVC 管理 | P1 | ❌ 未实现 | 持久卷管理 |
| Docker 运行时 | P2 | ❌ 未实现 | 直接 Docker 支持 |
| Podman 运行时 | P3 | ❌ 未实现 | Podman 支持 |
| 物理机运行时 | P3 | ❌ 未实现 | 裸机部署支持 |
| 虚拟机运行时 | P3 | ❌ 未实现 | 虚拟机部署支持 |
| 云厂商 API 集成 | P3 | ❌ 未实现 | AWS/阿里云等云服务 |
| 自动扩缩容 | P2 | ❌ 未实现 | 水平扩缩容逻辑 |

---

### 3.5 Storage Layer - 未完成项

| 功能 | 优先级 | 状态 | 说明 |
|------|--------|------|------|
| SQLite 数据库 | P1 | ✅ 已完成 | 使用 GORM + SQLite |
| GORM AutoMigrate | P1 | ✅ 已完成 | `cmd/controlplane/main.go:runMigrations()` |
| 数据库表结构 | P1 | ✅ 已完成 | users/tenants/projects/config_templates/claw_instances |
| 对象存储集成 (MinIO/S3) | P1 | ❌ 未实现 | 文件存储与共享 |
| PVC 管理 | P1 | ❌ 未实现 | Kubernetes 持久卷 |
| 文件上传/下载 API | P1 | ❌ 未实现 | 文件管理接口 |
| 在线文件编辑 | P2 | ❌ 未实现 | Web 编辑器 |
| 文件共享 (多实例) | P2 | ❌ 未实现 | 云盘模式 |
| 数据与实例解耦 | P1 | 🔄 部分完成 | 数据持久化独立于生命周期，但缺少 PVC |

---

### 3.6 前端 - 未完成页面

| 页面 | 优先级 | 状态 | 说明 |
|------|--------|------|------|
| 租户管理页面 | P1 | ❌ 未实现 | 租户列表/创建/编辑 |
| 项目管理页面 | P1 | ❌ 未实现 | 项目列表/创建/编辑 |
| 配置模板管理页面 | P1 | ❌ 未实现 | 配置模板 CRUD |
| 配置版本历史页面 | P2 | ❌ 未实现 | 版本对比与回滚 |
| 技能市场页面 | P1 | ❌ 未实现 | 技能浏览与安装 |
| 插件管理页面 | P1 | ❌ 未实现 | 插件安装与配置 |
| 监控仪表板 | P1 | ❌ 未实现 | 实例指标可视化 |
| 使用统计页面 | P2 | ❌ 未实现 | 资源使用统计 |
| 日志查看页面 | P1 | ❌ 未实现 | 实例日志实时查看 |
| 事件查看页面 | P2 | ❌ 未实现 | 实例事件查看 |
| 文件管理页面 | P1 | ❌ 未实现 | 实例文件管理 |
| 权限管理页面 | P2 | ❌ 未实现 | 用户与角色管理 |

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

#### Sprint 1.1: 认证与授权 ✅ 已完成
- [x] JWT 认证实现 (`internal/pkg/jwt/`)
- [x] 用户登录/登出 API (`internal/api/auth.go`)
- [x] 中间件：租户上下文注入 (`internal/middleware/auth.go`)
- [x] RequireAdmin 中间件 (`internal/middleware/auth.go`)
- [x] 前端登录页面 (`frontend/src/pages/Login.tsx`)
- [x] OTP 双因素认证 (`internal/pkg/otp/`, `internal/service/otp.go`)
- [x] OTP 前端设置页面 (`frontend/src/pages/OTPSettings.tsx`)
- [x] 用户模型与 Repository (`internal/model/user.go`, `internal/repository/user.go`)
- [x] 默认租户和管理员用户初始化 (`cmd/controlplane/main.go`)
#### Sprint 1.2: 配置模板管理
- [ ] ConfigTemplateRepository 实现 (`internal/repository/config_template.go`)
- [ ] 配置模板 CRUD Service (`internal/service/config_template.go`)
- [ ] 配置模板 CRUD API (`internal/api/config.go`)
- [ ] 前端配置模板管理页面
- [ ] 配置模板关联到实例创建

#### Sprint 1.3: 租户与项目管理
- [ ] TenantRepository 实现 (`internal/repository/tenant.go`)
- [ ] ProjectRepository 实现 (`internal/repository/project.go`)
- [ ] 租户/项目 CRUD API (`internal/api/tenant.go`, `internal/api/project.go`)
- [ ] 前端租户/项目管理页面
- [ ] 租户隔离逻辑实现 (资源查询过滤)

#### Sprint 1.4: 实例监控与日志
- [ ] 监控指标采集 (CPU/Memory)
- [ ] 日志查询 API (`GET /instances/:id/logs`)
- [ ] 日志流 WebSocket (`WS /instances/:id/logs/stream`)
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
| TODO 注释清理 | P2 | `internal/service/instance.go` 有 TODO 未完成 |
| 测试覆盖 | P1 | 缺少单元测试和集成测试 |
| 文档完善 | P2 | API 文档缺失 |
| 错误处理 | P1 | 部分错误被忽略 (如 ConfigMap 创建失败) |
| 分页参数解析 | P2 | `internal/api/handler.go` page/pageSize 未正确解析 |

### 5.2 架构改进

| 问题 | 优先级 | 说明 |
|------|--------|------|
| 状态同步机制 | P1 | 目前为异步 goroutine (`syncPodStatus`)，缺少错误处理和重试 |
| 配置下发优化 | P2 | 每次更新都创建 ConfigMap，可增加 diff 检查 |
| 事件驱动 | P2 | 缺少事件总线 (NATS/Kafka) |
| 重试机制 | P1 | K8S 操作缺少重试 (使用 kom 库) |
| 计数查询缺失 | P1 | `ListInstances` 返回实例数量而非总数 |
| 粗粒度错误处理 | P2 | 大量使用 `fmt.Errorf`，可考虑自定义错误类型 |

### 5.3 性能优化

| 问题 | 优先级 | 说明 |
|------|--------|------|
| 数据库连接池 | P2 | SQLite 连接管理优化 (当前使用默认设置) |
| API 响应缓存 | P2 | 静态数据缓存 (如租户列表) |
| 分页查询优化 | P1 | `ListInstances` 缺少总数查询 |
| 资源限制 | P2 | 防止资源耗尽 (K8S Pod 资源限制已实现) |
| API 响应缓存 | P2 | 静态数据缓存 |
| 分页查询优化 | P1 | 目前未实现总数查询 |
| 资源限制 | P2 | 防止资源耗尽 |

---

## 六、依赖与集成点

### 6.1 外部依赖

| 组件 | 用途 | 集成状态 | 备注 |
|------|------|---------|------|
| Kubernetes | 容器编排 | ✅ 已集成 | 使用 kom 库 |
| SQLite | 元数据存储 | ✅ 已使用 | GORM + SQLite |
| MinIO/S3 | 对象存储 | ❌ 未集成 | 计划中 |
| Prometheus | 指标采集 | ❌ 未集成 | 计划中 |
| Grafana | 可视化 | ❌ 未集成 | 计划中 |
| NATS/Kafka | 事件队列 | ❌ 未集成 | 计划中 |

### 6.2 核心库依赖

| 库 | 用途 | 版本 |
|------|------|------|
| gin-gonic/gin | Web 框架 | 最新 |
| gorm.io/gorm | ORM | 最新 |
| gorm.io/driver/sqlite | SQLite 驱动 | 最新 |
| golang-jwt/jwt | JWT 实现 | v5 |
| golang.org/x/crypto/bcrypt | 密码哈希 | 最新 |
| github.com/pquerna/otp | TOTP 实现 | 最新 |
| github.com/pquerna/otp/totp | TOTP | 最新 |
| github.com/skip2/go-qrcode | QR 码生成 | 最新 |
| github.com/weibaohui/kom | K8S 客户端 | 最新 |
| gopkg.in/yaml.v3 | YAML 解析 | 最新 |
| github.com/google/uuid | UUID 生成 | 最新 |
| viper | 配置管理 | 最新 |

### 6.3 Claw 依赖

| 组件 | 说明 | 状态 |
|------|------|------|
| OpenClaw | 目标类型 1 | ✅ Adapter 已实现 (`internal/adapter/openclaw.go`) |
| NanoClaw | 目标类型 2 | ❌ Adapter 未实现 |
| 自定义 Claw | 扩展支持 | ✅ 接口已定义 (`internal/adapter/adapter.go`) |

---

## 七、风险评估

| 风险 | 级别 | 缓解措施 | 当前状态 |
|------|--------|----------|---------|
| K8S 操作失败影响状态一致性 | 中 | 增加重试机制 + 状态修复任务 | 部分实现 (缺少重试) |
| 多租户隔离不彻底 | 高 | 加强中间件 + 数据库查询隔离 | 中间件完成，隔离逻辑未实现 |
| 配置下发失败导致实例启动异常 | 中 | 配置验证 + 回滚机制 | 部分实现 (有验证，无回滚) |
| 高并发下资源竞争 | 中 | 使用分布式锁 + 事件队列 | 未实现 |
| 数据库迁移风险 | 中 | 版本化迁移脚本 + 回滚方案 | 使用 GORM AutoMigrate |
| OTP Secret 加密密钥泄露 | 高 | 使用环境变量 + 密钥轮换 | 使用环境变量 |
| SQLite 单写性能限制 | 中 | 后期迁移到 PostgreSQL | 计划中 |
| 前端状态管理混乱 | 低 | 使用 Zustand + 统一 API 调用 | Zustand 已实现 |

---

## 八、附录

### 8.1 文件索引

| 模块 | 关键文件 | 说明 |
|------|----------|------|
| 主程序 | `cmd/controlplane/main.go` | 服务入口、数据库初始化、路由设置 |
| API 层 | `internal/api/`, `internal/router.go` | RESTful API 路由、Handler |
| 认证 API | `internal/api/auth.go`, `internal/api/otp.go` | 登录、OTP、用户管理 |
| 服务层 | `internal/service/` | 业务逻辑层 |
| 数据访问层 | `internal/repository/`, `internal/model/` | 数据持久化 |
| 中间件 | `internal/middleware/` | 认证、权限控制 |
| 适配器 | `internal/adapter/` | OpenClaw 适配器 |
| K8S 运行时 | `internal/runtime/k8s/` | Pod、ConfigMap 管理 |
| JWT | `internal/pkg/jwt/jwt.go` | JWT Token 生成和验证 |
| OTP | `internal/pkg/otp/totp.go` | TOTP 双因素认证 |
| 配置 | `config/` | 服务配置 |
| 前端 | `frontend/`, `internal/embed/frontend.go` | React 前端应用 |

### 8.2 数据库表详情

**users 表** (已实现)
```sql
id, username, password_hash, tenant_id, role, is_active,
otp_secret, otp_enabled, otp_backup_codes, temp_otp_token, temp_otp_token_expires_at,
created_at, updated_at
```

**tenants 表** (已实现)
```sql
id, name, max_instances, max_cpu, max_memory, max_storage,
created_at, updated_at
```

**projects 表** (已实现)
```sql
id, tenant_id (FK), name, created_at, updated_at
```

**config_templates 表** (已实现 Model，缺 Repository)
```sql
id, name, description, variables (BLOB), adapter_type, version,
created_at, updated_at
```

**claw_instances 表** (已实现)
```sql
id, name, tenant_id (FK), project_id (FK), type, version,
status, config (BLOB), cpu, memory, config_dir, data_dir,
storage_size, created_at, updated_at
```

### 8.3 关键技术细节

#### 认证流程

1. **普通登录流程** (`LoginWithOTP` → 无 OTP)
   ```
   POST /api/v1/auth/login (username, password)
   → 验证密码
   → 检查 OTP 状态
   → 如果未启用 OTP: 直接返回 access_token + refresh_token
   ```

2. **OTP 登录流程** (`LoginWithOTP` → 有 OTP)
   ```
   POST /api/v1/auth/login (username, password)
   → 验证密码
   → 检查 OTP 状态
   → 如果启用 OTP: 返回 temp_token (5分钟有效期)

   POST /api/v1/auth/otp/verify (temp_token, otp_code)
   → 验证 temp_token
   → 验证 OTP (或备用码)
   → 返回 access_token + refresh_token
   ```

3. **OTP 设置流程**
   ```
   POST /api/v1/auth/otp/generate
   → 生成 secret + qr_code
   → 存储加密 secret (未启用)

   POST /api/v1/auth/otp/enable (code)
   → 验证 code
   → 生成备用码
   → 启用 OTP

   GET /api/v1/auth/otp/backup
   → 返回备用码 (仅一次)
   ```

#### 实例生命周期

```
Creating (创建中)
    ↓ (Pod 创建成功)
Running (运行中)
    ↓ (Stop)
Stopped (已停止)
    ↓ (Start)
Running
    ↓ (失败)
Failed (失败)
    ↓ (Delete)
Destroyed (已销毁)
```

#### K8S 资源命名规则

```go
Pod 名称:      claw-{instanceID}
ConfigMap 名称: claw-config-{instanceID}
Label:         app=claw, instanceId={instanceID}, tenantId={tenantID}
```

---

**文档版本:** v1.0
**最后更新:** 2026-02-26
**维护者:** Open Cluster Claw 开发团队

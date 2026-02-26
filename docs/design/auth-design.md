# 认证与授权设计文档

> 变更记录表

| 版本 | 日期 | 作者 | 变更说明 |
|------|------|------|----------|
| v1.0 | 2026-02-27 | AI | 初始版本 |

---

## 一、要解决的问题

当前系统缺乏用户身份验证和授权机制，任何访问 API 的用户都可以执行所有操作，存在严重的安全风险。需要实现：
1. 用户身份验证机制
2. 租户隔离和上下文管理
3. API 访问保护

---

## 二、功能目标与边界

### 2.1 功能目标

1. **JWT 认证**
   - 生成和验证 JWT Token
   - Token 包含用户基本信息和租户信息
   - 支持访问令牌（Access Token）和刷新令牌（Refresh Token）

2. **用户登录/登出**
   - 用户登录接口返回 JWT Token
   - 登出使 Token 失效（通过维护黑名单或设置过期时间）

3. **租户上下文注入中间件**
   - 从 JWT 中提取租户信息
   - 将租户 ID 和用户信息注入到请求上下文
   - 保护需要认证的 API 端点

4. **前端登录页面**
   - 提供登录表单界面
   - 处理登录请求和 Token 存储
   - 自动在请求中携带 Token

### 2.2 功能边界

**本阶段不包含：**
- 用户注册功能（使用默认管理员账户或手动创建）
- RBAC 权限控制（基于角色的访问控制）
- OAuth2 集成
- 多因素认证（MFA）
- Token 黑名单机制（使用过期时间简化）

**本阶段包含：**
- 简单的用户模型（ID、用户名、密码哈希）
- 用户数据持久化
- 密码加密存储
- JWT Token 生成与验证
- 基本的认证中间件

---

## 三、核心设计思路

### 3.1 数据模型设计

#### users 表

```sql
CREATE TABLE IF NOT EXISTS users (
    id TEXT PRIMARY KEY,
    username TEXT NOT NULL UNIQUE,
    password_hash TEXT NOT NULL,
    tenant_id TEXT REFERENCES tenants(id) ON DELETE CASCADE,
    role TEXT DEFAULT 'user',
    is_active BOOLEAN DEFAULT true,
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    updated_at DATETIME DEFAULT CURRENT_TIMESTAMP
)
```

**字段说明：**
- `id`: 用户唯一标识符（UUID）
- `username`: 登录用户名
- `password_hash`: 密码哈希（使用 bcrypt）
- `tenant_id`: 所属租户 ID（管理员用户可以为空）
- `role`: 用户角色（admin/user）
- `is_active`: 账户是否激活

### 3.2 JWT Token 设计

#### Payload 结构

```go
type Claims struct {
    UserID   string `json:"user_id"`
    Username string `json:"username"`
    TenantID string `json:"tenant_id"`
    Role     string `json:"role"`
    jwt.RegisteredClaims
}
```

#### Token 类型

1. **Access Token**
   - 有效期：24 小时（可配置）
   - 用途：访问受保护的 API
   - 存储方式：前端内存（或 localStorage）

2. **Refresh Token**
   - 有效期：7 天（可配置）
   - 用途：刷新 Access Token
   - 存储方式：前端内存（或 localStorage）

### 3.3 API 设计

#### 认证相关 API

| 方法 | 路径 | 说明 |
|------|------|------|
| POST | `/api/v1/auth/login` | 用户登录 |
| POST | `/api/v1/auth/logout` | 用户登出 |
| POST | `/api/v1/auth/refresh` | 刷新 Token |

#### 登录请求

```json
{
  "username": "admin",
  "password": "password123"
}
```

#### 登录响应

```json
{
  "code": 0,
  "message": "success",
  "data": {
    "access_token": "eyJhbGciOiJIUzI1NiIs...",
    "refresh_token": "eyJhbGciOiJIUzI1NiIs...",
    "expires_in": 86400,
    "user": {
      "id": "user-123",
      "username": "admin",
      "tenant_id": "tenant-001",
      "role": "admin"
    }
  }
}
```

### 3.4 中间件设计

#### AuthMiddleware

职责：
1. 从请求头 `Authorization: Bearer <token>` 提取 Token
2. 验证 Token 有效性
3. 将用户信息注入到 `gin.Context`
4. 验证失败返回 401 未授权

```go
func AuthMiddleware(jwtService JWTService) gin.HandlerFunc {
    return func(c *gin.Context) {
        token := c.GetHeader("Authorization")
        // 验证逻辑...
        c.Next()
    }
}
```

#### TenantContextMiddleware

职责：
1. 从已认证的 Context 中提取租户 ID
2. 将租户 ID 设置到请求上下文中
3. 确保后续处理可以使用租户信息进行数据隔离

### 3.5 密码安全

- 使用 `golang.org/x/crypto/bcrypt` 进行密码哈希
- Hash 算法：bcrypt with cost 12
- 验证时使用 CompareHashAndPassword

---

## 四、可能影响的模块

### 4.1 新增模块

| 模块 | 位置 | 说明 |
|------|------|------|
| User Model | `internal/model/user.go` | 用户数据模型 |
| User Repository | `internal/repository/user.go` | 用户数据访问层 |
| Auth Service | `internal/service/auth.go` | 认证业务逻辑 |
| JWT Service | `internal/pkg/jwt/jwt.go` | JWT 生成和验证 |
| Auth Handler | `internal/api/auth.go` | 认证 API 处理器 |
| Auth Middleware | `internal/middleware/auth.go` | 认证中间件 |

### 4.2 修改模块

| 模块 | 修改内容 |
|------|----------|
| Database Migration | 添加 users 表创建语句 |
| Router | 添加认证路由和中间件配置 |
| Frontend | 添加登录页面和认证逻辑 |
| API Client | 添加 Token 拦截器 |

### 4.3 配置文件修改

`config/config.yaml` 中已有 JWT 配置，无需修改：
```yaml
jwt:
  secret: your-secret-key-change-in-production
  expire_time: 86400 # 24 hours in seconds
```

---

## 五、实现步骤

### Step 1: 数据库和模型
1. 创建 users 数据库表
2. 实现 User Model
3. 实现 UserRepository

### Step 2: JWT 服务
1. 实现 JWT Token 生成
2. 实现 JWT Token 验证
3. 实现 Claims 结构

### Step 3: 认证服务
1. 实现密码哈希和验证
2. 实现用户登录逻辑
3. 实现用户登出逻辑（可选，使用过期时间）
4. 实现刷新 Token 逻辑

### Step 4: 中间件
1. 实现 AuthMiddleware
2. 实现 TenantContextMiddleware
3. 集成到 Router

### Step 5: API Handler
1. 实现登录 Handler
2. 实现登出 Handler（可选）
3. 实现刷新 Token Handler
4. 添加路由配置

### Step 6: 前端实现
1. 创建登录页面组件
2. 实现 API 客户端认证拦截器
3. 实现 Token 存储和管理
4. 更新路由配置

### Step 7: 集成测试
1. 测试登录流程
2. 测试 Token 验证
3. 测试未授权访问
4. 测试租户隔离

---

## 六、已知限制与待改进点

### 6.1 已知限制

1. **Token 撤销问题**
   - 当前设计不支持主动撤销 Token
   - 用户登出只是客户端删除 Token
   - 解决方案：引入 Token 黑名单或使用短有效期

2. **密码重置**
   - 不支持密码重置功能
   - 需要管理员手动修改

3. **用户管理**
   - 不支持用户注册
   - 需要通过数据库直接创建用户

### 6.2 待改进点

1. **RBAC 权限控制**（Phase 3）
   - 基于角色的权限控制
   - 细粒度的 API 访问控制

2. **OAuth2 集成**（Phase 5）
   - 支持 OAuth2 认证
   - 支持第三方登录

3. **多因素认证**（Phase 5）
   - 支持 TOTP
   - 支持 SMS 验证码

4. **审计日志**（Phase 4）
   - 记录用户登录/登出操作
   - 记录敏感操作

---

## 七、技术栈

| 组件 | 技术选型 |
|------|----------|
| Web 框架 | Gin |
| JWT | golang-jwt/jwt/v5 |
| 密码哈希 | golang.org/x/crypto/bcrypt |
| 数据库 | SQLite |
| 前端框架 | React + TypeScript |
| 前端 UI | Ant Design |

---

## 八、安全性考虑

1. **密码安全**
   - 使用 bcrypt 进行密码哈希
   - 不在日志中记录明文密码
   - 前端传输使用 HTTPS（生产环境）

2. **Token 安全**
   - JWT Secret 必须配置强密钥
   - Token 过期时间不宜过长
   - 不在 URL 中传递 Token

3. **输入验证**
   - 验证用户名格式
   - 验证密码强度
   - 防止 SQL 注入

4. **敏感信息处理**
   - 不在响应中返回密码哈希
   - 登录失败不泄露用户名是否存在
   - 使用通用错误消息

---

**文档版本:** v1.0
**最后更新:** 2026-02-27

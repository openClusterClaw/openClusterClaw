# 认证与授权实现总结

> 变更记录表

| 版本 | 日期 | 作者 | 变更说明 |
|------|------|------|----------|
| v1.0 | 2026-02-27 | AI | 初始版本：完成 JWT 认证和授权功能 |

---

## 一、实现内容

### 1.1 后端实现

#### 数据库层
- **用户模型** (`internal/model/user.go`)
  - 定义了 `User` 数据模型
  - 定义了 `UserRole` 类型（admin/user）
  - 定义了 `UserResponse` 用于 API 响应（不包含敏感信息）

- **用户仓储** (`internal/repository/user.go`)
  - `Create` - 创建新用户
  - `GetByID` - 按 ID 获取用户
  - `GetByUsername` - 按用户名获取用户
  - `Update` - 更新用户信息
  - `Delete` - 删除用户
  - `List` - 获取用户列表（支持租户过滤和分页）
  - `GenerateID` - 生成用户 ID

- **数据库迁移** (`cmd/controlplane/main.go`)
  - 添加 `users` 表创建语句
  - 添加索引：`idx_users_tenant`, `idx_users_username`

#### JWT 服务
- **JWT 服务** (`internal/pkg/jwt/jwt.go`)
  - `GenerateAccessToken` - 生成访问令牌
  - `GenerateRefreshToken` - 生成刷新令牌
  - `ValidateToken` - 验证令牌有效性
  - 使用 HMAC-SHA256 签名算法
  - Claims 包含：UserID, Username, TenantID, Role

#### 认证服务
- **认证服务** (`internal/service/auth.go`)
  - `Login` - 用户登录，返回访问令牌和刷新令牌
  - `RefreshToken` - 使用刷新令牌获取新的访问令牌
  - `CreateUser` - 创建新用户（管理员功能）
  - `GetUserByID` - 获取用户信息
  - `HashPassword` - 密码哈希（bcrypt）
  - `GenerateUserID` - 生成用户 ID

#### 中间件
- **认证中间件** (`internal/middleware/auth.go`)
  - `AuthMiddleware` - 验证 JWT 并注入用户信息到上下文
  - `RequireAdmin` - 要求管理员权限
  - `OptionalAuthMiddleware` - 可选认证
  - 上下文键：`user_id`, `username`, `tenant_id`, `role`
  - 工具函数：`GetUserID`, `GetTenantID`, `GetUsername`, `GetUserRole`

#### API 处理器
- **认证处理器** (`internal/api/auth.go`)
  - `Login` - 处理登录请求
  - `RefreshToken` - 处理令牌刷新请求
  - `Logout` - 处理登出请求
  - `GetCurrentUser` - 获取当前用户信息
  - `CreateUser` - 创建用户（需要管理员权限）

#### 路由配置
- **更新路由** (`internal/api/router.go`)
  - 公开路由：`/api/v1/auth/login`, `/api/v1/auth/refresh`
  - 需要认证的路由：`/api/v1/instances/*`
  - 受保护的用户路由：`/api/v1/auth/*`（需要认证）
  - 管理员路由：`/api/v1/auth/users`（需要管理员权限）

#### 主程序
- **初始化更新** (`cmd/controlplane/main.go`)
  - 添加 JWT 服务初始化
  - 添加认证服务初始化
  - 添加默认租户和管理员用户创建逻辑
  - 默认账户：username=`admin`, password=`admin123`

### 1.2 前端实现

#### API 客户端
- **认证 API** (`frontend/src/api/auth.ts`)
  - `authApi.login` - 登录
  - `authApi.logout` - 登出
  - `authApi.getCurrentUser` - 获取当前用户
  - `authApi.refreshToken` - 刷新令牌
  - `tokenManager` - Token 管理工具
    - 存储到 localStorage
    - 自动在请求中添加 Authorization 头
    - 自动刷新令牌（401 时）
    - 检查认证状态

#### 页面组件
- **登录页面** (`frontend/src/pages/Login.tsx`)
  - 登录表单（用户名/密码）
  - 集成认证 API
  - 成功后跳转到首页
  - 显示默认账户信息

#### 布局组件更新
- **应用布局** (`frontend/src/App.tsx`)
  - 添加 `ProtectedRoute` 组件
  - 检查认证状态
  - 未认证时重定向到登录页
  - 添加登录路由

- **公共布局** (`frontend/src/components/common/Layout.tsx`)
  - 添加 Header 显示用户信息
  - 添加用户下拉菜单
  - 添加退出登录功能
  - 显示用户名和角色

---

## 二、与需求对应关系

| 需求项 | 实现方式 | 状态 |
|---------|---------|------|
| JWT 认证实现 | JWT 服务 + 认证服务 | ✅ 完成 |
| 用户登录/登出 API | `/api/v1/auth/login`, `/api/v1/auth/logout` | ✅ 完成 |
| 中间件：租户上下文注入 | AuthMiddleware 注入 tenant_id | ✅ 完成 |
| 前端登录页面 | `Login.tsx` | ✅ 完成 |

---

## 三、关键实现点

### 3.1 密码安全
- 使用 bcrypt 进行密码哈希（cost 12）
- 不在响应中返回密码哈希
- 不在日志中记录明文密码

### 3.2 Token 安全
- 使用强密钥（通过配置文件设置）
- Access Token 有效期 24 小时
- Refresh Token 有效期 7 天
- Token 包含租户信息，用于数据隔离

### 3.3 租户隔离
- JWT Claims 包含 tenant_id
- 中间件将租户 ID 注入到请求上下文
- 后续处理可通过 `middleware.GetTenantID(c)` 获取租户信息

### 3.4 自动 Token 刷新
- 前端 axios 拦截器捕获 401 错误
- 自动使用 refresh token 获取新 access token
- 重试原始请求
- 刷新失败时清除 token 并重定向到登录页

---

## 四、已知限制

1. **Token 撤销问题**
   - 当前设计不支持主动撤销 Token
   - 用户登出只是客户端删除 Token
   - 依赖 Token 过期时间保证安全性

2. **密码重置**
   - 不支持密码重置功能
   - 需要管理员手动修改数据库

3. **用户注册**
   - 不支持用户注册
   - 需要通过管理员创建用户

4. **RBAC 权限控制**
   - 当前仅有基本角色区分（admin/user）
   - 不支持细粒度的权限控制

5. **前端构建**
   - ✅ 前端依赖已安装
   - ✅ 前端构建成功
   - 建议在生产环境中预先构建静态文件

---

## 五、待改进点（按 Phase）

### Phase 3: 企业能力
- RBAC 权限控制
- 用户注册功能
- 密码重置功能

### Phase 4: 运维与扩展
- Token 黑名单机制
- 操作审计日志
- 多因素认证（MFA）

### Phase 5: 移动端与高级特性
- OAuth2 集成
- 第三方登录支持

---

## 六、API 端点列表

| 方法 | 路径 | 认证要求 | 说明 |
|------|------|-----------|------|
| POST | `/api/v1/auth/login` | 无需认证 | 用户登录 |
| POST | `/api/v1/auth/refresh` | 无需认证 | 刷新 Token |
| POST | `/api/v1/auth/logout` | 需要认证 | 用户登出 |
| GET | `/api/v1/auth/me` | 需要认证 | 获取当前用户 |
| POST | `/api/v1/auth/users` | 需要认证+管理员 | 创建用户 |

---

## 七、配置说明

### JWT 配置 (`config/config.yaml`)
```yaml
jwt:
  secret: your-secret-key-change-in-production  # JWT 签名密钥（必须修改）
  expire_time: 86400  # Access Token 有效期（秒），默认 24 小时
```

### 默认账户
- 用户名：`admin`
- 密码：`admin123`
- **重要**：首次登录后请立即修改密码！

---

## 八、测试建议

### 8.1 后端测试
```bash
# 登录测试
curl -X POST http://localhost:8080/api/v1/auth/login \
  -H "Content-Type: application/json" \
  -d '{"username":"admin","password":"admin123"}'

# 访问受保护资源（使用返回的 token）
curl http://localhost:8080/api/v1/auth/me \
  -H "Authorization: Bearer <access_token>"

# 创建用户（需要管理员权限）
curl -X POST http://localhost:8080/api/v1/auth/users \
  -H "Authorization: Bearer <access_token>" \
  -H "Content-Type: application/json" \
  -d '{"username":"testuser","password":"password123"}'
```

### 8.2 前端测试
1. 启动后端服务
2. 启动前端开发服务器：`cd frontend && npm run dev`
3. 访问 http://localhost:5173
4. 使用默认账户登录
5. 测试实例管理功能

---

**文档版本:** v1.0
**最后更新:** 2026-02-27
**对应设计文档:** `docs/design/auth-design.md`

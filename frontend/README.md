# Open Cluster Claw - Frontend

Open Cluster Claw 前端应用

## 技术栈

- **框架**: React 18
- **语言**: TypeScript
- **构建工具**: Vite
- **UI 组件库**: Ant Design 5
- **路由**: React Router 6
- **状态管理**: Zustand
- **HTTP 客户端**: Axios
- **包管理器**: pnpm

## 项目结构

```
frontend/
├── src/
│   ├── api/              # API 接口
│   ├── components/       # 组件
│   │   ├── common/      # 通用组件
│   │   ├── instances/   # 实例相关组件
│   │   └── config/      # 配置相关组件
│   ├── pages/           # 页面组件
│   ├── hooks/           # 自定义 Hooks
│   ├── services/        # 业务逻辑
│   ├── store/           # 状态管理
│   ├── types/           # TypeScript 类型定义
│   ├── utils/           # 工具函数
│   ├── App.tsx          # 根组件
│   └── main.tsx         # 应用入口
├── public/              # 静态资源
├── index.html           # HTML 模板
├── vite.config.ts       # Vite 配置
├── tsconfig.json        # TypeScript 配置
└── package.json         # 项目配置
```

## 快速开始

### 安装依赖

```bash
pnpm install
```

### 开发模式

```bash
pnpm dev
```

访问 http://localhost:3000

### 构建

```bash
pnpm build
```

### 预览构建

```bash
pnpm preview
```

## 命令别名

也可以使用 Make 命令：

```bash
make install    # 安装依赖
make dev        # 启动开发服务器
make build      # 构建生产版本
make lint       # 代码检查
make format     # 代码格式化
```

## 开发规范

详见项目文档 `docs/开发规范/前端规范/` 目录。

### 组件编写

- 使用函数组件和 Hooks
- 组件文件命名采用 PascalCase
- Props 使用 TypeScript 类型定义

### 样式规范

- 使用 Ant Design 组件库
- 使用 CSS Modules 或 styled-components（待实现）

### API 调用

- 使用 `src/api/` 目录下定义的 API 方法
- 统一使用 axios 客户端
- 错误处理在拦截器中统一处理

### 状态管理

- 使用 Zustand 进行全局状态管理
- 局部状态使用 React useState

## 环境变量

创建 `.env` 文件：

```
VITE_API_BASE_URL=http://localhost:8080/api
```

## License

MIT
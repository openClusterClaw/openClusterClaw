# API 调用规范

本文档定义 API 服务的封装方式和错误处理。

## API 服务文件结构

```tsx
// services/api.ts
import axios from 'axios';
import type { Instance, Config, Tenant } from '@/types';

const API_BASE = import.meta.env.VITE_API_BASE || 'http://localhost:8080/api';

const api = axios.create({
  baseURL: API_BASE,
  headers: { 'Content-Type': 'application/json' },
});

// 按资源分组
export const instanceApi = {
  list: (params?: ListParams) => api.get<Instance[]>('/instances', { params }),
  get: (id: string) => api.get<Instance>(`/instances/${id}`),
  create: (data: CreateInstanceRequest) => api.post<Instance>('/instances', data),
  update: (id: string, data: UpdateInstanceRequest) => api.put<Instance>(`/instances/${id}`, data),
  delete: (id: string) => api.delete(`/instances/${id}`),
  start: (id: string) => api.post(`/instances/${id}/start`),
  stop: (id: string) => api.post(`/instances/${id}/stop`),
  restart: (id: string) => api.post(`/instances/${id}/restart`),
};

export const configApi = {
  list: () => api.get<Config[]>('/configs'),
  get: (id: string) => api.get<Config>(`/configs/${id}`),
  create: (data: CreateConfigRequest) => api.post<Config>('/configs', data),
  update: (id: string, data: UpdateConfigRequest) => api.put<Config>(`/configs/${id}`, data),
  delete: (id: string) => api.delete(`/configs/${id}`),
  publish: (id: string) => api.post(`/configs/${id}/publish`),
  rollback: (id: string) => api.post(`/configs/${id}/rollback`),
};

export const tenantApi = {
  list: () => api.get<Tenant[]>('/tenants'),
  get: (id: string) => api.get<Tenant>(`/tenants/${id}`),
  create: (data: CreateTenantRequest) => api.post<Tenant>('/tenants', data),
  update: (id: string, data: UpdateTenantRequest) => api.put<Tenant>(`/tenants/${id}`, data),
  delete: (id: string) => api.delete(`/tenants/${id}`),
};
```

## 错误处理

```tsx
import { message } from 'antd';

// 统一错误处理
const handleApiError = (error: unknown) => {
  if (axios.isAxiosError(error)) {
    const errorMsg = error.response?.data?.message || error.message;
    message.error(errorMsg);
  }
};

// 使用示例
try {
  await instanceApi.create(data);
  message.success('Instance created');
} catch (error) {
  handleApiError(error);
}
```

## 请求拦截器

```tsx
// 请求拦截器（如需添加 token）
api.interceptors.request.use(
  (config) => {
    const token = localStorage.getItem('token');
    if (token) {
      config.headers.Authorization = `Bearer ${token}`;
    }
    return config;
  },
  (error) => Promise.reject(error)
);

// 响应拦截器
api.interceptors.response.use(
  (response) => response,
  (error) => {
    if (error.response?.status === 401) {
      // 处理未授权
    }
    return Promise.reject(error);
  }
);
```

## 在组件中使用

```tsx
function InstanceList() {
  const [instances, setInstances] = useState<Instance[]>([]);
  const [loading, setLoading] = useState(true);

  useEffect(() => {
    fetchInstances();
  }, []);

  const fetchInstances = async () => {
    setLoading(true);
    try {
      const res = await instanceApi.list();
      setInstances(res.data);
    } catch (error) {
      handleApiError(error);
    } finally {
      setLoading(false);
    }
  };

  const handleDelete = async (id: string) => {
    try {
      await instanceApi.delete(id);
      message.success('Instance deleted');
      fetchInstances(); // 刷新列表
    } catch (error) {
      handleApiError(error);
    }
  };

  // ...
}
```

## API 命名规范

| 操作 | HTTP 方法 | 命名 | 示例 |
| ---- | --------- | ---- | -------------------------------- |
| 获取列表 | GET | list | `instanceApi.list()` |
| 获取单个 | GET | get | `instanceApi.get(id)` |
| 创建 | POST | create | `instanceApi.create(data)` |
| 更新 | PUT/PATCH | update | `instanceApi.update(id, data)` |
| 删除 | DELETE | delete | `instanceApi.delete(id)` |
| 自定义操作 | POST | 动词 | `instanceApi.start(id)` |

## 相关文档

- [类型定义规范](./10-类型定义规范.md)
- [状态管理规范](./06-状态管理规范.md)
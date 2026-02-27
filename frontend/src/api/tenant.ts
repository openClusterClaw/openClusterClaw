import client from './client';
import type { ApiResponse } from '@/types';

export interface Tenant {
  id: string;
  name: string;
  max_instances: number;
  max_cpu: string;
  max_memory: string;
  max_storage: string;
  created_at: string;
  updated_at: string;
}

export interface Project {
  id: string;
  tenant_id: string;
  name: string;
  created_at: string;
  updated_at: string;
}

export interface CreateTenantRequest {
  name: string;
  max_instances?: number;
  max_cpu?: string;
  max_memory?: string;
  max_storage?: string;
}

export interface UpdateTenantRequest {
  name?: string;
  max_instances?: number;
  max_cpu?: string;
  max_memory?: string;
  max_storage?: string;
}

export interface CreateProjectRequest {
  tenant_id: string;
  name: string;
}

export interface UpdateProjectRequest {
  name?: string;
}

export interface TenantListResponse {
  tenants: Tenant[];
  total: number;
  page: number;
  page_size: number;
}

export interface ProjectListResponse {
  projects: Project[];
  total: number;
  page: number;
  page_size: number;
}

export const tenantApi = {
  async listTenants(): Promise<Tenant[]> {
    const { data } = await client.get<ApiResponse<TenantListResponse>>('/tenants');
    return data.data?.tenants || [];
  },

  async getTenant(id: string): Promise<Tenant> {
    const { data } = await client.get<ApiResponse<Tenant>>(`/tenants/${id}`);
    return data.data!;
  },

  async createTenant(reqData: CreateTenantRequest): Promise<Tenant> {
    const { data } = await client.post<ApiResponse<Tenant>>('/tenants', reqData);
    return data.data!;
  },

  async updateTenant(id: string, reqData: UpdateTenantRequest): Promise<Tenant> {
    const { data } = await client.put<ApiResponse<Tenant>>(`/tenants/${id}`, reqData);
    return data.data!;
  },

  async deleteTenant(id: string): Promise<void> {
    await client.delete(`/tenants/${id}`);
  },
};

export const projectApi = {
  async listProjects(tenantId?: string): Promise<Project[]> {
    const params = tenantId ? `?tenant_id=${tenantId}` : '';
    const { data } = await client.get<ApiResponse<ProjectListResponse>>(`/projects${params}`);
    return data.data?.projects || [];
  },

  async getProject(id: string): Promise<Project> {
    const { data } = await client.get<ApiResponse<Project>>(`/projects/${id}`);
    return data.data!;
  },

  async createProject(reqData: CreateProjectRequest): Promise<Project> {
    const { data } = await client.post<ApiResponse<Project>>('/projects', reqData);
    return data.data!;
  },

  async updateProject(id: string, reqData: UpdateProjectRequest): Promise<Project> {
    const { data } = await client.put<ApiResponse<Project>>(`/projects/${id}`, reqData);
    return data.data!;
  },

  async deleteProject(id: string): Promise<void> {
    await client.delete(`/projects/${id}`);
  },
};
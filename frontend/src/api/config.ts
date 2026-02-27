import client from './client';
import type { ApiResponse } from '@/types';

export interface TemplateVariable {
  name: string;
  type: string;
  default: any;
  required: boolean;
  description: string;
  secret: boolean;
}

export interface ConfigTemplate {
  id: string;
  name: string;
  description: string;
  variables: TemplateVariable[];
  adapter_type: string;
  version: string;
  created_at: string;
  updated_at: string;
}

export interface CreateTemplateRequest {
  name: string;
  description?: string;
  variables: TemplateVariable[];
  adapter_type: string;
  version?: string;
}

export interface UpdateTemplateRequest {
  name?: string;
  description?: string;
  variables?: TemplateVariable[];
  version?: string;
}

export interface ConfigTemplateListResponse {
  templates: ConfigTemplate[];
  total: number;
  page: number;
  page_size: number;
}

export const configApi = {
  async listTemplates(adapterType?: string): Promise<ConfigTemplate[]> {
    const params = new URLSearchParams();
    if (adapterType) {
      params.append('adapter_type', adapterType);
    }
    const { data } = await client.get<ApiResponse<ConfigTemplateListResponse>>(`/configs?${params}`);
    return data.data?.templates || [];
  },

  async getTemplate(id: string): Promise<ConfigTemplate> {
    const { data } = await client.get<ApiResponse<ConfigTemplate>>(`/configs/${id}`);
    return data.data!;
  },

  async createTemplate(requestData: CreateTemplateRequest): Promise<ConfigTemplate> {
    const { data } = await client.post<ApiResponse<ConfigTemplate>>('/configs', requestData);
    return data.data!;
  },

  async updateTemplate(id: string, requestData: UpdateTemplateRequest): Promise<ConfigTemplate> {
    const { data } = await client.put<ApiResponse<ConfigTemplate>>(`/configs/${id}`, requestData);
    return data.data!;
  },

  async deleteTemplate(id: string): Promise<void> {
    await client.delete(`/configs/${id}`);
  },
};
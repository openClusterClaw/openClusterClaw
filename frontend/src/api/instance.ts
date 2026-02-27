import apiClient from './client';
import type {
  ClawInstance,
  CreateInstanceRequest,
  UpdateInstanceRequest,
  InstanceListResponse,
  ApiResponse,
} from '@/types';

export const instanceApi = {
  // Get instance list
  list: (params?: { tenant_id?: string; project_id?: string; page?: number; page_size?: number }) =>
    apiClient.get<ApiResponse<InstanceListResponse>>('/instances', { params }),

  // Get instance by ID
  get: (id: string) => apiClient.get<ApiResponse<ClawInstance>>(`/instances/${id}`),

  // Create instance
  create: (data: CreateInstanceRequest) => apiClient.post<ApiResponse<ClawInstance>>('/instances', data),

  // Update instance
  update: (id: string, data: UpdateInstanceRequest) => apiClient.put<ApiResponse<ClawInstance>>(`/instances/${id}`, data),

  // Delete instance
  delete: (id: string) => apiClient.delete<ApiResponse<void>>(`/instances/${id}`),

  // Start instance
  start: (id: string) => apiClient.post<ApiResponse<{ message: string }>>(`/instances/${id}/start`),

  // Stop instance
  stop: (id: string) => apiClient.post<ApiResponse<{ message: string }>>(`/instances/${id}/stop`),

  // Restart instance
  restart: (id: string) => apiClient.post<ApiResponse<{ message: string }>>(`/instances/${id}/restart`),

  // Get instance logs
  getLogs: (id: string, tailLines?: number) =>
    apiClient.get<ApiResponse<{ logs: string; tail_lines: number }>>(`/instances/${id}/logs`, {
      params: tailLines ? { tail_lines: tailLines } : undefined,
    }),

  // Helper method to get logs string directly
  getInstanceLogs: async (id: string, tailLines = 100): Promise<string> => {
    const response = await instanceApi.getLogs(id, tailLines);
    return response.data.data?.logs || '';
  },
};
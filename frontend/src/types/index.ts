export interface ApiResponse<T = any> {
  code: number;
  message: string;
  data?: T;
}

export interface ErrorResponse {
  code: number;
  message: string;
  error?: string;
}

// Instance Types
export type InstanceStatus = 'Creating' | 'Running' | 'Stopped' | 'Failed' | 'Destroyed';

export interface ResourceSpec {
  cpu: string;
  memory: string;
}

export interface StorageSpec {
  config_dir: string;
  data_dir: string;
  size: string;
}

export interface InstanceConfig {
  template_name: string;
  overrides: Record<string, string>;
}

export interface ClawInstance {
  id: string;
  name: string;
  tenant_id: string;
  project_id: string;
  type: string;
  version: string;
  status: InstanceStatus;
  config?: InstanceConfig;
  resources?: ResourceSpec;
  storage?: StorageSpec;
  created_at: string;
  updated_at: string;
}

export interface CreateInstanceRequest {
  name: string;
  tenant_id: string;
  project_id: string;
  type: string;
  version: string;
  config?: InstanceConfig;
  cpu?: string;
  memory?: string;
}

export interface UpdateInstanceRequest {
  name?: string;
  config?: InstanceConfig;
  cpu?: string;
  memory?: string;
}

export interface InstanceListResponse {
  instances: ClawInstance[];
  total: number;
  page: number;
  page_size: number;
}

// Config Template Types
export interface ConfigVariable {
  name: string;
  type: 'string' | 'number' | 'boolean';
  default: string;
  required: boolean;
  secret: boolean;
  description: string;
}

export interface ConfigTemplate {
  id: string;
  name: string;
  description: string;
  variables: ConfigVariable[];
  adapter_type: string;
  version: string;
  created_at: string;
  updated_at: string;
}

// Tenant Types
export interface Quota {
  max_instances: number;
  max_cpu: string;
  max_memory: string;
  max_storage: string;
}

export interface Tenant {
  id: string;
  name: string;
  quota?: Quota;
  created_at: string;
  updated_at: string;
}

// Project Types
export interface Project {
  id: string;
  tenant_id: string;
  name: string;
  created_at: string;
  updated_at: string;
}
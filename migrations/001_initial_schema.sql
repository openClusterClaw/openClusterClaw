-- Create tenants table
CREATE TABLE IF NOT EXISTS tenants (
    id TEXT PRIMARY KEY,
    name TEXT NOT NULL UNIQUE,
    max_instances INTEGER DEFAULT 10,
    max_cpu TEXT DEFAULT '10',
    max_memory TEXT DEFAULT '20Gi',
    max_storage TEXT DEFAULT '100Gi',
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    updated_at DATETIME DEFAULT CURRENT_TIMESTAMP
);

-- Create projects table
CREATE TABLE IF NOT EXISTS projects (
    id TEXT PRIMARY KEY,
    tenant_id TEXT NOT NULL REFERENCES tenants(id) ON DELETE CASCADE,
    name TEXT NOT NULL,
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    updated_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    UNIQUE(tenant_id, name)
);

-- Create config templates table
CREATE TABLE IF NOT EXISTS config_templates (
    id TEXT PRIMARY KEY,
    name TEXT NOT NULL UNIQUE,
    description TEXT,
    variables BLOB,
    adapter_type TEXT NOT NULL,
    version TEXT DEFAULT '1.0.0',
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    updated_at DATETIME DEFAULT CURRENT_TIMESTAMP
);

-- Create claw instances table
CREATE TABLE IF NOT EXISTS claw_instances (
    id TEXT PRIMARY KEY,
    name TEXT NOT NULL,
    tenant_id TEXT NOT NULL REFERENCES tenants(id) ON DELETE CASCADE,
    project_id TEXT REFERENCES projects(id) ON DELETE CASCADE,
    type TEXT NOT NULL,
    version TEXT NOT NULL,
    status TEXT NOT NULL DEFAULT 'Creating',
    config BLOB,
    cpu TEXT,
    memory TEXT,
    config_dir TEXT,
    data_dir TEXT,
    storage_size TEXT,
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    updated_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    UNIQUE(tenant_id, name)
);

-- Create indexes
CREATE INDEX IF NOT EXISTS idx_claw_instances_tenant ON claw_instances(tenant_id);
CREATE INDEX IF NOT EXISTS idx_claw_instances_project ON claw_instances(project_id);
CREATE INDEX IF NOT EXISTS idx_claw_instances_status ON claw_instances(status);
CREATE INDEX IF NOT EXISTS idx_projects_tenant ON projects(tenant_id);
package repository

import (
	"context"
	"database/sql"
	"github.com/google/uuid"
	"github.com/weibh/openClusterClaw/internal/model"
)

// InstanceRepository defines the interface for instance data access
type InstanceRepository interface {
	Create(ctx context.Context, instance *model.ClawInstance) error
	GetByID(ctx context.Context, id string) (*model.ClawInstance, error)
	List(ctx context.Context, tenantID, projectID string, limit, offset int) ([]*model.ClawInstance, error)
	Update(ctx context.Context, instance *model.ClawInstance) error
	UpdateStatus(ctx context.Context, id string, status model.InstanceStatus) error
	Delete(ctx context.Context, id string) error
}

// ConfigTemplateRepository defines the interface for config template data access
type ConfigTemplateRepository interface {
	Create(ctx context.Context, template *model.ConfigTemplate) error
	GetByID(ctx context.Context, id string) (*model.ConfigTemplate, error)
	List(ctx context.Context, limit, offset int) ([]*model.ConfigTemplate, error)
	Update(ctx context.Context, template *model.ConfigTemplate) error
	Delete(ctx context.Context, id string) error
}

// TenantRepository defines the interface for tenant data access
type TenantRepository interface {
	Create(ctx context.Context, tenant *model.Tenant) error
	GetByID(ctx context.Context, id string) (*model.Tenant, error)
	List(ctx context.Context, limit, offset int) ([]*model.Tenant, error)
	Update(ctx context.Context, tenant *model.Tenant) error
}

// ProjectRepository defines the interface for project data access
type ProjectRepository interface {
	Create(ctx context.Context, project *model.Project) error
	GetByID(ctx context.Context, id string) (*model.Project, error)
	ListByTenant(ctx context.Context, tenantID string, limit, offset int) ([]*model.Project, error)
	Update(ctx context.Context, project *model.Project) error
}

// DB defines the database interface
type DB interface {
	ExecContext(ctx context.Context, query string, args ...any) (sql.Result, error)
	QueryContext(ctx context.Context, query string, args ...any) (*sql.Rows, error)
	QueryRowContext(ctx context.Context, query string, args ...any) *sql.Row
	BeginTx(ctx context.Context, opts *sql.TxOptions) (*sql.Tx, error)
}

// instanceRepository implements InstanceRepository
type instanceRepository struct {
	db DB
}

// NewInstanceRepository creates a new instance repository
func NewInstanceRepository(db DB) InstanceRepository {
	return &instanceRepository{db: db}
}

func (r *instanceRepository) Create(ctx context.Context, instance *model.ClawInstance) error {
	if instance.ID == "" {
		instance.ID = uuid.New().String()
	}

	query := `
		INSERT INTO claw_instances (id, name, tenant_id, project_id, type, version, status, config, cpu, memory, config_dir, data_dir, storage_size)
		VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)
	`
	_, err := r.db.ExecContext(ctx, query,
		instance.ID, instance.Name, instance.TenantID, instance.ProjectID,
		instance.Type, instance.Version, string(instance.Status), instance.Config,
		instance.CPU, instance.Memory, instance.ConfigDir, instance.DataDir, instance.StorageSize,
	)
	return err
}

func (r *instanceRepository) GetByID(ctx context.Context, id string) (*model.ClawInstance, error) {
	query := `SELECT id, name, tenant_id, project_id, type, version, status, config, cpu, memory, config_dir, data_dir, storage_size, created_at, updated_at FROM claw_instances WHERE id = ?`
	row := r.db.QueryRowContext(ctx, query, id)

	instance := &model.ClawInstance{}
	var status string
	err := row.Scan(
		&instance.ID, &instance.Name, &instance.TenantID, &instance.ProjectID,
		&instance.Type, &instance.Version, &status, &instance.Config,
		&instance.CPU, &instance.Memory, &instance.ConfigDir, &instance.DataDir, &instance.StorageSize,
		&instance.CreatedAt, &instance.UpdatedAt,
	)
	if err != nil {
		return nil, err
	}
	instance.Status = model.InstanceStatus(status)
	return instance, nil
}

func (r *instanceRepository) List(ctx context.Context, tenantID, projectID string, limit, offset int) ([]*model.ClawInstance, error) {
	query := `
		SELECT id, name, tenant_id, project_id, type, version, status, config, cpu, memory, config_dir, data_dir, storage_size, created_at, updated_at
		FROM claw_instances
		WHERE (? = '' OR tenant_id = ?)
		AND (? = '' OR project_id = ?)
		ORDER BY created_at DESC
		LIMIT ? OFFSET ?
	`
	rows, err := r.db.QueryContext(ctx, query, tenantID, tenantID, projectID, projectID, limit, offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var instances []*model.ClawInstance
	for rows.Next() {
		instance := &model.ClawInstance{}
		var status string
		err := rows.Scan(
			&instance.ID, &instance.Name, &instance.TenantID, &instance.ProjectID,
			&instance.Type, &instance.Version, &status, &instance.Config,
			&instance.CPU, &instance.Memory, &instance.ConfigDir, &instance.DataDir, &instance.StorageSize,
			&instance.CreatedAt, &instance.UpdatedAt,
		)
		if err != nil {
			return nil, err
		}
		instance.Status = model.InstanceStatus(status)
		instances = append(instances, instance)
	}

	return instances, nil
}

func (r *instanceRepository) Update(ctx context.Context, instance *model.ClawInstance) error {
	query := `
		UPDATE claw_instances
		SET name = ?, type = ?, version = ?, status = ?, config = ?, cpu = ?, memory = ?, config_dir = ?, data_dir = ?, storage_size = ?, updated_at = datetime('now')
		WHERE id = ?
	`
	_, err := r.db.ExecContext(ctx, query,
		instance.Name, instance.Type, instance.Version, string(instance.Status),
		instance.Config, instance.CPU, instance.Memory, instance.ConfigDir, instance.DataDir, instance.StorageSize,
		instance.ID,
	)
	return err
}

func (r *instanceRepository) UpdateStatus(ctx context.Context, id string, status model.InstanceStatus) error {
	query := `UPDATE claw_instances SET status = ?, updated_at = datetime('now') WHERE id = ?`
	_, err := r.db.ExecContext(ctx, query, string(status), id)
	return err
}

func (r *instanceRepository) Delete(ctx context.Context, id string) error {
	query := `DELETE FROM claw_instances WHERE id = ?`
	_, err := r.db.ExecContext(ctx, query, id)
	return err
}
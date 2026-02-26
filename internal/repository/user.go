package repository

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/google/uuid"
	"github.com/weibh/openClusterClaw/internal/model"
)

// UserRepository handles user data persistence
type UserRepository struct {
	db *sql.DB
}

// NewUserRepository creates a new user repository
func NewUserRepository(db *sql.DB) *UserRepository {
	return &UserRepository{db: db}
}

// Create creates a new user
func (r *UserRepository) Create(ctx context.Context, user *model.User) error {
	query := `
		INSERT INTO users (id, username, password_hash, tenant_id, role, is_active)
		VALUES (?, ?, ?, ?, ?, ?)
	`
	_, err := r.db.ExecContext(ctx, query, user.ID, user.Username, user.PasswordHash,
		user.TenantID, user.Role, user.IsActive)
	if err != nil {
		return fmt.Errorf("failed to create user: %w", err)
	}
	return nil
}

// GetByID retrieves a user by ID
func (r *UserRepository) GetByID(ctx context.Context, id string) (*model.User, error) {
	query := `
		SELECT id, username, password_hash, tenant_id, role, is_active, created_at, updated_at
		FROM users WHERE id = ?
	`
	var user model.User
	err := r.db.QueryRowContext(ctx, query, id).Scan(
		&user.ID, &user.Username, &user.PasswordHash, &user.TenantID,
		&user.Role, &user.IsActive, &user.CreatedAt, &user.UpdatedAt,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("user not found")
		}
		return nil, fmt.Errorf("failed to get user: %w", err)
	}
	return &user, nil
}

// GetByUsername retrieves a user by username
func (r *UserRepository) GetByUsername(ctx context.Context, username string) (*model.User, error) {
	query := `
		SELECT id, username, password_hash, tenant_id, role, is_active, created_at, updated_at
		FROM users WHERE username = ?
	`
	var user model.User
	err := r.db.QueryRowContext(ctx, query, username).Scan(
		&user.ID, &user.Username, &user.PasswordHash, &user.TenantID,
		&user.Role, &user.IsActive, &user.CreatedAt, &user.UpdatedAt,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("user not found")
		}
		return nil, fmt.Errorf("failed to get user: %w", err)
	}
	return &user, nil
}

// Update updates a user
func (r *UserRepository) Update(ctx context.Context, user *model.User) error {
	query := `
		UPDATE users
		SET username = ?, password_hash = ?, tenant_id = ?, role = ?, is_active = ?, updated_at = CURRENT_TIMESTAMP
		WHERE id = ?
	`
	_, err := r.db.ExecContext(ctx, query, user.Username, user.PasswordHash,
		user.TenantID, user.Role, user.IsActive, user.ID)
	if err != nil {
		return fmt.Errorf("failed to update user: %w", err)
	}
	return nil
}

// Delete deletes a user by ID
func (r *UserRepository) Delete(ctx context.Context, id string) error {
	query := `DELETE FROM users WHERE id = ?`
	result, err := r.db.ExecContext(ctx, query, id)
	if err != nil {
		return fmt.Errorf("failed to delete user: %w", err)
	}
	rowsAffected, _ := result.RowsAffected()
	if rowsAffected == 0 {
		return fmt.Errorf("user not found")
	}
	return nil
}

// List retrieves a list of users, optionally filtered by tenant
func (r *UserRepository) List(ctx context.Context, tenantID string, limit, offset int) ([]*model.User, int, error) {
	var query string
	var args []interface{}

	if tenantID != "" {
		query = `
			SELECT id, username, password_hash, tenant_id, role, is_active, created_at, updated_at
			FROM users WHERE tenant_id = ?
			ORDER BY created_at DESC LIMIT ? OFFSET ?
		`
		args = []interface{}{tenantID, limit, offset}
	} else {
		query = `
			SELECT id, username, password_hash, tenant_id, role, is_active, created_at, updated_at
			FROM users
			ORDER BY created_at DESC LIMIT ? OFFSET ?
		`
		args = []interface{}{limit, offset}
	}

	rows, err := r.db.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to list users: %w", err)
	}
	defer rows.Close()

	var users []*model.User
	for rows.Next() {
		var user model.User
		if err := rows.Scan(&user.ID, &user.Username, &user.PasswordHash, &user.TenantID,
			&user.Role, &user.IsActive, &user.CreatedAt, &user.UpdatedAt); err != nil {
			return nil, 0, fmt.Errorf("failed to scan user: %w", err)
		}
		users = append(users, &user)
	}

	// Get total count
	var countQuery string
	var countArgs []interface{}
	if tenantID != "" {
		countQuery = `SELECT COUNT(*) FROM users WHERE tenant_id = ?`
		countArgs = []interface{}{tenantID}
	} else {
		countQuery = `SELECT COUNT(*) FROM users`
		countArgs = []interface{}{}
	}

	var total int
	if err := r.db.QueryRowContext(ctx, countQuery, countArgs...).Scan(&total); err != nil {
		return nil, 0, fmt.Errorf("failed to count users: %w", err)
	}

	return users, total, nil
}

// GenerateID generates a new UUID for a user
func (r *UserRepository) GenerateID() string {
	return uuid.New().String()
}

package main

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"path/filepath"
	"syscall"
	"time"

	"github.com/weibh/openClusterClaw/config"
	"github.com/weibh/openClusterClaw/internal/api"
	"github.com/weibh/openClusterClaw/internal/pkg/jwt"
	"github.com/weibh/openClusterClaw/internal/repository"
	"github.com/weibh/openClusterClaw/internal/runtime/k8s"
	"github.com/weibh/openClusterClaw/internal/service"
	_ "modernc.org/sqlite"
)

func main() {
	// Load configuration
	cfg, err := config.Load("./config/config.yaml")
	if err != nil {
		log.Fatalf("Failed to load config: %v", err)
	}

	// Set gin mode
	ginMode := cfg.Server.Mode
	if ginMode == "" {
		ginMode = "debug"
	}

	// Initialize database connection
	db, err := initDB(cfg)
	if err != nil {
		log.Fatalf("Failed to initialize database: %v", err)
	}
	defer db.Close()

	// Initialize repositories
	instanceRepo := repository.NewInstanceRepository(db)
	userRepo := repository.NewUserRepository(db)

	// Initialize K8S client
	var podManager *k8s.PodManager
	var configMapManager *k8s.ConfigMapManager
	if cfg.K8S.Enabled {
		kubeconfig := cfg.K8S.Kubeconfig
		if kubeconfig == "" {
			// Try environment variable or default path
			kubeconfig = os.Getenv("KUBECONFIG")
		}

		err := k8s.Initialize(kubeconfig)
		if err != nil {
			log.Printf("Warning: Failed to initialize K8S client: %v", err)
			log.Println("Continuing without K8S integration...")
		} else {
			namespace := cfg.K8S.Namespace
			if namespace == "" {
				namespace = "default"
			}
			podManager = k8s.NewPodManager(namespace)
			configMapManager = k8s.NewConfigMapManager(namespace)
			log.Printf("K8S client initialized, using namespace: %s", namespace)
		}
	}

	// Initialize services
	instanceService := service.NewInstanceService(instanceRepo, podManager, configMapManager)

	// Initialize JWT service
	jwtService := jwt.NewJWTService(cfg)

	// Initialize auth service
	authService := service.NewAuthService(userRepo, jwtService)

	// Initialize default tenant and admin user if they don't exist
	if err := initializeDefaultData(db, authService, userRepo); err != nil {
		log.Printf("Warning: Failed to initialize default data: %v", err)
	}

	// Initialize router
	router := api.NewRouter(instanceService, authService, jwtService)
	router.SetupRoutes()
	engine := router.Engine()

	// Create HTTP server
	srv := &http.Server{
		Addr:         fmt.Sprintf(":%d", cfg.Server.Port),
		Handler:      engine,
		ReadTimeout:  time.Duration(cfg.Server.ReadTimeout) * time.Second,
		WriteTimeout: time.Duration(cfg.Server.WriteTimeout) * time.Second,
	}

	// Start server in a goroutine
	go func() {
		log.Printf("Server starting on port %d (mode: %s)", cfg.Server.Port, ginMode)
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("Failed to start server: %v", err)
		}
	}()

	// Wait for interrupt signal
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	log.Println("Shutting down server...")

	// Graceful shutdown with timeout
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := srv.Shutdown(ctx); err != nil {
		log.Fatalf("Server forced to shutdown: %v", err)
	}

	log.Println("Server exited")
}

// initDB initializes the database connection and creates tables
func initDB(cfg *config.Config) (*sql.DB, error) {
	// Ensure data directory exists
	dataDir := filepath.Dir(cfg.Database.Path)
	if dataDir != "." && dataDir != "" {
		if err := os.MkdirAll(dataDir, 0755); err != nil {
			return nil, fmt.Errorf("failed to create data directory: %w", err)
		}
	}

	// Open database connection
	db, err := sql.Open("sqlite", cfg.Database.Path+"?_pragma=foreign_keys(1)")
	if err != nil {
		return nil, fmt.Errorf("failed to open database: %w", err)
	}

	// Test connection
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := db.PingContext(ctx); err != nil {
		return nil, fmt.Errorf("failed to ping database: %w", err)
	}

	// Run migrations
	if err := runMigrations(db); err != nil {
		return nil, fmt.Errorf("failed to run migrations: %w", err)
	}

	log.Println("Database connection established")

	return db, nil
}

// runMigrations creates the database tables
func runMigrations(db *sql.DB) error {
	migrations := []string{
		`CREATE TABLE IF NOT EXISTS tenants (
			id TEXT PRIMARY KEY,
			name TEXT NOT NULL UNIQUE,
			max_instances INTEGER DEFAULT 10,
			max_cpu TEXT DEFAULT '10',
			max_memory TEXT DEFAULT '20Gi',
			max_storage TEXT DEFAULT '100Gi',
			created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
			updated_at DATETIME DEFAULT CURRENT_TIMESTAMP
		)`,
		`CREATE TABLE IF NOT EXISTS projects (
			id TEXT PRIMARY KEY,
			tenant_id TEXT NOT NULL REFERENCES tenants(id) ON DELETE CASCADE,
			name TEXT NOT NULL,
			created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
			updated_at DATETIME DEFAULT CURRENT_TIMESTAMP,
			UNIQUE(tenant_id, name)
		)`,
		`CREATE TABLE IF NOT EXISTS config_templates (
			id TEXT PRIMARY KEY,
			name TEXT NOT NULL UNIQUE,
			description TEXT,
			variables BLOB,
			adapter_type TEXT NOT NULL,
			version TEXT DEFAULT '1.0.0',
			created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
			updated_at DATETIME DEFAULT CURRENT_TIMESTAMP
		)`,
		`CREATE TABLE IF NOT EXISTS claw_instances (
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
		)`,
		`CREATE INDEX IF NOT EXISTS idx_claw_instances_tenant ON claw_instances(tenant_id)`,
		`CREATE INDEX IF NOT EXISTS idx_claw_instances_project ON claw_instances(project_id)`,
		`CREATE INDEX IF NOT EXISTS idx_claw_instances_status ON claw_instances(status)`,
		`CREATE INDEX IF NOT EXISTS idx_projects_tenant ON projects(tenant_id)`,
		`CREATE TABLE IF NOT EXISTS users (
			id TEXT PRIMARY KEY,
			username TEXT NOT NULL UNIQUE,
			password_hash TEXT NOT NULL,
			tenant_id TEXT REFERENCES tenants(id) ON DELETE CASCADE,
			role TEXT DEFAULT 'user',
			is_active BOOLEAN DEFAULT true,
			created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
			updated_at DATETIME DEFAULT CURRENT_TIMESTAMP
		)`,
		`CREATE INDEX IF NOT EXISTS idx_users_tenant ON users(tenant_id)`,
		`CREATE INDEX IF NOT EXISTS idx_users_username ON users(username)`,
	}

	for _, migration := range migrations {
		if _, err := db.Exec(migration); err != nil {
			return fmt.Errorf("migration failed: %w", err)
		}
	}

	return nil
}

// initializeDefaultData creates default tenant and admin user if they don't exist
func initializeDefaultData(db *sql.DB, authService *service.AuthService, userRepo *repository.UserRepository) error {
	ctx := context.Background()

	// Check if default admin user exists
	_, err := userRepo.GetByUsername(ctx, "admin")
	if err == nil {
		// Admin user already exists
		return nil
	}

	// Create default tenant
	tenantID := "default-tenant"
	_, err = db.Exec(`
		INSERT OR IGNORE INTO tenants (id, name, max_instances, max_cpu, max_memory, max_storage)
		VALUES (?, ?, ?, ?, ?, ?)
	`, tenantID, "Default Tenant", 100, "100", "200Gi", "1Ti")
	if err != nil {
		return fmt.Errorf("failed to create default tenant: %w", err)
	}

	// Create default admin user
	defaultAdmin := &service.CreateUserRequest{
		Username: "admin",
		Password: "admin123", // Default password, should be changed in production
		TenantID: tenantID,
		Role:     "admin",
	}

	_, err = authService.CreateUser(ctx, defaultAdmin)
	if err != nil {
		return fmt.Errorf("failed to create default admin user: %w", err)
	}

	log.Println("Default admin user created: username=admin, password=admin123")
	log.Println("Please change the default password after first login!")

	return nil
}
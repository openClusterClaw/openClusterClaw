package main

import (
	"context"
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
	"github.com/weibh/openClusterClaw/internal/model"
	"github.com/weibh/openClusterClaw/internal/pkg/jwt"
	"github.com/weibh/openClusterClaw/internal/repository"
	"github.com/weibh/openClusterClaw/internal/runtime/k8s"
	"github.com/weibh/openClusterClaw/internal/service"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func main() {
	// Load configuration - try multiple paths
	cfgPaths := []string{
		"../../config/config.yaml",
		"../config/config.yaml",
		"./config/config.yaml",
	}
	var cfg *config.Config
	var err error
	for _, path := range cfgPaths {
		cfg, err = config.Load(path)
		if err == nil {
			break
		}
	}
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

	// Initialize repositories
	instanceRepo := repository.NewInstanceRepository(db)
	userRepo := repository.NewUserRepository(db)
	configTemplateRepo := repository.NewConfigTemplateRepository(db)
	tenantRepo := repository.NewTenantRepository(db)
	projectRepo := repository.NewProjectRepository(db)

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
	configTemplateService := service.NewConfigTemplateService(configTemplateRepo)
	tenantService := service.NewTenantService(tenantRepo, instanceRepo)
	projectService := service.NewProjectService(projectRepo)

	// Initialize JWT service
	jwtService := jwt.NewJWTService(cfg)

	// Initialize auth service
	authService := service.NewAuthService(userRepo, jwtService)

	// Initialize default tenant and admin user if they don't exist
	if err := initializeDefaultData(db, authService, userRepo); err != nil {
		log.Printf("Warning: Failed to initialize default data: %v", err)
	}

	// Initialize router
	router := api.NewRouter(instanceService, configTemplateService, tenantService, projectService, authService, jwtService, userRepo, cfg)
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
func initDB(cfg *config.Config) (*gorm.DB, error) {
	// Ensure data directory exists
	dataDir := filepath.Dir(cfg.Database.Path)
	if dataDir != "." && dataDir != "" {
		if err := os.MkdirAll(dataDir, 0755); err != nil {
			return nil, fmt.Errorf("failed to create data directory: %w", err)
		}
	}

	// Open database connection with foreign key support
	dsn := cfg.Database.Path + "?_pragma=foreign_keys(1)"
	db, err := gorm.Open(sqlite.Open(dsn), &gorm.Config{})
	if err != nil {
		return nil, fmt.Errorf("failed to open database: %w", err)
	}

	// Get underlying sql.DB for configuration
	sqlDB, err := db.DB()
	if err != nil {
		return nil, fmt.Errorf("failed to get database instance: %w", err)
	}

	// Test connection
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := sqlDB.PingContext(ctx); err != nil {
		return nil, fmt.Errorf("failed to ping database: %w", err)
	}

	// Run migrations
	if err := runMigrations(db); err != nil {
		return nil, fmt.Errorf("failed to run migrations: %w", err)
	}

	log.Println("Database connection established")

	return db, nil
}

// runMigrations creates the database tables using GORM AutoMigrate
func runMigrations(db *gorm.DB) error {
	return db.AutoMigrate(
		&model.Tenant{},
		&model.Project{},
		&model.ConfigTemplate{},
		&model.ClawInstance{},
		&model.User{},
	)
}

// initializeDefaultData creates default tenant and admin user if they don't exist
func initializeDefaultData(db *gorm.DB, authService *service.AuthService, userRepo *repository.UserRepository) error {
	ctx := context.Background()

	// Check if default admin user exists
	_, err := userRepo.GetByUsername(ctx, "admin")
	if err == nil {
		// Admin user already exists
		return nil
	}

	// Create default tenant
	tenantID := "default-tenant"
	tenant := &model.Tenant{
		ID:           tenantID,
		Name:         "Default Tenant",
		MaxInstances: 100,
		MaxCPU:       "100",
		MaxMemory:    "200Gi",
		MaxStorage:   "1Ti",
	}
	if err := db.Create(tenant).Error; err != nil {
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
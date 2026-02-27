package api

import (
	"encoding/hex"
	"github.com/gin-gonic/gin"
	"github.com/weibh/openClusterClaw/internal/embed"
	"github.com/weibh/openClusterClaw/internal/middleware"
	"github.com/weibh/openClusterClaw/internal/pkg/jwt"
	"github.com/weibh/openClusterClaw/internal/pkg/otp"
	"github.com/weibh/openClusterClaw/internal/service"
	"github.com/weibh/openClusterClaw/internal/repository"
	"github.com/weibh/openClusterClaw/config"
	"log"
)

// Router sets up API routes
type Router struct {
	handler         *Handler
	authHandler     *AuthHandler
	otpHandler      *OTPHandler
	configHandler   *ConfigTemplateHandler
	tenantHandler   *TenantHandler
	projectHandler  *ProjectHandler
	engine          *gin.Engine
	jwtService      *jwt.JWTService
}

// NewRouter creates a new router
func NewRouter(
	instanceService service.InstanceService,
	configTemplateService service.ConfigTemplateService,
	tenantService service.TenantService,
	projectService service.ProjectService,
	authService *service.AuthService,
	jwtService *jwt.JWTService,
	userRepo *repository.UserRepository,
	cfg *config.Config,
) *Router {
	handler := NewHandler(instanceService)
	authHandler := NewAuthHandler(authService)
	configHandler := NewConfigTemplateHandler(configTemplateService)
	tenantHandler := NewTenantHandler(tenantService)
	projectHandler := NewProjectHandler(projectService)
	engine := gin.Default()

	// Create OTP service from config
	encryptionKeyBytes, err := hex.DecodeString(cfg.OTP.EncryptionKey)
	if err != nil {
		log.Fatalf("Failed to decode encryption key: %v", err)
	}
	if len(encryptionKeyBytes) != 32 {
		log.Fatalf("Encryption key must be 32 bytes, got %d", len(encryptionKeyBytes))
	}
	otpService := otp.NewService(string(encryptionKeyBytes), cfg.OTP.Issuer)
	otpSvc := service.NewOTPService(userRepo, otpService, jwtService)
	otpHandler := NewOTPHandler(otpSvc, authService, userRepo)

	return &Router{
		handler:         handler,
		authHandler:     authHandler,
		otpHandler:      otpHandler,
		configHandler:   configHandler,
		tenantHandler:   tenantHandler,
		projectHandler:  projectHandler,
		engine:          engine,
		jwtService:      jwtService,
	}
}

// SetupRoutes configures all routes
func (r *Router) SetupRoutes() {
	api := r.engine.Group("/api/v1")
	{
		// Auth routes (public)
		auth := api.Group("/auth")
		{
			auth.POST("/login", r.otpHandler.LoginWithOTP)
			auth.POST("/refresh", r.authHandler.RefreshToken)
			auth.POST("/otp/verify", r.otpHandler.VerifyOTP)
		}

		// Authenticated routes
		authenticated := api.Group("")
		authenticated.Use(middleware.AuthMiddleware(r.jwtService))
		{
			// User routes
			auth := authenticated.Group("/auth")
			{
				auth.POST("/logout", r.authHandler.Logout)
				auth.GET("/me", r.authHandler.GetCurrentUser)
			}

			// OTP routes
			otp := authenticated.Group("/auth/otp")
			{
				otp.POST("/generate", r.otpHandler.GenerateSecret)
				otp.POST("/enable", r.otpHandler.EnableOTP)
				otp.POST("/disable", r.otpHandler.DisableOTP)
				otp.GET("/backup", r.otpHandler.GetBackupCodes)
				otp.GET("/status", r.otpHandler.GetOTPStatus)
			}

			// User management (admin only)
			users := authenticated.Group("/auth/users")
			users.Use(middleware.RequireAdmin())
			{
				users.POST("", r.authHandler.CreateUser)
			}

			// Config template routes (admin only)
			configs := authenticated.Group("/configs")
			configs.Use(middleware.RequireAdmin())
			{
				configs.POST("", r.configHandler.Create)
				configs.GET("", r.configHandler.List)
				configs.GET("/:id", r.configHandler.Get)
				configs.PUT("/:id", r.configHandler.Update)
				configs.DELETE("/:id", r.configHandler.Delete)
			}

			// Tenant routes (admin only)
			tenants := authenticated.Group("/tenants")
			tenants.Use(middleware.RequireAdmin())
			{
				tenants.POST("", r.tenantHandler.Create)
				tenants.GET("", r.tenantHandler.List)
				tenants.GET("/:id", r.tenantHandler.Get)
				tenants.PUT("/:id", r.tenantHandler.Update)
				tenants.DELETE("/:id", r.tenantHandler.Delete)
			}

			// Project routes (admin only)
			projects := authenticated.Group("/projects")
			projects.Use(middleware.RequireAdmin())
			{
				projects.POST("", r.projectHandler.Create)
				projects.GET("", r.projectHandler.List)
				projects.GET("/:id", r.projectHandler.Get)
				projects.PUT("/:id", r.projectHandler.Update)
				projects.DELETE("/:id", r.projectHandler.Delete)
			}

			// Instance routes (need authentication)
			instanceHandler := NewInstanceHandler(r.handler.instanceService)
			instances := authenticated.Group("/instances")
			{
				instances.POST("", instanceHandler.Create)
				instances.GET("", instanceHandler.List)
				instances.GET("/:id", instanceHandler.Get)
				instances.DELETE("/:id", instanceHandler.Delete)
				instances.POST("/:id/start", instanceHandler.Start)
				instances.POST("/:id/stop", instanceHandler.Stop)
				instances.POST("/:id/restart", instanceHandler.Restart)
				instances.GET("/:id/logs", instanceHandler.Logs)
			}
		}
	}

	// Health check
	r.engine.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{"status": "ok"})
	})

	// Setup embedded frontend static file serving
	embed.SetupRouter(r.engine)
}

// Engine returns gin engine
func (r *Router) Engine() *gin.Engine {
	return r.engine
}

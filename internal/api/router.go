package api

import (
	"github.com/gin-gonic/gin"
	"github.com/weibh/openClusterClaw/internal/embed"
	"github.com/weibh/openClusterClaw/internal/middleware"
	"github.com/weibh/openClusterClaw/internal/pkg/jwt"
	"github.com/weibh/openClusterClaw/internal/service"
)

// Router sets up the API routes
type Router struct {
	handler      *Handler
	authHandler  *AuthHandler
	engine       *gin.Engine
	jwtService   *jwt.JWTService
}

// NewRouter creates a new router
func NewRouter(instanceService service.InstanceService, authService *service.AuthService, jwtService *jwt.JWTService) *Router {
	handler := NewHandler(instanceService)
	authHandler := NewAuthHandler(authService)
	engine := gin.Default()

	return &Router{
		handler:      handler,
		authHandler:  authHandler,
		engine:       engine,
		jwtService:   jwtService,
	}
}

// SetupRoutes configures all routes
func (r *Router) SetupRoutes() {
	api := r.engine.Group("/api/v1")
	{
		// Auth routes (public)
		auth := api.Group("/auth")
		{
			auth.POST("/login", r.authHandler.Login)
			auth.POST("/refresh", r.authHandler.RefreshToken)
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

			// User management (admin only)
			users := authenticated.Group("/auth/users")
			users.Use(middleware.RequireAdmin())
			{
				users.POST("", r.authHandler.CreateUser)
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

// Engine returns the gin engine
func (r *Router) Engine() *gin.Engine {
	return r.engine
}

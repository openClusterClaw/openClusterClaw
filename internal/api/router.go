package api

import (
	"github.com/gin-gonic/gin"
	"github.com/weibh/openClusterClaw/internal/embed"
	"github.com/weibh/openClusterClaw/internal/service"
)

// Router sets up the API routes
type Router struct {
	handler *Handler
	engine  *gin.Engine
}

// NewRouter creates a new router
func NewRouter(instanceService service.InstanceService) *Router {
	handler := NewHandler(instanceService)
	engine := gin.Default()

	return &Router{
		handler: handler,
		engine:  engine,
	}
}

// SetupRoutes configures all routes
func (r *Router) SetupRoutes() {
	api := r.engine.Group("/api/v1")
	{
		instanceHandler := NewInstanceHandler(r.handler.instanceService)
		instances := api.Group("/instances")
		{
			instances.POST("", instanceHandler.Create)
			instances.GET("", instanceHandler.List)
			instances.GET("/:id", instanceHandler.Get)
			instances.DELETE("/:id", instanceHandler.Delete)
			instances.POST("/:id/start", instanceHandler.Start)
			instances.POST("/:id/stop", instanceHandler.Stop)
			instances.POST("/:id/restart", instanceHandler.Restart)
		}

		// TODO: Add other routes (configs, tenants, etc.)
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
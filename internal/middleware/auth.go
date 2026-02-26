package middleware

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/weibh/openClusterClaw/internal/pkg/jwt"
	"github.com/weibh/openClusterClaw/internal/model"
)

const (
	// ContextUserIDKey is the key for user ID in gin context
	ContextUserIDKey = "user_id"
	// ContextUsernameKey is the key for username in gin context
	ContextUsernameKey = "username"
	// ContextTenantIDKey is the key for tenant ID in gin context
	ContextTenantIDKey = "tenant_id"
	// ContextUserRoleKey is the key for user role in gin context
	ContextUserRoleKey = "role"
)

// AuthMiddleware validates JWT tokens and injects user info into context
func AuthMiddleware(jwtService *jwt.JWTService) gin.HandlerFunc {
	return func(c *gin.Context) {
		// Get Authorization header
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			c.JSON(http.StatusUnauthorized, gin.H{
				"code":    401,
				"message": "missing authorization header",
			})
			c.Abort()
			return
		}

		// Check Bearer format
		parts := strings.Split(authHeader, " ")
		if len(parts) != 2 || parts[0] != "Bearer" {
			c.JSON(http.StatusUnauthorized, gin.H{
				"code":    401,
				"message": "invalid authorization header format",
			})
			c.Abort()
			return
		}

		// Validate token
		token := parts[1]
		claims, err := jwtService.ValidateToken(token)
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{
				"code":    401,
				"message": fmt.Sprintf("invalid token: %v", err),
			})
			c.Abort()
			return
		}

		// Inject user info into context
		c.Set(ContextUserIDKey, claims.UserID)
		c.Set(ContextUsernameKey, claims.Username)
		c.Set(ContextTenantIDKey, claims.TenantID)
		c.Set(ContextUserRoleKey, claims.Role)

		c.Next()
	}
}

// RequireAdmin requires user to have admin role
func RequireAdmin() gin.HandlerFunc {
	return func(c *gin.Context) {
		role, exists := c.Get(ContextUserRoleKey)
		if !exists {
			c.JSON(http.StatusUnauthorized, gin.H{
				"code":    401,
				"message": "unauthorized",
			})
			c.Abort()
			return
		}

		if role != model.RoleAdmin {
			c.JSON(http.StatusForbidden, gin.H{
				"code":    403,
				"message": "admin role required",
			})
			c.Abort()
			return
		}

		c.Next()
	}
}

// GetUserID retrieves user ID from context
func GetUserID(c *gin.Context) string {
	userID, _ := c.Get(ContextUserIDKey)
	if userID == nil {
		return ""
	}
	return userID.(string)
}

// GetTenantID retrieves tenant ID from context
func GetTenantID(c *gin.Context) string {
	tenantID, _ := c.Get(ContextTenantIDKey)
	if tenantID == nil {
		return ""
	}
	return tenantID.(string)
}

// GetUsername retrieves username from context
func GetUsername(c *gin.Context) string {
	username, _ := c.Get(ContextUsernameKey)
	if username == nil {
		return ""
	}
	return username.(string)
}

// GetUserRole retrieves user role from context
func GetUserRole(c *gin.Context) string {
	role, _ := c.Get(ContextUserRoleKey)
	if role == nil {
		return ""
	}
	return role.(string)
}

// OptionalAuthMiddleware authenticates if token is provided, but doesn't require it
func OptionalAuthMiddleware(jwtService *jwt.JWTService) gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			// No token provided, continue without auth
			c.Next()
			return
		}

		parts := strings.Split(authHeader, " ")
		if len(parts) != 2 || parts[0] != "Bearer" {
			c.Next()
			return
		}

		token := parts[1]
		claims, err := jwtService.ValidateToken(token)
		if err != nil {
			// Invalid token, continue without auth
			c.Next()
			return
		}

		// Inject user info into context
		c.Set(ContextUserIDKey, claims.UserID)
		c.Set(ContextUsernameKey, claims.Username)
		c.Set(ContextTenantIDKey, claims.TenantID)
		c.Set(ContextUserRoleKey, claims.Role)

		c.Next()
	}
}

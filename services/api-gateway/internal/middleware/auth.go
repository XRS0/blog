package middleware

import (
	"context"
	"log/slog"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"

	authpb "github.com/XRS0/blog/services/api-gateway/proto/auth"
)

func AuthMiddleware(authClient authpb.AuthServiceClient, logger *slog.Logger, required bool) gin.HandlerFunc {
	return func(c *gin.Context) {
		// Extract token from Authorization header
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			if required {
				c.JSON(http.StatusUnauthorized, gin.H{"error": "authorization required"})
				c.Abort()
				return
			}
			c.Set("user_id", uint64(0))
			c.Next()
			return
		}

		// Remove "Bearer " prefix
		token := strings.TrimPrefix(authHeader, "Bearer ")
		if token == authHeader {
			if required {
				c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid authorization format"})
				c.Abort()
				return
			}
			c.Set("user_id", uint64(0))
			c.Next()
			return
		}

		// Validate token with auth service
		resp, err := authClient.ValidateToken(context.Background(), &authpb.ValidateTokenRequest{
			Token: token,
		})
		if err != nil || !resp.Valid {
			logger.Error("token validation failed", "error", err)
			if required {
				c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid token"})
				c.Abort()
				return
			}
			c.Set("user_id", uint64(0))
			c.Next()
			return
		}

		// Set user ID in context
		c.Set("user_id", resp.UserId)
		c.Next()
	}
}

func RequireAuth(authClient authpb.AuthServiceClient, logger *slog.Logger) gin.HandlerFunc {
	return AuthMiddleware(authClient, logger, true)
}

func OptionalAuth(authClient authpb.AuthServiceClient, logger *slog.Logger) gin.HandlerFunc {
	return AuthMiddleware(authClient, logger, false)
}

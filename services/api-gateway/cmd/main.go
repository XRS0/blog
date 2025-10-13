package main

import (
	"log"
	"log/slog"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/gin-gonic/gin"

	"github.com/XRS0/blog/services/api-gateway/internal/client"
	"github.com/XRS0/blog/services/api-gateway/internal/handlers"
	"github.com/XRS0/blog/services/api-gateway/internal/middleware"
	sharedLogger "github.com/XRS0/blog/shared/logger"
)

func main() {
	gin.SetMode(gin.ReleaseMode)

	// Logger setup
	logLevel := parseLogLevel(getEnv("LOG_LEVEL", "info"))
	logger := sharedLogger.New(logLevel, "api-gateway")

	// Connect to microservices
	authURL := getEnv("AUTH_SERVICE_URL", "localhost:50051")
	articleURL := getEnv("ARTICLE_SERVICE_URL", "localhost:50052")
	statsURL := getEnv("STATS_SERVICE_URL", "localhost:50053")

	clients, err := client.NewServiceClients(authURL, articleURL, statsURL)
	if err != nil {
		log.Fatalf("failed to connect to services: %v", err)
	}
	logger.Logger.Info("connected to microservices")

	// Create handlers
	authHandler := handlers.NewAuthHandler(clients.Auth, logger.Logger)
	articleHandler := handlers.NewArticleHandler(clients.Article, clients.Stats, logger.Logger)

	// Setup router
	router := gin.Default()
	router.Use(corsMiddleware())

	// Health check
	router.GET("/healthz", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"status": "ok", "timestamp": time.Now()})
	})

	// API routes
	api := router.Group("/api")
	{
		// Auth routes
		auth := api.Group("/auth")
		{
			auth.POST("/register", authHandler.Register)
			auth.POST("/login", authHandler.Login)
			auth.GET("/me", middleware.RequireAuth(clients.Auth, logger.Logger), authHandler.GetMe)
		}

		// Article routes
		articles := api.Group("/articles")
		{
			// Public routes (with optional auth)
			articles.GET("", middleware.OptionalAuth(clients.Auth, logger.Logger), articleHandler.ListArticles)
			articles.GET("/:id", middleware.OptionalAuth(clients.Auth, logger.Logger), articleHandler.GetArticle)

			// Protected routes (require auth)
			articles.POST("", middleware.RequireAuth(clients.Auth, logger.Logger), articleHandler.CreateArticle)
			articles.PUT("/:id", middleware.RequireAuth(clients.Auth, logger.Logger), articleHandler.UpdateArticle)
			articles.DELETE("/:id", middleware.RequireAuth(clients.Auth, logger.Logger), articleHandler.DeleteArticle)
			articles.POST("/:id/like", middleware.RequireAuth(clients.Auth, logger.Logger), articleHandler.LikeArticle)
		}
	}

	// Start server
	port := getEnv("PORT", "8080")
	logger.Logger.Info("api gateway starting", "port", port)
	if err := router.Run(":" + port); err != nil {
		log.Fatalf("failed to start server: %v", err)
	}
}

func corsMiddleware() gin.HandlerFunc {
	allowedOrigin := getEnv("ALLOW_ORIGIN", "*")

	return func(c *gin.Context) {
		c.Writer.Header().Set("Access-Control-Allow-Origin", allowedOrigin)
		c.Writer.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		c.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")
		c.Writer.Header().Set("Access-Control-Allow-Credentials", "true")

		if c.Request.Method == http.MethodOptions {
			c.AbortWithStatus(http.StatusNoContent)
			return
		}

		c.Next()
	}
}

func getEnv(key, fallback string) string {
	if value, ok := os.LookupEnv(key); ok {
		return value
	}
	return fallback
}

func parseLogLevel(value string) slog.Leveler {
	switch strings.ToLower(strings.TrimSpace(value)) {
	case "debug":
		return slog.LevelDebug
	case "warn", "warning":
		return slog.LevelWarn
	case "error":
		return slog.LevelError
	default:
		return slog.LevelInfo
	}
}

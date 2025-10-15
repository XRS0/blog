package main

import (
	"context"
	"fmt"
	"log"
	"log/slog"
	"net"
	"os"
	"strings"

	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"

	"github.com/XRS0/blog/services/article-service/internal/repository"
	"github.com/XRS0/blog/services/article-service/internal/server"
	"github.com/XRS0/blog/services/article-service/internal/service"
	pb "github.com/XRS0/blog/services/article-service/proto/article"
	sharedDB "github.com/XRS0/blog/shared/database"
	sharedLogger "github.com/XRS0/blog/shared/logger"
	"github.com/XRS0/blog/shared/rabbitmq"
)

func main() {
	// Logger setup
	logLevel := parseLogLevel(getEnv("LOG_LEVEL", "info"))
	logger := sharedLogger.New(logLevel, "article-service")

	// Database connection with Bun
	dsn := getEnv("DATABASE_URL", "postgres://blog:blog@localhost:5432/blog?sslmode=disable")
	db, err := sharedDB.NewDB(dsn, logger.Logger)
	if err != nil {
		log.Fatalf("failed to connect to database: %v", err)
	}
	defer db.Close()

	// Run auto migrations
	ctx := context.Background()
	models := []interface{}{
		(*repository.Article)(nil),
	}
	if err := sharedDB.RunMigrations(ctx, db, models, logger.Logger); err != nil {
		log.Fatalf("failed to run migrations: %v", err)
	}

	// RabbitMQ connection
	mqURL := getEnv("RABBITMQ_URL", "amqp://blog:blog@localhost:5672/")
	mq, err := rabbitmq.NewClient(mqURL, logger.Logger)
	if err != nil {
		log.Fatalf("failed to connect to RabbitMQ: %v", err)
	}
	defer mq.Close()

	// Setup exchanges and queues
	if err := mq.DeclareExchange("articles"); err != nil {
		log.Fatalf("failed to declare exchange: %v", err)
	}
	logger.Info("connected to RabbitMQ")

	// Initialize repository and service
	articleRepo := repository.NewArticleRepository(db)
	authServiceURL := getEnv("AUTH_SERVICE_URL", "localhost:50051")
	articleService, err := service.NewArticleService(articleRepo, authServiceURL, mq, logger.Logger)
	if err != nil {
		log.Fatalf("failed to create article service: %v", err)
	}

	// Create gRPC server
	grpcServer := grpc.NewServer()
	pb.RegisterArticleServiceServer(grpcServer, server.NewArticleServer(articleService, logger.Logger))

	// Enable reflection for debugging
	reflection.Register(grpcServer)

	// Start listening
	port := getEnv("GRPC_PORT", "50052")
	listener, err := net.Listen("tcp", fmt.Sprintf(":%s", port))
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	logger.Info("article service starting", "port", port)
	if err := grpcServer.Serve(listener); err != nil {
		log.Fatalf("failed to serve: %v", err)
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

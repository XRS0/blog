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

	"github.com/XRS0/blog/services/stats-service/internal/repository"
	"github.com/XRS0/blog/services/stats-service/internal/server"
	"github.com/XRS0/blog/services/stats-service/internal/service"
	pb "github.com/XRS0/blog/services/stats-service/proto"
	sharedDB "github.com/XRS0/blog/shared/database"
	sharedLogger "github.com/XRS0/blog/shared/logger"
	"github.com/XRS0/blog/shared/rabbitmq"
)

func main() {
	// Logger setup
	logLevel := parseLogLevel(getEnv("LOG_LEVEL", "info"))
	logger := sharedLogger.New(logLevel, "stats-service")

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
		(*repository.ArticleView)(nil),
		(*repository.ArticleLike)(nil),
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
	logger.Logger.Info("connected to RabbitMQ")

	// Initialize repository and service
	statsRepo := repository.NewStatsRepository(db)
	statsService := service.NewStatsService(statsRepo, mq, logger.Logger)

	// Start event consumer
	if err := statsService.StartEventConsumer(ctx); err != nil {
		log.Fatalf("failed to start event consumer: %v", err)
	}
	logger.Logger.Info("started event consumer")

	// Create gRPC server
	grpcServer := grpc.NewServer()
	pb.RegisterStatsServiceServer(grpcServer, server.NewStatsServer(statsService, logger.Logger))

	// Enable reflection for debugging
	reflection.Register(grpcServer)

	// Start listening
	port := getEnv("GRPC_PORT", "50053")
	listener, err := net.Listen("tcp", fmt.Sprintf(":%s", port))
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	logger.Logger.Info("stats service starting", "port", port)
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

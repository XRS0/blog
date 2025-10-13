.PHONY: proto clean proto-auth proto-article proto-stats

# Generate all proto files in services
proto: proto-auth proto-article proto-stats

# Generate auth service proto (in auth-service and api-gateway)
proto-auth:
	@echo "Generating auth proto for auth-service..."
	@mkdir -p services/auth-service/proto
	protoc --proto_path=proto \
		--go_out=services/auth-service/proto --go_opt=paths=source_relative \
		--go-grpc_out=services/auth-service/proto --go-grpc_opt=paths=source_relative \
		auth.proto
	@echo "Generating auth proto for api-gateway..."
	@mkdir -p services/api-gateway/proto/auth
	protoc --proto_path=proto \
		--go_out=services/api-gateway/proto/auth --go_opt=paths=source_relative \
		--go-grpc_out=services/api-gateway/proto/auth --go-grpc_opt=paths=source_relative \
		auth.proto

# Generate article service proto (in article-service and api-gateway)
proto-article:
	@echo "Generating article proto for article-service..."
	@mkdir -p services/article-service/proto
	protoc --proto_path=proto \
		--go_out=services/article-service/proto --go_opt=paths=source_relative \
		--go-grpc_out=services/article-service/proto --go-grpc_opt=paths=source_relative \
		article.proto
	@echo "Generating article proto for api-gateway..."
	@mkdir -p services/api-gateway/proto/article
	protoc --proto_path=proto \
		--go_out=services/api-gateway/proto/article --go_opt=paths=source_relative \
		--go-grpc_out=services/api-gateway/proto/article --go-grpc_opt=paths=source_relative \
		article.proto

# Generate stats service proto (in stats-service and api-gateway)
proto-stats:
	@echo "Generating stats proto for stats-service..."
	@mkdir -p services/stats-service/proto
	protoc --proto_path=proto \
		--go_out=services/stats-service/proto --go_opt=paths=source_relative \
		--go-grpc_out=services/stats-service/proto --go-grpc_opt=paths=source_relative \
		stats.proto
	@echo "Generating stats proto for api-gateway..."
	@mkdir -p services/api-gateway/proto/stats
	protoc --proto_path=proto \
		--go_out=services/api-gateway/proto/stats --go_opt=paths=source_relative \
		--go-grpc_out=services/api-gateway/proto/stats --go-grpc_opt=paths=source_relative \
		stats.proto

# Clean generated files in all services
clean:
	rm -rf services/auth-service/proto
	rm -rf services/article-service/proto
	rm -rf services/stats-service/proto
	rm -rf services/api-gateway/proto

# Build all services
build-services:
	@echo "Building API Gateway..."
	cd services/api-gateway && go build -o ../../bin/api-gateway ./cmd/main.go
	@echo "Building Auth Service..."
	cd services/auth-service && go build -o ../../bin/auth-service ./cmd/main.go
	@echo "Building Article Service..."
	cd services/article-service && go build -o ../../bin/article-service ./cmd/main.go
	@echo "Building Stats Service..."
	cd services/stats-service && go build -o ../../bin/stats-service ./cmd/main.go

# Run with docker-compose
docker-up:
	docker-compose up -d

docker-down:
	docker-compose down

docker-logs:
	docker-compose logs -f

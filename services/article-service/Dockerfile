FROM golang:1.24.5-alpine AS builder

WORKDIR /app

# Copy shared module and proto
COPY shared/ shared/
COPY proto/ proto/

# Copy service
COPY services/article-service/ services/article-service/

# Build the service
WORKDIR /app/services/article-service
RUN go mod download
RUN CGO_ENABLED=0 GOOS=linux go build -o /article-service ./cmd/main.go

# Final stage
FROM alpine:latest

RUN apk --no-cache add ca-certificates

WORKDIR /root/

COPY --from=builder /article-service .

EXPOSE 50052

CMD ["./article-service"]

FROM golang:1.21-alpine AS builder

WORKDIR /app

# Copy shared module and proto
COPY shared/ shared/
COPY proto/ proto/

# Copy service
COPY services/auth-service/ services/auth-service/

# Build the service
WORKDIR /app/services/auth-service
RUN go mod download
RUN CGO_ENABLED=0 GOOS=linux go build -o /auth-service ./cmd/main.go

# Final stage
FROM alpine:latest

RUN apk --no-cache add ca-certificates

WORKDIR /root/

COPY --from=builder /auth-service .

EXPOSE 50051

CMD ["./auth-service"]

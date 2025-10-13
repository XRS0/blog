FROM golang:1.21-alpine AS builder

WORKDIR /app

# Copy shared module and proto
COPY shared/ shared/
COPY proto/ proto/

# Copy service
COPY services/api-gateway/ services/api-gateway/

# Build the service
WORKDIR /app/services/api-gateway
RUN go mod download
RUN CGO_ENABLED=0 GOOS=linux go build -o /api-gateway ./cmd/main.go

# Final stage
FROM alpine:latest

RUN apk --no-cache add ca-certificates

WORKDIR /root/

COPY --from=builder /api-gateway .

EXPOSE 8080

CMD ["./api-gateway"]

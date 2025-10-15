FROM golang:1.24.5-alpine AS builder

WORKDIR /app

# Copy shared module and proto
COPY shared/ shared/
COPY proto/ proto/

# Copy service
COPY services/stats-service/ services/stats-service/

# Build the service
WORKDIR /app/services/stats-service
RUN go mod download
RUN CGO_ENABLED=0 GOOS=linux go build -o /stats-service ./cmd/main.go

# Final stage
FROM alpine:latest

RUN apk --no-cache add ca-certificates

WORKDIR /root/

COPY --from=builder /stats-service .

EXPOSE 50053

CMD ["./stats-service"]

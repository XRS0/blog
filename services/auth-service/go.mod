module github.com/XRS0/blog/services/auth-service

go 1.21

require (
	github.com/golang-jwt/jwt/v5 v5.2.0
	github.com/XRS0/blog/proto/gen/auth v0.0.0
	github.com/XRS0/blog/shared v0.0.0
	github.com/uptrace/bun v1.1.16
	github.com/uptrace/bun/dialect/pgdialect v1.1.16
	github.com/uptrace/bun/driver/pgdriver v1.1.16
	golang.org/x/crypto v0.18.0
	google.golang.org/grpc v1.60.1
	google.golang.org/protobuf v1.32.0
)

replace github.com/XRS0/blog/proto/gen/auth => ../../proto/gen/auth

replace github.com/XRS0/blog/shared => ../../shared

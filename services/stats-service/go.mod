module github.com/XRS0/blog/services/stats-service

go 1.21

require (
	github.com/XRS0/blog/proto/gen/stats v0.0.0
	github.com/XRS0/blog/shared v0.0.0
	github.com/uptrace/bun v1.1.16
	github.com/uptrace/bun/dialect/pgdialect v1.1.16
	github.com/uptrace/bun/driver/pgdriver v1.1.16
	google.golang.org/grpc v1.60.1
	google.golang.org/protobuf v1.32.0
)

replace github.com/XRS0/blog/proto/gen/stats => ../../proto/gen/stats

replace github.com/XRS0/blog/shared => ../../shared

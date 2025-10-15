module github.com/XRS0/blog/services/article-service

go 1.24.5

replace github.com/XRS0/blog/shared => ../../shared

require (
	github.com/XRS0/blog/shared v0.0.0-20251014090659-db8c789b0e32
	github.com/google/uuid v1.6.0
	github.com/uptrace/bun v1.2.5
	google.golang.org/grpc v1.76.0
	google.golang.org/protobuf v1.36.10
)

require (
	github.com/fatih/color v1.18.0 // indirect
	github.com/jinzhu/inflection v1.0.0 // indirect
	github.com/mattn/go-colorable v0.1.14 // indirect
	github.com/mattn/go-isatty v0.0.20 // indirect
	github.com/puzpuzpuz/xsync/v3 v3.5.1 // indirect
	github.com/rabbitmq/amqp091-go v1.9.0 // indirect
	github.com/tmthrgd/go-hex v0.0.0-20190904060850-447a3041c3bc // indirect
	github.com/uptrace/bun/dialect/pgdialect v1.2.5 // indirect
	github.com/uptrace/bun/driver/pgdriver v1.2.5 // indirect
	github.com/uptrace/bun/extra/bundebug v1.2.5 // indirect
	github.com/vmihailenco/msgpack/v5 v5.4.1 // indirect
	github.com/vmihailenco/tagparser/v2 v2.0.0 // indirect
	golang.org/x/crypto v0.40.0 // indirect
	golang.org/x/net v0.42.0 // indirect
	golang.org/x/sys v0.34.0 // indirect
	golang.org/x/text v0.27.0 // indirect
	google.golang.org/genproto/googleapis/rpc v0.0.0-20250804133106-a7a43d27e69b // indirect
	mellium.im/sasl v0.3.2 // indirect
)

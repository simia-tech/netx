
GO ?= go

test-nats:
	LISTEN_NETWORK=nats DIAL_NETWORK=nats NODES=nats://localhost:4222 $(GO) test -v ./...

bench-nats:
	LISTEN_NETWORK=nats DIAL_NETWORK=nats NODES=nats://localhost:4222 $(GO) test -run=none -bench=. -benchmem ./...

test-consul:
	LISTEN_NETWORK=consul://localhost@localhost:8500 DIAL_NETWORK=consul://localhost:8500 $(GO) test -v ./...

bench-consul:
	LISTEN_NETWORK=consul://localhost@localhost:8500 DIAL_NETWORK=consul://localhost:8500 $(GO) test -run=none -bench=. -benchmem ./...

test-consul-dnssrv:
	LISTEN_NETWORK=consul://localhost@localhost:8500 DIAL_NETWORK=dnssrv://localhost:8600 $(GO) test -v ./...

bench-consul-dnssrv:
	LISTEN_NETWORK=consul://localhost@localhost:8500 DIAL_NETWORK=dnssrv://localhost:8600 $(GO) test -run=none -bench=. -benchmem ./...

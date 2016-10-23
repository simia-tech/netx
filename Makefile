
GO ?= go

test-nats:
	LISTEN_NETWORK=nats://localhost:4222 DIAL_NETWORK=nats://localhost:4222 $(GO) test -v -bench=. -benchmem ./...

test-consul:
	LISTEN_NETWORK=consul://localhost@localhost:8500 DIAL_NETWORK=consul://localhost:8500 $(GO) test -v -bench=. -benchmem ./...

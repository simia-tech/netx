
GO ?= go

test-nats:
	LISTEN_NETWORK=nats \
		LISTEN_NETWORK_NODES=nats://localhost:4222 \
		DIAL_NETWORK=nats \
		DIAL_NETWORK_NODES=nats://localhost:4222 \
		$(GO) test -v ./...

bench-nats:
	LISTEN_NETWORK=nats \
		LISTEN_NETWORK_NODES=nats://localhost:4222 \
		DIAL_NETWORK=nats \
		DIAL_NETWORK_NODES=nats://localhost:4222 \
		$(GO) test -run=none -bench=. -benchmem ./...

test-consul:
	LISTEN_NETWORK=consul \
		LISTEN_NETWORK_NODES=http://localhost:8500 \
		LISTEN_LOCAL_ADDRESS=localhost:0 \
		DIAL_NETWORK=consul \
		DIAL_NETWORK_NODES=http://localhost:8500 \
		$(GO) test -v ./...

bench-consul:
	LISTEN_NETWORK=consul \
		LISTEN_NETWORK_NODES=http://localhost:8500 \
		LISTEN_LOCAL_ADDRESS=localhost:0 \
		DIAL_NETWORK=consul \
		DIAL_NETWORK_NODES=http://localhost:8500 \
	  $(GO) test -run=none -bench=. -benchmem ./...

test-consul-dnssrv:
	LISTEN_NETWORK=consul \
		LISTEN_NETWORK_NODES=http://localhost:8500 \
		LISTEN_LOCAL_ADDRESS=localhost:0 \
		DIAL_NETWORK=dnssrv \
		DIAL_NETWORK_NODES=localhost:8600 \
		$(GO) test -v ./...

bench-consul-dnssrv:
	LISTEN_NETWORK=consul \
		LISTEN_NETWORK_NODES=http://localhost:8500 \
		LISTEN_LOCAL_ADDRESS=localhost:0 \
		DIAL_NETWORK=dnssrv \
		DIAL_NETWORK_NODES=localhost:8600 \
	  $(GO) test -run=none -bench=. -benchmem ./...

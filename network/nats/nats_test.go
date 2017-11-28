package nats_test

import (
	"testing"

	_ "github.com/simia-tech/netx/network/nats"
	"github.com/simia-tech/netx/test"
	"github.com/simia-tech/netx/value"
)

var options = &test.Options{
	ListenNetwork:    "nats",
	ListenOptions:    []value.Option{value.Nodes("nats://localhost:4222")},
	DialNetwork:      "nats",
	DialOptions:      []value.Option{value.Nodes("nats://localhost:4222")},
	MulticastNetwork: "nats",
	MulticastOptions: []value.Option{value.Nodes("nats://localhost:4222")},
}

func TestConnection(t *testing.T) {
	test.ConnectionTest(t, options)
}

func BenchmarkConnection(b *testing.B) {
	test.ConnectionBenchmark(b, options)
}

func TestMulticast(t *testing.T) {
	test.MulticastTest(t, options)
}

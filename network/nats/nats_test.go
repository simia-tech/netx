package nats_test

import (
	"testing"

	"github.com/simia-tech/netx"
	_ "github.com/simia-tech/netx/network/nats"
	"github.com/simia-tech/netx/test"
)

var options = &test.Options{
	ListenNetwork: "nats",
	ListenOptions: []netx.Option{netx.Nodes("nats://localhost:4222")},
	DialNetwork:   "nats",
	DialOptions:   []netx.Option{netx.Nodes("nats://localhost:4222")},
}

func TestConnection(t *testing.T) {
	test.ConnectionTest(t, options)
}

func BenchmarkConnection(b *testing.B) {
	test.ConnectionBenchmark(b, options)
}

func TestRandomBalancing(t *testing.T) {
	test.RandomBalancingTest(t, options)
}

func BenchmarkBalancing(b *testing.B) {
	test.RandomBalancingBenchmark(b, options)
}

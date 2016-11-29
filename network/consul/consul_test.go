package consul_test

import (
	"testing"

	"github.com/simia-tech/netx"
	_ "github.com/simia-tech/netx/network/consul"
	"github.com/simia-tech/netx/test"
)

var options = &test.Options{
	ListenNetwork: "consul",
	ListenOptions: []netx.Option{netx.Nodes("http://localhost:8500"), netx.PublicAddress("127.0.0.1:0")},
	DialNetwork:   "consul",
	DialOptions:   []netx.Option{netx.Nodes("http://localhost:8500")},
}

func TestConnection(t *testing.T) {
	test.ConnectionTest(t, options)
}

func BenchmarkConnection(b *testing.B) {
	test.ConnectionBenchmark(b, options)
}

func TestBalancing(t *testing.T) {
	test.BalancingTest(t, options)
}

func BenchmarkBalancing(b *testing.B) {
	test.BalancingBenchmark(b, options)
}

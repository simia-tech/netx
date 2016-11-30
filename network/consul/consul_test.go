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

func TestFailover(t *testing.T) {
	test.FailoverTest(t, options)
}

func TestRandomBalancing(t *testing.T) {
	test.RandomBalancingTest(t, options)
}

func TestRoundRobinBalancing(t *testing.T) {
	test.RoundRobinBalancingTest(t, options)
}

func BenchmarkConnection(b *testing.B) {
	test.ConnectionBenchmark(b, options)
}

func BenchmarkRandomBalancing(b *testing.B) {
	test.RandomBalancingBenchmark(b, options)
}

func BenchmarkRoundRobinBalancing(b *testing.B) {
	test.RoundRobinBalancingBenchmark(b, options)
}

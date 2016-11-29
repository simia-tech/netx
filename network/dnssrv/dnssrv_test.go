package dnssrv_test

import (
	"testing"

	"github.com/simia-tech/netx"
	_ "github.com/simia-tech/netx/network/consul"
	"github.com/simia-tech/netx/test"
)

var options = &test.Options{
	ListenNetwork: "consul",
	ListenOptions: []netx.Option{netx.Nodes("http://127.0.0.1:8500"), netx.PublicAddress("127.0.0.1:0")},
	DialNetwork:   "dnssrv",
	DialOptions:   []netx.Option{netx.Nodes("127.0.0.1:8600")},
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

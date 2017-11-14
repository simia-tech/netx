package dnssrv_test

import (
	"testing"

	"github.com/simia-tech/netx"
	"github.com/simia-tech/netx/test"
	"github.com/simia-tech/netx/value"
)

var options = &test.Options{
	ListenNetwork: "consul",
	ListenOptions: []netx.Option{netx.Nodes("http://127.0.0.1:8500"), netx.PublicAddress("127.0.0.1:0")},
	DialNetwork:   "dnssrv",
	DialOptions:   []value.DialOption{value.Nodes("127.0.0.1:8600")},
}

func TestConnection(t *testing.T) {
	test.ConnectionTest(t, options)
}

func TestFailover(t *testing.T) {
	test.FailoverTest(t, options)
}

func BenchmarkConnection(b *testing.B) {
	test.ConnectionBenchmark(b, options)
}

package quic_test

import (
	"testing"

	_ "github.com/simia-tech/netx/network/quic"
	"github.com/simia-tech/netx/test"
)

var options = &test.Options{
	ListenNetwork: "quic",
	ListenAddress: "localhost:0",
	DialNetwork:   "quic",
}

func TestConnection(t *testing.T) {
	test.ConnectionTest(t, options)
}

func BenchmarkConnection(b *testing.B) {
	test.ConnectionBenchmark(b, options)
}

package netx_test

import (
	"testing"

	"github.com/simia-tech/netx/test"
)

var options = &test.Options{
	ListenNetwork: "tcp",
	ListenAddress: "localhost:0",
	DialNetwork:   "tcp",
}

func TestConnection(t *testing.T) {
	test.ConnectionTest(t, options)
}

func BenchmarkConnection(b *testing.B) {
	test.ConnectionBenchmark(b, options)
}

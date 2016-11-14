package test

import "testing"

// BalancingTest runs a balancing test.
func BalancingTest(t *testing.T, options *Options) {
	address, close := makeEchoListeners(t, 2, options)
	defer close()
	makeEchoCalls(t, 4, address, options)
}

// BalancingBenchmark runs a balancing benchmark.
func BalancingBenchmark(b *testing.B, options *Options) {
	address, close := makeEchoListeners(b, 2, options)
	defer close()
	b.ResetTimer()
	makeEchoCalls(b, b.N, address, options)
}

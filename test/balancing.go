package test

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

// BalancingTest runs a balancing test.
func BalancingTest(t *testing.T, options *Options) {
	address, counters, close := makeEchoListeners(t, 2, options)
	defer close()

	makeEchoCalls(t, 4, address, options)

	assert.Equal(t, 4, sum(counters()))
}

// BalancingBenchmark runs a balancing benchmark.
func BalancingBenchmark(b *testing.B, options *Options) {
	address, counters, close := makeEchoListeners(b, 2, options)
	defer close()

	b.ResetTimer()
	makeEchoCalls(b, b.N, address, options)
	b.StopTimer()

	assert.Equal(b, b.N, sum(counters()))
}

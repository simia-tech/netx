package test

import (
	"math/rand"
	"testing"
	"time"

	"github.com/simia-tech/netx"
	"github.com/stretchr/testify/assert"
)

// BalancingTest runs balancing tests.
func BalancingTest(t *testing.T, options *Options) {
	t.Run("Random", func(t *testing.T) {
		rand.Seed(time.Now().UnixNano())

		options.DialOptions = append(options.DialOptions, netx.Balancing(netx.RandomBalancing))
		address, counters, close := makeEchoListeners(t, 2, options)
		defer close()

		makeEchoCalls(t, 4, address, options)

		assert.Equal(t, 4, sum(counters()))
	})
}

// BalancingBenchmark runs balancing benchmarks.
func BalancingBenchmark(b *testing.B, options *Options) {
	b.Run("Random", func(b *testing.B) {
		rand.Seed(time.Now().UnixNano())

		options.DialOptions = append(options.DialOptions, netx.Balancing(netx.RandomBalancing))
		address, counters, close := makeEchoListeners(b, 2, options)
		defer close()

		b.ResetTimer()
		makeEchoCalls(b, b.N, address, options)
		b.StopTimer()

		assert.Equal(b, b.N, sum(counters()))
	})
}

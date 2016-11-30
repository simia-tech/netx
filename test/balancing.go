package test

import (
	"math/rand"
	"testing"
	"time"

	"github.com/simia-tech/netx"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// RandomBalancingTest runs random balancing tests.
func RandomBalancingTest(t *testing.T, options *Options) {
	rand.Seed(time.Now().UnixNano())
	balancer := netx.RandomBalancer()
	options.DialOptions = append(options.DialOptions, netx.Balancer(balancer), netx.DialTimeout(100*time.Millisecond))

	t.Run("ZeroNodes", func(t *testing.T) {
		err := makeEchoCalls(1, "missing", options)
		assert.Equal(t, netx.ErrServiceUnavailable, err)
	})
	t.Run("TwoNodes", func(t *testing.T) {
		address, counters, close := makeEchoListeners(t, 2, options)
		defer close()

		require.NoError(t, makeEchoCalls(4, address, options))

		time.Sleep(100 * time.Millisecond)
		assert.Equal(t, 4, sum(counters()))
	})
}

// RoundRobinBalancingTest runs random balancing tests.
func RoundRobinBalancingTest(t *testing.T, options *Options) {
	balancer := netx.RoundRobinBalancer()
	options.DialOptions = append(options.DialOptions, netx.Balancer(balancer), netx.DialTimeout(100*time.Millisecond))

	t.Run("ZeroNodes", func(t *testing.T) {
		err := makeEchoCalls(1, "missing", options)
		assert.Equal(t, netx.ErrServiceUnavailable, err)
	})
	t.Run("TwoNodes", func(t *testing.T) {
		address, counters, close := makeEchoListeners(t, 2, options)
		defer close()

		require.NoError(t, makeEchoCalls(4, address, options))

		time.Sleep(100 * time.Millisecond)
		assert.Equal(t, []int{2, 2}, counters())
	})
}

// RandomBalancingBenchmark runs balancing benchmarks.
func RandomBalancingBenchmark(b *testing.B, options *Options) {
	rand.Seed(time.Now().UnixNano())
	balancer := netx.RandomBalancer()
	options.DialOptions = append(options.DialOptions, netx.Balancer(balancer), netx.DialTimeout(100*time.Millisecond))

	b.Run("TwoNodes", func(b *testing.B) {
		address, counters, close := makeEchoListeners(b, 2, options)
		defer close()

		b.ResetTimer()
		require.NoError(b, makeEchoCalls(b.N, address, options))
		b.StopTimer()

		time.Sleep(100 * time.Millisecond)
		assert.Equal(b, b.N, sum(counters()))
	})
}

// RoundRobinBalancingBenchmark runs balancing benchmarks.
func RoundRobinBalancingBenchmark(b *testing.B, options *Options) {
	rand.Seed(time.Now().UnixNano())
	balancer := netx.RoundRobinBalancer()
	options.DialOptions = append(options.DialOptions, netx.Balancer(balancer), netx.DialTimeout(100*time.Millisecond))

	b.Run("TwoNodes", func(b *testing.B) {
		address, counters, close := makeEchoListeners(b, 2, options)
		defer close()

		b.ResetTimer()
		require.NoError(b, makeEchoCalls(b.N, address, options))
		b.StopTimer()

		time.Sleep(100 * time.Millisecond)
		assert.Equal(b, b.N, sum(counters()))
	})
}

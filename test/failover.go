package test

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func FailoverTest(t *testing.T, options *Options) {
	t.Run("OneMissingNode", func(t *testing.T) {
		address, counters, kill, close := makeFailingEchoListeners(t, 2, options.Clone())
		defer close()

		require.NoError(t, kill(0))

		require.NoError(t, makeEchoCalls(4, address, options))

		time.Sleep(100 * time.Millisecond)
		assert.Equal(t, 4, sum(counters()))
	})
}

package test

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func FailoverTest(t *testing.T, options *Options) {
	t.Run("OneMissingNode", func(t *testing.T) {
		address, counters, kill, close := makeKillableListeners(t, 2, echoServer, options.Clone())
		defer close()

		kill(0)

		require.NoError(t, makeCalls(4, address, echoClient, options))

		time.Sleep(100 * time.Millisecond)
		assert.Equal(t, 4, sum(counters()))
	})
}

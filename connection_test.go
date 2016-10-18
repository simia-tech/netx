package netx_test

import (
	"testing"

	"github.com/simia-tech/netx"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestConnection(t *testing.T) {
	listener := setUpTestListener(t)
	defer listener.Close()

	conn, err := netx.Dial(defaultNatsURL, "test")
	require.NoError(t, err)
	defer conn.Close()

	requireWrite(t, conn, []byte("test"))
	buffer := requireRead(t, conn, 4)

	assert.Equal(t, "test", string(buffer))
}

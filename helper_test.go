package netx_test

import (
	"io"
	"io/ioutil"
	"testing"

	"github.com/stretchr/testify/require"
)

func requireRead(tb testing.TB, r io.Reader, n int64) []byte {
	data, err := ioutil.ReadAll(io.LimitReader(r, n))
	require.NoError(tb, err)
	require.Len(tb, data, int(n))
	return data
}

func requireWrite(tb testing.TB, w io.Writer, data []byte) {
	n, err := w.Write(data)
	require.NoError(tb, err)
	require.Equal(tb, len(data), n)
}

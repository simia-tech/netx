package netx_test

import (
	"encoding/binary"
	"io"
	"io/ioutil"
	"testing"

	"github.com/stretchr/testify/require"
)

func requireRead(tb testing.TB, r io.Reader) []byte {
	length := uint16(0)
	require.NoError(tb, binary.Read(r, binary.BigEndian, &length))

	data, err := ioutil.ReadAll(io.LimitReader(r, int64(length)))
	require.NoError(tb, err)
	require.Len(tb, data, int(length))

	return data
}

func requireWrite(tb testing.TB, w io.Writer, data []byte) {
	length := uint16(len(data))
	require.NoError(tb, binary.Write(w, binary.BigEndian, &length))

	n, err := w.Write(data)
	require.NoError(tb, err)
	require.Equal(tb, len(data), n)
}

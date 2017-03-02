package test

import (
	"encoding/binary"
	"fmt"
	"io"
	"io/ioutil"
	"testing"

	"github.com/stretchr/testify/require"
)

func ReadBlock(r io.Reader) ([]byte, error) {
	length := uint32(0)
	if err := binary.Read(r, binary.BigEndian, &length); err != nil {
		return nil, err
	}

	data, err := ioutil.ReadAll(io.LimitReader(r, int64(length)))
	if err != nil {
		return nil, err
	}

	return data, nil
}

func WriteBlock(w io.Writer, data []byte) error {
	length := uint32(len(data))

	buffer := make([]byte, len(data)+4)
	binary.BigEndian.PutUint32(buffer, length)
	copy(buffer[4:], data)

	n, err := w.Write(buffer)
	if err != nil {
		return err
	}
	if n != len(buffer) {
		return fmt.Errorf("not all data was written (%d of %d bytes)", n, len(data))
	}

	return nil
}

func RequireReadBlock(tb testing.TB, r io.Reader) []byte {
	data, err := ReadBlock(r)
	require.NoError(tb, err)
	return data
}

func RequireWriteBlock(tb testing.TB, w io.Writer, data []byte) {
	require.NoError(tb, WriteBlock(w, data))
}

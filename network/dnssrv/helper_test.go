package dnssrv_test

import (
	"bytes"
	"encoding/binary"
	"io"
	"io/ioutil"
	"testing"

	"github.com/pkg/errors"
	"github.com/stretchr/testify/require"
)

func readBlock(r io.Reader) ([]byte, error) {
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

func writeBlock(w io.Writer, data []byte) error {
	length := uint32(len(data))
	if err := binary.Write(w, binary.BigEndian, &length); err != nil {
		return err
	}

	n, err := io.Copy(w, bytes.NewBuffer(data))
	if err != nil {
		return err
	}
	if n != int64(len(data)) {
		return errors.Errorf("not all data was written (%d of %d bytes)", n, len(data))
	}

	return nil
}

func requireRead(tb testing.TB, r io.Reader) []byte {
	data, err := readBlock(r)
	require.NoError(tb, err)
	return data
}

func requireWrite(tb testing.TB, w io.Writer, data []byte) {
	require.NoError(tb, writeBlock(w, data))
}

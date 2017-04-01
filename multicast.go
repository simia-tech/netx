package netx

import "io"

type multicast struct {
	listener io.ReadCloser
	conn     io.WriteCloser
}

// ListenAndDialMulticast listens to the provided readAddress and dials the provided writeAddress. The result is combined
// in the returned io.ReadWriteCloser interface.
func ListenAndDialMulticast(network, readAddress, writeAddress string, options ...Option) (io.ReadWriteCloser, error) {
	listener, err := ListenMulticast(network, readAddress, options...)
	if err != nil {
		return nil, err
	}

	conn, err := DialMulticast(network, writeAddress, options...)
	if err != nil {
		return nil, err
	}

	return &multicast{listener: listener, conn: conn}, nil
}

func (m *multicast) Read(buffer []byte) (int, error) {
	return m.listener.Read(buffer)
}

func (m *multicast) Write(buffer []byte) (int, error) {
	return m.conn.Write(buffer)
}

func (m *multicast) Close() error {
	if err := m.listener.Close(); err != nil {
		return err
	}
	if err := m.conn.Close(); err != nil {
		return err
	}
	return nil
}

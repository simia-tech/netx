package udp

import (
	"io"
	"net"

	"github.com/simia-tech/netx"
	"github.com/simia-tech/netx/value"
)

type multicastListener struct {
	conn       *net.UDPConn
	readBuffer []byte
}

func init() {
	netx.RegisterListenMulticast("udp", ListenMulticast)
}

func ListenMulticast(address string, options *value.Options) (io.ReadCloser, error) {
	addr, err := net.ResolveUDPAddr("udp", address)
	if err != nil {
		return nil, err
	}

	listener, err := net.ListenMulticastUDP("udp", nil, addr)
	if err != nil {
		return nil, err
	}
	listener.SetReadBuffer(ReadBufferSize)

	return &multicastListener{
		conn: listener,
	}, nil
}

func (ml *multicastListener) Read(readBuffer []byte) (int, error) {
	if len(ml.readBuffer) > 0 {
		n := copy(readBuffer, ml.readBuffer)
		if n < len(ml.readBuffer) {
			ml.readBuffer = ml.readBuffer[n:]
		} else {
			ml.readBuffer = nil
		}
		return n, nil
	}

	ml.readBuffer = make([]byte, ReadBufferSize)
	n, _, err := ml.conn.ReadFromUDP(ml.readBuffer)
	if err != nil {
		return 0, err
	}
	ml.readBuffer = ml.readBuffer[:n]
	n = copy(readBuffer, ml.readBuffer)
	if n < len(ml.readBuffer) {
		ml.readBuffer = ml.readBuffer[n:]
	}
	return n, nil
}

func (ml *multicastListener) Close() error {
	if err := ml.conn.Close(); err != nil {
		return err
	}
	return nil
}

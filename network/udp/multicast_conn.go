package udp

import (
	"io"
	"net"

	"github.com/simia-tech/netx"
)

type multicastConn struct {
	conn *net.UDPConn
}

func init() {
	netx.RegisterDialMulticast("udp", DialMulticast)
}

func DialMulticast(address string, options *netx.Options) (io.WriteCloser, error) {
	addr, err := net.ResolveUDPAddr("udp", address)
	if err != nil {
		return nil, err
	}

	conn, err := net.DialUDP("udp", nil, addr)
	if err != nil {
		return nil, err
	}

	return &multicastConn{
		conn: conn,
	}, nil
}

func (mc *multicastConn) Write(data []byte) (int, error) {
	return mc.conn.Write(data)
}

func (mc *multicastConn) Close() error {
	if err := mc.conn.Close(); err != nil {
		return err
	}
	return nil
}

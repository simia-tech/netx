package nats

import (
	"fmt"
	"net"
	"strings"
	"time"

	n "github.com/nats-io/nats"

	"github.com/simia-tech/netx"
	"github.com/simia-tech/netx/model"
)

type multicastConn struct {
	conn    *n.Conn
	address string

	writeDeadline time.Time

	maxPacketSize int
}

func init() {
	netx.RegisterDialMulticast("nats", DialMulticast)
}

func DialMulticast(address string, options *netx.Options) (net.Conn, error) {
	o := []n.Option{}
	if options.TLSConfig != nil {
		o = append(o, n.Secure(options.TLSConfig))
	}
	conn, err := n.Connect(strings.Join(options.Nodes, ","), o...)
	if err != nil {
		return nil, err
	}

	host, _, err := net.SplitHostPort(address)
	if err != nil {
		host = address
	}

	return newMulticastConn(conn, host)
}

func newMulticastConn(nc *n.Conn, address string) (*multicastConn, error) {
	return &multicastConn{
		conn:          nc,
		address:       address,
		maxPacketSize: int(nc.MaxPayload() - overheadSize),
	}, nil
}

func (mc *multicastConn) Read(readBuffer []byte) (int, error) {
	return 0, netx.ErrNotSupported
}

func (mc *multicastConn) Write(data []byte) (int, error) {
	total := 0
	for len(data) > mc.maxPacketSize {
		if err := mc.sendPacket(model.Packet_DATA, data[:mc.maxPacketSize]); err != nil {
			return total, err
		}
		data = data[mc.maxPacketSize:]
		total += mc.maxPacketSize
	}

	if err := mc.sendPacket(model.Packet_DATA, data); err != nil {
		return 0, err
	}
	total += len(data)

	return total, nil
}

func (mc *multicastConn) Close() error {
	return nil
}

func (mc *multicastConn) LocalAddr() net.Addr {
	return &addr{net: "nats", address: "multi"}
}

func (mc *multicastConn) RemoteAddr() net.Addr {
	return &addr{net: "nats", address: mc.address}
}

func (mc *multicastConn) SetDeadline(t time.Time) error {
	mc.writeDeadline = t
	return nil
}

func (mc *multicastConn) SetReadDeadline(t time.Time) error {
	return netx.ErrNotSupported
}

func (mc *multicastConn) SetWriteDeadline(t time.Time) error {
	mc.writeDeadline = t
	return nil
}

func (mc *multicastConn) String() string {
	return fmt.Sprintf("(%s -> %s)", mc.LocalAddr(), mc.RemoteAddr())
}

func (mc *multicastConn) sendPacket(t model.Packet_Type, payload []byte) error {
	return sendPacket(mc.conn, mc.address, t, payload)
}

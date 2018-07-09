package nats

import (
	"context"
	"fmt"
	"io"
	"net"
	"strings"
	"time"

	n "github.com/nats-io/nats"

	"github.com/simia-tech/netx"
	"github.com/simia-tech/netx/model"
	"github.com/simia-tech/netx/value"
)

const overheadSize = 100

type conn struct {
	conn          *n.Conn
	connDedicated bool
	localInbox    string
	remoteInbox   string

	subscription *n.Subscription

	readDeadline  time.Time
	writeDeadline time.Time

	readBuffer    []byte
	maxPacketSize int
}

func init() {
	netx.RegisterDial("nats", Dial)
}

// Dial establishes a connection to the provided address on the provided network.
func Dial(ctx context.Context, address string, options *value.Options) (net.Conn, error) {
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

	localInbox := n.NewInbox()
	c, err := newConn(conn, true, localInbox, host)
	if err != nil {
		return nil, err
	}

	if err = c.sendPacket(model.Packet_NEW, []byte(localInbox)); err != nil {
		return nil, err
	}

	packet := (*model.Packet)(nil)
	if timeout := options.Timeout; timeout == 0 {
		packet, err = c.receivePacket()
	} else {
		packet, err = receivePacket(c.subscription, timeout)
	}
	if err == n.ErrTimeout {
		return nil, netx.ErrServiceUnavailable
	}
	if err != nil {
		return nil, err
	}
	if packet.Type != model.Packet_ACCEPT {
		return nil, fmt.Errorf("expected ACCEPT packet, got %s", packet.Type)
	}
	c.remoteInbox = string(packet.Payload)

	return c, nil
}

func newConn(nc *n.Conn, connDedicated bool, localInbox, remoteInbox string) (*conn, error) {
	subscription, err := nc.SubscribeSync(localInbox)
	if err != nil {
		return nil, err
	}

	return &conn{
		conn:          nc,
		connDedicated: connDedicated,
		localInbox:    localInbox,
		remoteInbox:   remoteInbox,
		subscription:  subscription,
		maxPacketSize: int(nc.MaxPayload() - overheadSize),
	}, nil
}

func (c *conn) Read(readBuffer []byte) (int, error) {
	if len(c.readBuffer) > 0 {
		n := copy(readBuffer, c.readBuffer)
		if n < len(c.readBuffer) {
			c.readBuffer = c.readBuffer[n:]
		} else {
			c.readBuffer = nil
		}
		return n, nil
	}

	if c.subscription == nil {
		return 0, io.EOF
	}

	packet, err := c.receivePacket()
	if err != nil {
		return 0, err
	}
	switch packet.Type {
	case model.Packet_DATA:
		n := copy(readBuffer, packet.Payload)
		if n < len(packet.Payload) {
			c.readBuffer = packet.Payload[n:]
		}
		return n, nil
	case model.Packet_CLOSE:
		return 0, io.EOF
	default:
		return 0, fmt.Errorf("expected DATA packet, got %s", packet.Type)
	}
}

func (c *conn) Write(data []byte) (int, error) {
	if c.subscription == nil {
		return 0, io.ErrClosedPipe
	}

	total := 0
	for len(data) > c.maxPacketSize {
		if err := c.sendPacket(model.Packet_DATA, data[:c.maxPacketSize]); err != nil {
			return total, err
		}
		data = data[c.maxPacketSize:]
		total += c.maxPacketSize
	}

	if err := c.sendPacket(model.Packet_DATA, data); err != nil {
		return 0, err
	}
	total += len(data)

	return total, nil
}

func (c *conn) Close() error {
	if c.subscription == nil {
		return nil
	}
	if err := c.sendPacket(model.Packet_CLOSE, nil); err != nil {
		return err
	}
	if err := c.subscription.Unsubscribe(); err != nil {
		return err
	}
	c.subscription = nil

	if c.connDedicated {
		c.conn.Close()
	}

	return nil
}

func (c *conn) LocalAddr() net.Addr {
	return &addr{net: c.conn.Opts.Name, address: c.localInbox}
}

func (c *conn) RemoteAddr() net.Addr {
	return &addr{net: c.conn.Opts.Name, address: c.remoteInbox}
}

func (c *conn) SetDeadline(t time.Time) error {
	c.readDeadline = t
	c.writeDeadline = t
	return nil
}

func (c *conn) SetReadDeadline(t time.Time) error {
	c.readDeadline = t
	return nil
}

func (c *conn) SetWriteDeadline(t time.Time) error {
	c.writeDeadline = t
	return nil
}

func (c *conn) String() string {
	return fmt.Sprintf("(%s -> %s)", c.LocalAddr(), c.RemoteAddr())
}

func (c *conn) sendPacket(t model.Packet_Type, payload []byte) error {
	return sendPacket(c.conn, c.remoteInbox, t, payload)
}

func (c *conn) receivePacket() (*model.Packet, error) {
	if c.readDeadline.IsZero() {
		return receivePacket(c.subscription, endlessTimeout)
	}
	return receivePacket(c.subscription, c.readDeadline.Sub(time.Now()))
}

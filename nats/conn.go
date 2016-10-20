package nats

import (
	"fmt"
	"io"
	"log"
	"net"
	"time"

	n "github.com/nats-io/nats"
	"github.com/pkg/errors"
	"github.com/simia-tech/netx/model"
)

type conn struct {
	conn        *n.Conn
	localInbox  string
	remoteInbox string

	subscription *n.Subscription
	dataChan     chan []byte

	readDeadline  time.Time
	writeDeadline time.Time
}

func Dial(network, address string) (net.Conn, error) {
	conn, err := n.Connect(network)
	if err != nil {
		return nil, err
	}

	host, _, err := net.SplitHostPort(address)
	if err != nil {
		host = address
	}

	message, err := conn.Request(host, []byte{}, 2*time.Second)
	if err != nil {
		return nil, errors.Wrapf(err, "requesting address from [%s] failed", host)
	}

	packet, err := receivePacket(message.Data)
	if err != nil {
		return nil, err
	}
	if packet.Type != model.Packet_ACCEPT {
		return nil, errors.Errorf("unexpected packet type %s", packet.Type)
	}

	remoteInbox := string(packet.Payload)

	return newConn(conn, message.Subject, remoteInbox)
}

func newConn(nc *n.Conn, localInbox, remoteInbox string) (*conn, error) {
	c := &conn{
		conn:        nc,
		localInbox:  localInbox,
		remoteInbox: remoteInbox,
	}

	dataChan := make(chan []byte)
	subscription, err := nc.Subscribe(localInbox, func(message *n.Msg) {
		packet, err := receivePacket(message.Data)
		if err != nil {
			log.Println(err)
			return
		}

		switch packet.Type {
		case model.Packet_DATA:
			dataChan <- packet.Payload
		case model.Packet_CLOSE:
			c.subscription.Unsubscribe()
			c.subscription = nil
			close(dataChan)
		default:
			log.Printf("unknown packet type %s", packet.Type)
		}
	})
	if err != nil {
		return nil, err
	}

	c.subscription = subscription
	c.dataChan = dataChan

	return c, nil
}

func (c *conn) Read(buffer []byte) (int, error) {
	if c.subscription == nil {
		return 0, io.ErrClosedPipe
	}
	if c.readDeadline.IsZero() {
		data, ok := <-c.dataChan
		if !ok {
			return 0, io.ErrClosedPipe
		}
		return copy(buffer, data), nil
	} else {
		select {
		case data, ok := <-c.dataChan:
			if !ok {
				return 0, io.ErrClosedPipe
			}
			return copy(buffer, data), nil
		case <-time.After(c.readDeadline.Sub(time.Now())):
			return 0, errors.New("timeout")
		}
	}
}

func (c *conn) Write(buffer []byte) (int, error) {
	if c.subscription == nil {
		return 0, io.ErrClosedPipe
	}
	if err := c.sendPacket(model.Packet_DATA, buffer); err != nil {
		return 0, err
	}
	return len(buffer), nil
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

	c.conn.Close()

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

func sendPacket(conn *n.Conn, address string, t model.Packet_Type, payload []byte) error {
	packet := &model.Packet{
		Type:    t,
		Payload: payload,
	}
	data, err := packet.MarshalBinary()
	if err != nil {
		return err
	}
	if err := conn.Publish(address, data); err != nil {
		return err
	}
	return nil
}

func receivePacket(data []byte) (*model.Packet, error) {
	packet := &model.Packet{}
	if err := packet.UnmarshalBinary(data); err != nil {
		return nil, err
	}
	return packet, nil
}

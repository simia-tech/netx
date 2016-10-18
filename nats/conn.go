package nats

import (
	"fmt"
	"io"
	"log"
	"net"
	"time"

	n "github.com/nats-io/nats"
	"github.com/simia-tech/netx/model"
)

type conn struct {
	network     *network
	localInbox  string
	remoteInbox string

	subscription *n.Subscription
	dataChan     chan []byte
}

func newConn(network *network, localInbox, remoteInbox string) (*conn, error) {
	dataChan := make(chan []byte)
	subscription, err := network.conn.Subscribe(localInbox, func(message *n.Msg) {
		packet, err := receivePacket(message.Data)
		if err != nil {
			log.Println(err)
			return
		}

		switch packet.Type {
		case model.Packet_DATA:
			dataChan <- packet.Payload
		case model.Packet_CLOSE:
			close(dataChan)
		default:
			log.Printf("unknown packet type %s", packet.Type)
		}
	})
	if err != nil {
		return nil, err
	}

	return &conn{
		network:      network,
		localInbox:   localInbox,
		remoteInbox:  remoteInbox,
		subscription: subscription,
		dataChan:     dataChan,
	}, nil
}

func (c *conn) Read(buffer []byte) (int, error) {
	data, ok := <-c.dataChan
	if !ok {
		return 0, io.ErrClosedPipe
	}
	return copy(buffer, data), nil
}

func (c *conn) Write(buffer []byte) (int, error) {
	if err := c.sendPacket(model.Packet_DATA, buffer); err != nil {
		return 0, err
	}
	return len(buffer), nil
}

func (c *conn) Close() error {
	if err := c.sendPacket(model.Packet_CLOSE, nil); err != nil {
		return err
	}
	if err := c.subscription.Unsubscribe(); err != nil {
		return err
	}
	return nil
}

func (c *conn) LocalAddr() net.Addr {
	return &addr{network: c.network, address: c.localInbox}
}

func (c *conn) RemoteAddr() net.Addr {
	return &addr{network: c.network, address: c.remoteInbox}
}

func (c *conn) SetDeadline(t time.Time) error {
	return nil
}

func (c *conn) SetReadDeadline(t time.Time) error {
	return nil
}

func (c *conn) SetWriteDeadline(t time.Time) error {
	return nil
}

func (c *conn) String() string {
	return fmt.Sprintf("(%s -> %s)", c.LocalAddr(), c.RemoteAddr())
}

func (c *conn) sendPacket(t model.Packet_Type, payload []byte) error {
	return sendPacket(c.network.conn, c.remoteInbox, t, payload)
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

package nats

import (
	"fmt"
	"net"
	"time"

	n "github.com/nats-io/nats"
)

type conn struct {
	network     *network
	localInbox  string
	remoteInbox string

	messageSubscription *n.Subscription
	messageChan         chan *n.Msg
}

func newConn(network *network, localInbox, remoteInbox string) (*conn, error) {
	messageChan := make(chan *n.Msg)
	subscription, err := network.conn.ChanSubscribe(localInbox, messageChan)
	if err != nil {
		return nil, err
	}

	return &conn{
		network:             network,
		localInbox:          localInbox,
		remoteInbox:         remoteInbox,
		messageSubscription: subscription,
		messageChan:         messageChan,
	}, nil
}

func (c *conn) Read(buffer []byte) (int, error) {
	message := <-c.messageChan
	return copy(buffer, message.Data), nil
}

func (c *conn) Write(buffer []byte) (int, error) {
	if err := c.network.conn.Publish(c.remoteInbox, buffer); err != nil {
		return 0, err
	}
	return len(buffer), nil
}

func (c *conn) Close() error {
	return nil
}

func (c *conn) LocalAddr() net.Addr {
	return nil
}

func (c *conn) RemoteAddr() net.Addr {
	return nil
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
	return fmt.Sprintf("(%s -> %s)", c.localInbox, c.remoteInbox)
}

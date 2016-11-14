package nats

import (
	"io"
	"math"
	"net"
	"strings"
	"time"

	n "github.com/nats-io/nats"
	"github.com/pkg/errors"

	"github.com/simia-tech/netx"
	"github.com/simia-tech/netx/model"
)

var endlessTimeout = time.Duration(math.MaxInt64)

type listener struct {
	conn         *n.Conn
	subscription *n.Subscription
}

func init() {
	netx.RegisterListen("nats", Listen)
}

// Listen starts a listener at the provided address on the nats network.
func Listen(address string, options *netx.Options) (net.Listener, error) {
	o := []n.Option{}
	if options.TLSConfig != nil {
		o = append(o, n.Secure(options.TLSConfig))
	}
	conn, err := n.Connect(strings.Join(options.Nodes, ","), o...)
	if err != nil {
		return nil, err
	}
	subscription, err := conn.QueueSubscribeSync(address, address)
	if err != nil {
		return nil, err
	}
	return &listener{
		conn:         conn,
		subscription: subscription,
	}, nil
}

func (l *listener) Accept() (net.Conn, error) {
	if l.subscription == nil {
		return nil, io.ErrClosedPipe
	}

	packet, err := receivePacket(l.subscription, endlessTimeout)
	if err != nil {
		return nil, err
	}
	if packet.Type != model.Packet_NEW {
		return nil, errors.Errorf("expected NEW packet, got %s", packet.Type)
	}

	localInbox := n.NewInbox()
	c, err := newConn(l.conn, false, localInbox, string(packet.Payload))
	if err != nil {
		return nil, err
	}

	if err := c.sendPacket(model.Packet_ACCEPT, []byte(localInbox)); err != nil {
		return nil, err
	}

	return c, nil
}

func (l *listener) Close() error {
	if l.subscription == nil {
		return io.ErrClosedPipe
	}

	if err := l.subscription.Unsubscribe(); err != nil {
		return err
	}

	return nil
}

func (l *listener) Addr() net.Addr {
	return &addr{net: l.conn.Opts.Name, address: l.subscription.Subject}
}

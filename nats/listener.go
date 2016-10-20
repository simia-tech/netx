package nats

import (
	"net"
	"time"

	n "github.com/nats-io/nats"

	"github.com/simia-tech/netx/model"
)

type listener struct {
	conn         *n.Conn
	address      string
	subscription *n.Subscription
}

func Listen(net, address string) (net.Listener, error) {
	conn, err := n.Connect(net)
	if err != nil {
		return nil, err
	}
	return &listener{
		conn:    conn,
		address: address,
	}, nil
}

func (l *listener) Accept() (net.Conn, error) {
	if l.subscription == nil {
		subscription, err := l.conn.QueueSubscribeSync(l.address, l.address)
		if err != nil {
			return nil, err
		}
		l.subscription = subscription
	}

	message, err := l.subscription.NextMsg(100 * time.Second)
	if err != nil {
		return nil, err
	}

	localInbox := n.NewInbox()
	if err := sendPacket(l.conn, message.Reply, model.Packet_ACCEPT, []byte(localInbox)); err != nil {
		return nil, err
	}

	return newConn(l.conn, localInbox, message.Reply)
}

func (l *listener) Close() error {
	if l.subscription != nil {
		if err := l.subscription.Unsubscribe(); err != nil {
			return err
		}
		l.subscription = nil
	}
	return nil
}

func (l *listener) Addr() net.Addr {
	return &addr{net: l.conn.Opts.Name, address: l.address}
}

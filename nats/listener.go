package nats

import (
	"net"
	"time"

	n "github.com/nats-io/nats"

	"github.com/simia-tech/netx/model"
)

type listener struct {
	network *network
	address string
}

func (l *listener) Accept() (net.Conn, error) {
	subscription, err := l.network.conn.QueueSubscribeSync(l.address, l.address)
	if err != nil {
		return nil, err
	}

	message, err := subscription.NextMsg(100 * time.Millisecond)
	if err != nil {
		return nil, err
	}

	localInbox := n.NewInbox()
	if err := sendPacket(l.network.conn, message.Reply, model.Packet_ACCEPT, []byte(localInbox)); err != nil {
		return nil, err
	}

	return newConn(l.network, localInbox, message.Reply)
}

func (l *listener) Close() error {
	return nil
}

func (l *listener) Addr() net.Addr {
	return &addr{network: l.network, address: l.address}
}

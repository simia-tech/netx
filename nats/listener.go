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
	subscription, err := l.network.conn.QueueSubscribeSync(l.address, "queue")
	if err != nil {
		return nil, err
	}

	message, err := subscription.NextMsg(100 * time.Millisecond)
	if err != nil {
		return nil, err
	}

	localInbox := n.NewInbox()
	packet := &model.Packet{
		Type:    model.Packet_ACCEPT,
		Payload: []byte(localInbox),
	}
	data, err := packet.MarshalBinary()
	if err != nil {
		return nil, err
	}

	if err := l.network.conn.Publish(message.Reply, data); err != nil {
		return nil, err
	}

	return newConn(l.network, localInbox, message.Reply)
}

func (l *listener) Close() error {
	return nil
}

func (l *listener) Addr() net.Addr {
	return nil
}

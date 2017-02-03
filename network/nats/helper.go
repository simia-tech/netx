package nats

import (
	"time"

	n "github.com/nats-io/nats"

	"github.com/simia-tech/netx/model"
)

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
	if err := conn.Flush(); err != nil {
		return err
	}
	return nil
}

func receivePacket(subscription *n.Subscription, timeout time.Duration) (*model.Packet, error) {
	message, err := subscription.NextMsg(timeout)
	if err != nil {
		return nil, err
	}
	packet := &model.Packet{}
	if err := packet.UnmarshalBinary(message.Data); err != nil {
		return nil, err
	}
	return packet, nil
}

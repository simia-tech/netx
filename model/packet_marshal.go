package model

import "github.com/gogo/protobuf/proto"

func (p *Packet) MarshalBinary() ([]byte, error) {
	return proto.Marshal(p)
}

func (p *Packet) UnmarshalBinary(data []byte) error {
	return proto.Unmarshal(data, p)
}

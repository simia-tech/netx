// Code generated by protoc-gen-go.
// source: packet.proto
// DO NOT EDIT!

/*
Package model is a generated protocol buffer package.

It is generated from these files:
	packet.proto

It has these top-level messages:
	Packet
*/
package model

import proto "github.com/golang/protobuf/proto"
import fmt "fmt"
import math "math"

// Reference imports to suppress errors if they are not otherwise used.
var _ = proto.Marshal
var _ = fmt.Errorf
var _ = math.Inf

// This is a compile-time assertion to ensure that this generated file
// is compatible with the proto package it is being compiled against.
// A compilation error at this line likely means your copy of the
// proto package needs to be updated.
const _ = proto.ProtoPackageIsVersion2 // please upgrade the proto package

type Packet_Type int32

const (
	Packet_ACCEPT Packet_Type = 0
	Packet_DATA   Packet_Type = 1
	Packet_CLOSE  Packet_Type = 2
)

var Packet_Type_name = map[int32]string{
	0: "ACCEPT",
	1: "DATA",
	2: "CLOSE",
}
var Packet_Type_value = map[string]int32{
	"ACCEPT": 0,
	"DATA":   1,
	"CLOSE":  2,
}

func (x Packet_Type) String() string {
	return proto.EnumName(Packet_Type_name, int32(x))
}
func (Packet_Type) EnumDescriptor() ([]byte, []int) { return fileDescriptor0, []int{0, 0} }

type Packet struct {
	Type    Packet_Type `protobuf:"varint,1,opt,name=type,enum=model.Packet_Type" json:"type,omitempty"`
	Payload []byte      `protobuf:"bytes,2,opt,name=payload,proto3" json:"payload,omitempty"`
}

func (m *Packet) Reset()                    { *m = Packet{} }
func (m *Packet) String() string            { return proto.CompactTextString(m) }
func (*Packet) ProtoMessage()               {}
func (*Packet) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{0} }

func init() {
	proto.RegisterType((*Packet)(nil), "model.Packet")
	proto.RegisterEnum("model.Packet_Type", Packet_Type_name, Packet_Type_value)
}

var fileDescriptor0 = []byte{
	// 143 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x09, 0x6e, 0x88, 0x02, 0xff, 0xe2, 0xe2, 0x29, 0x48, 0x4c, 0xce,
	0x4e, 0x2d, 0xd1, 0x2b, 0x28, 0xca, 0x2f, 0xc9, 0x17, 0x62, 0xcd, 0xcd, 0x4f, 0x49, 0xcd, 0x51,
	0x2a, 0xe6, 0x62, 0x0b, 0x00, 0x0b, 0x0b, 0xa9, 0x71, 0xb1, 0x94, 0x54, 0x16, 0xa4, 0x4a, 0x30,
	0x2a, 0x30, 0x6a, 0xf0, 0x19, 0x09, 0xe9, 0x81, 0xe5, 0xf5, 0x20, 0x92, 0x7a, 0x21, 0x40, 0x99,
	0x20, 0xb0, 0xbc, 0x90, 0x04, 0x17, 0x7b, 0x41, 0x62, 0x65, 0x4e, 0x7e, 0x62, 0x8a, 0x04, 0x13,
	0x50, 0x29, 0x4f, 0x10, 0x8c, 0xab, 0xa4, 0xce, 0xc5, 0x02, 0x52, 0x27, 0xc4, 0xc5, 0xc5, 0xe6,
	0xe8, 0xec, 0xec, 0x1a, 0x10, 0x22, 0xc0, 0x20, 0xc4, 0xc1, 0xc5, 0xe2, 0xe2, 0x18, 0xe2, 0x28,
	0xc0, 0x28, 0xc4, 0xc9, 0xc5, 0xea, 0xec, 0xe3, 0x1f, 0xec, 0x2a, 0xc0, 0x94, 0xc4, 0x06, 0x76,
	0x82, 0x31, 0x20, 0x00, 0x00, 0xff, 0xff, 0x98, 0x28, 0x31, 0xc6, 0x92, 0x00, 0x00, 0x00,
}

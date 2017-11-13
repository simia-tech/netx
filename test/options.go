package test

import (
	"github.com/simia-tech/netx"
	"github.com/simia-tech/netx/value"
)

// Options holds the test options.
type Options struct {
	ListenNetwork            string
	ListenAddress            string
	ListenOptions            []netx.Option
	DialNetwork              string
	DialOptions              []value.DialOption
	MulticastNetwork         string
	MulticastRequestAddress  string
	MulticastResponseAddress string
	MulticastOptions         []netx.Option
}

// Clone returns a clone of the Options struct.
func (o *Options) Clone() *Options {
	return &Options{
		ListenNetwork:            o.ListenNetwork,
		ListenAddress:            o.ListenAddress,
		ListenOptions:            cloneOptionSlice(o.ListenOptions),
		DialNetwork:              o.DialNetwork,
		DialOptions:              cloneDialOptionSlice(o.DialOptions),
		MulticastNetwork:         o.MulticastNetwork,
		MulticastRequestAddress:  o.MulticastRequestAddress,
		MulticastResponseAddress: o.MulticastResponseAddress,
		MulticastOptions:         cloneOptionSlice(o.MulticastOptions),
	}
}

func cloneOptionSlice(options []netx.Option) []netx.Option {
	result := []netx.Option{}
	for _, option := range options {
		result = append(result, option)
	}
	return result
}

func cloneDialOptionSlice(options []value.DialOption) []value.DialOption {
	result := []value.DialOption{}
	for _, option := range options {
		result = append(result, option)
	}
	return result
}

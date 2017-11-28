package test

import (
	"github.com/simia-tech/netx/value"
)

// Options holds the test options.
type Options struct {
	ListenNetwork            string
	ListenAddress            string
	ListenOptions            []value.Option
	DialNetwork              string
	DialOptions              []value.Option
	MulticastNetwork         string
	MulticastRequestAddress  string
	MulticastResponseAddress string
	MulticastOptions         []value.Option
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

func cloneOptionSlice(options []value.Option) []value.Option {
	result := []value.Option{}
	for _, option := range options {
		result = append(result, option)
	}
	return result
}

func cloneDialOptionSlice(options []value.Option) []value.Option {
	result := []value.Option{}
	for _, option := range options {
		result = append(result, option)
	}
	return result
}

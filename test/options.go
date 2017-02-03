package test

import "github.com/simia-tech/netx"

// Options holds the test options.
type Options struct {
	ListenNetwork          string
	ListenAddress          string
	ListenOptions          []netx.Option
	DialNetwork            string
	DialOptions            []netx.Option
	ListenMulticastNetwork string
	ListenMulticastOptions []netx.Option
	DialMulticastNetwork   string
	DialMulticastOptions   []netx.Option
}

func (o *Options) Clone() *Options {
	return &Options{
		ListenNetwork:          o.ListenNetwork,
		ListenAddress:          o.ListenAddress,
		ListenOptions:          cloneOptionSlice(o.ListenOptions),
		DialNetwork:            o.DialNetwork,
		DialOptions:            cloneOptionSlice(o.DialOptions),
		ListenMulticastNetwork: o.ListenMulticastNetwork,
		ListenMulticastOptions: cloneOptionSlice(o.ListenMulticastOptions),
		DialMulticastNetwork:   o.DialMulticastNetwork,
		DialMulticastOptions:   cloneOptionSlice(o.DialMulticastOptions),
	}
}

func cloneOptionSlice(options []netx.Option) []netx.Option {
	result := []netx.Option{}
	for _, option := range options {
		result = append(result, option)
	}
	return result
}

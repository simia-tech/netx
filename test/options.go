package test

import "github.com/simia-tech/netx"

// Options holds the test options.
type Options struct {
	ListenNetwork string
	ListenAddress string
	ListenOptions []netx.Option
	DialNetwork   string
	DialOptions   []netx.Option
}

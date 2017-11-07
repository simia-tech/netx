package provider

import "github.com/simia-tech/netx/value"

// Interface defines the provider interface.
type Interface interface {
	Addrs() (value.Addrs, error)
}

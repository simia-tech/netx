package netx

import (
	"crypto/tls"
	"net"
	"time"
)

var DefaultOptions = &Options{
	Nodes:          []string{},
	PublicListener: nil,
	PublicAddress:  "127.0.0.1:0",
	TLSConfig:      nil,
	Balancer:       RandomBalancer(),
	DialTimeout:    0,
}

// Options holds generic options for Listen and Dial functions.
type Options struct {
	Nodes          []string
	PublicListener net.Listener
	PublicAddress  string
	TLSConfig      *tls.Config
	Balancer       BalancerFn
	DialTimeout    time.Duration
}

// Option defines a generic option.
type Option func(*Options) error

// Nodes returns an option to set the network nodes.
func Nodes(values ...string) Option {
	return func(o *Options) error {
		o.Nodes = values
		return nil
	}
}

// PublicListener returns an option to set the public listener. The listener should be bound to a public
// network interface that can be reached by other nodes. The listener's address might be shared with
// other nodes.
func PublicListener(value net.Listener) Option {
	return func(o *Options) error {
		o.PublicListener = value
		return nil
	}
}

// PublicAddress returns an option to set the public address. The address is used to create a public listener
// unless the PublicListener option is used to set one.
func PublicAddress(value string) Option {
	return func(o *Options) error {
		o.PublicAddress = value
		return nil
	}
}

// TLS returns an option to set the tls configuration.
func TLS(value *tls.Config) Option {
	return func(o *Options) error {
		o.TLSConfig = value
		return nil
	}
}

// Balancer returns an option to set the balancer.
func Balancer(value BalancerFn) Option {
	return func(o *Options) error {
		o.Balancer = value
		return nil
	}
}

// DialTimeout returns on options to set the dial timeout.
func DialTimeout(value time.Duration) Option {
	return func(o *Options) error {
		o.DialTimeout = value
		return nil
	}
}

package netx

import "crypto/tls"

// Options holds generic options for Listen and Dial functions.
type Options struct {
	nodes        []string
	localAddress string
	tlsConfig    *tls.Config
}

// Option defines a generic option.
type Option func(*Options) error

// Nodes returns an option to set the network nodes.
func Nodes(nodes ...string) Option {
	return func(o *Options) error {
		o.nodes = nodes
		return nil
	}
}

// LocalAddress returns an option to set the local address that is used
// to bind a local listener.
func LocalAddress(value string) Option {
	return func(o *Options) error {
		o.localAddress = value
		return nil
	}
}

// TLS returns an option to set the tls configuration.
func TLS(value *tls.Config) Option {
	return func(o *Options) error {
		o.tlsConfig = value
		return nil
	}
}

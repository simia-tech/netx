package netx

import "crypto/tls"

// Options holds generic options for Listen and Dial functions.
type Options struct {
	Nodes        []string
	LocalAddress string
	TLSConfig    *tls.Config
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

// LocalAddress returns an option to set the local address that is used
// to bind a local listener.
func LocalAddress(value string) Option {
	return func(o *Options) error {
		o.LocalAddress = value
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

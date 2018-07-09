package value

import (
	"crypto/tls"
	"time"
)

// Options holds all options for a dial call.
type Options struct {
	TLSConfig *tls.Config
	Timeout   time.Duration
	Nodes     []string
}

// Option defines a function that can modify the provided DialOptions structure.
type Option func(*Options) error

// TLS returns an option to set the tls configuration.
func TLS(value *tls.Config) Option {
	return func(o *Options) error {
		o.TLSConfig = value
		return nil
	}
}

// Nodes returns on options to set the nodes.
func Nodes(value ...string) Option {
	return func(o *Options) error {
		o.Nodes = value
		return nil
	}
}

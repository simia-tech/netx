package value

import (
	"crypto/tls"
	"time"
)

// DialOptions holds all options for a dial call.
type DialOptions struct {
	TLSConfig *tls.Config
	Timeout   time.Duration
	Nodes     []string
}

// DialOption defines a function that can modify the provided DialOptions structure.
type DialOption func(*DialOptions) error

// TLS returns an option to set the tls configuration.
func TLS(value *tls.Config) DialOption {
	return func(do *DialOptions) error {
		do.TLSConfig = value
		return nil
	}
}

// Timeout returns on options to set the dial timeout.
func Timeout(value time.Duration) DialOption {
	return func(do *DialOptions) error {
		do.Timeout = value
		return nil
	}
}

// Nodes returns on options to set the nodes.
func Nodes(value ...string) DialOption {
	return func(do *DialOptions) error {
		do.Nodes = value
		return nil
	}
}

package static

import (
	"net"

	"github.com/simia-tech/netx/provider"
	"github.com/simia-tech/netx/value"
)

type static struct {
	addrs []net.Addr
}

// NewProvider returns a new static provider.
func NewProvider(addrs ...net.Addr) provider.Interface {
	return &static{addrs: addrs}
}

// NewProviderFromURLs returns a new static provider made from the provided urls.
func NewProviderFromURLs(urls ...string) (provider.Interface, error) {
	addrs := []net.Addr{}
	for _, url := range urls {
		addr, err := value.ParseAddrURL(url)
		if err != nil {
			return nil, err
		}
		addrs = append(addrs, addr)
	}
	return NewProvider(addrs...), nil
}

func (p *static) Addrs() (value.Addrs, error) {
	return p.addrs, nil
}

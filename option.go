package netx

type Options struct {
	nodes        []string
	localAddress string
}

type Option func(*Options) error

func Nodes(nodes []string) Option {
	return func(o *Options) error {
		o.nodes = nodes
		return nil
	}
}

func LocalAddress(value string) Option {
	return func(o *Options) error {
		o.localAddress = value
		return nil
	}
}

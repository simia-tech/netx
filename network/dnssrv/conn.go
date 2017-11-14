package dnssrv

import (
	"net"

	"github.com/simia-tech/netx"
	"github.com/simia-tech/netx/value"
)

func init() {
	netx.RegisterDial("dnssrv", Dial)
}

// Dial establishes a connection to the provided address over the provided network.
func Dial(address string, options *value.DialOptions) (net.Conn, error) {
	// if len(options.Nodes) < 1 {
	// 	return nil, fmt.Errorf("no nodes specified")
	// }
	//
	// host, _, err := net.SplitHostPort(address)
	// if err != nil {
	// 	host = address
	// }
	//
	// client := dns.Client{}
	// request := dns.Msg{}
	// request.SetQuestion(fmt.Sprintf("_%s._tcp.service.consul.", host), dns.TypeSRV)
	//
	// response, _, err := client.Exchange(&request, options.Nodes[0])
	// if err != nil {
	// 	return nil, err
	// }
	//
	// names := make(map[string]net.IP)
	// for _, rr := range response.Extra {
	// 	switch rr := rr.(type) {
	// 	case *dns.A:
	// 		names[rr.Hdr.Name] = rr.A
	// 	}
	// }
	//
	// addrs := []net.Addr{}
	// for _, rr := range response.Answer {
	// 	switch rr := rr.(type) {
	// 	case *dns.SRV:
	// 		ip, ok := names[rr.Target]
	// 		if !ok {
	// 			continue
	// 		}
	// 		addrs = append(addrs, &net.TCPAddr{IP: ip, Port: int(rr.Port)})
	// 	}
	// }
	//
	// return netx.DialOne(addrs, options)
	return nil, nil
}

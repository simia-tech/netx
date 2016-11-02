package dnssrv

import (
	"fmt"
	"math/rand"
	"net"

	"github.com/miekg/dns"
	"github.com/pkg/errors"
)

// Dial establishes a connection to the provided address over the provided network.
func Dial(address string, nodes []string) (net.Conn, error) {
	if len(nodes) < 1 {
		return nil, errors.New("no nodes specified")
	}

	host, _, err := net.SplitHostPort(address)
	if err != nil {
		host = address
	}

	client := dns.Client{}
	request := dns.Msg{}
	request.SetQuestion(fmt.Sprintf("_%s._tcp.service.consul.", host), dns.TypeSRV)

	response, _, err := client.Exchange(&request, nodes[0])
	if err != nil {
		return nil, err
	}

	names := make(map[string]net.IP)
	for _, rr := range response.Extra {
		switch rr := rr.(type) {
		case *dns.A:
			names[rr.Hdr.Name] = rr.A
		}
	}

	addrs := []net.Addr{}
	for _, rr := range response.Answer {
		switch rr := rr.(type) {
		case *dns.SRV:
			ip, ok := names[rr.Target]
			if !ok {
				continue
			}
			addrs = append(addrs, &net.TCPAddr{IP: ip, Port: int(rr.Port)})
		}
	}

	switch l := len(addrs); l {
	case 0:
		return nil, errors.Errorf("could find any instances for service [%s]", address)
	case 1:
		return net.Dial("tcp", addrs[0].String())
	default:
		return net.Dial("tcp", addrs[rand.Intn(l)].String())
	}
}

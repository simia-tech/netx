package dnssrv

import (
	"fmt"
	"net"
	"strconv"
	"strings"

	"github.com/miekg/dns"
	"github.com/simia-tech/netx/value"
)

// DNSSRV implements the dnssrv provider.
type DNSSRV struct {
	address         string
	endpointNetwork string
}

// NewProvider returns a new dnssrv provider.
func NewProvider(address, endpointNetwork string) *DNSSRV {
	return &DNSSRV{
		address:         address,
		endpointNetwork: endpointNetwork,
	}
}

// Endpoints returns the endpoints of the provided service.
func (p *DNSSRV) Endpoints(service string) (value.Endpoints, error) {
	client := dns.Client{}
	request := dns.Msg{}
	request.SetQuestion(fmt.Sprintf("_%s._.service.consul.", service), dns.TypeSRV)

	// log.Printf("request\n%s\n", request.String())

	response, _, err := client.Exchange(&request, p.address)
	if err != nil {
		return nil, err
	}
	// log.Printf("response\n%s\n", response)

	hosts := make(map[string]string)
	for _, rr := range response.Extra {
		switch rr := rr.(type) {
		case *dns.A:
			hosts[rr.Hdr.Name] = rr.A.String()
		case *dns.CNAME:
			hosts[rr.Hdr.Name] = strings.TrimSuffix(rr.Target, ".")
		}
	}

	endpoints := value.Endpoints{}
	for _, rr := range response.Answer {
		switch rr := rr.(type) {
		case *dns.SRV:
			host, ok := hosts[rr.Target]
			if !ok {
				continue
			}
			endpoint := value.NewEndpoint(p.endpointNetwork, net.JoinHostPort(host, strconv.Itoa(int(rr.Port))))
			endpoints = append(endpoints, endpoint)
		}
	}

	return endpoints, nil
}

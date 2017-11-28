package client

import (
	"bytes"
	"encoding/json"
	"fmt"
	"math/rand"
	"net"
	"net/http"
	"os"
	"strconv"
	"time"

	"github.com/simia-tech/netx"
	"github.com/simia-tech/netx/filter/blacklist"
	"github.com/simia-tech/netx/provider/static"
	"github.com/simia-tech/netx/selector/roundrobin"
	"github.com/simia-tech/netx/value"
)

// Consul implements a simple client to the consul HTTP api.
type Consul struct {
	client *http.Client
	node   string
}

// NewConsul returns a new consul client.
func NewConsul(urls []string, options ...value.Option) (*Consul, error) {
	p := static.NewProvider()
	for _, url := range urls {
		dial, err := value.ParseEndpointURL(url, options...)
		if err != nil {
			return nil, fmt.Errorf("parsing dial url [%s]: %v", url, err)
		}
		p.Add("consul", dial)
	}

	md, err := netx.NewMultiDialer(p, blacklist.NewFilter(blacklist.ConstantBackoff(100*time.Millisecond)), roundrobin.NewSelector())
	if err != nil {
		return nil, err
	}

	client := &http.Client{
		Transport: netx.NewHTTPMultiTransport(md),
	}

	node, err := os.Hostname()
	if err != nil {
		return nil, err
	}

	return &Consul{
		client: client,
		node:   node,
	}, nil
}

// Register adds a service with the provided name and endpoint at the consul agent.
func (c *Consul) Register(name string, endpoint value.Endpoint) (string, error) {
	host, p, err := net.SplitHostPort(endpoint.Address())
	if err != nil {
		return "", err
	}
	var port int
	if port, err = strconv.Atoi(p); err != nil {
		port = 0
	}

	id := makeID(name + "-")

	m := map[string]interface{}{
		"Node":    c.node,
		"Address": host,
		"Service": map[string]interface{}{
			"ID":      id,
			"Service": name,
			"Tags":    []interface{}{endpoint.Network()},
			"Address": host,
			"Port":    port,
		},
	}

	buffer := &bytes.Buffer{}
	if err = json.NewEncoder(buffer).Encode(m); err != nil {
		return "", err
	}

	request, err := http.NewRequest("PUT", "http://consul/v1/catalog/register", buffer)
	if err != nil {
		return "", err
	}

	response, err := c.client.Do(request)
	if err != nil {
		return "", err
	}
	defer response.Body.Close()

	if response.StatusCode != http.StatusOK {
		return "", fmt.Errorf("expected status code 200, got %d", response.StatusCode)
	}

	return id, nil
}

// Deregister removes the service with the provided id from the consul catalog.
func (c *Consul) Deregister(id string) error {
	m := map[string]interface{}{
		"Node":      c.node,
		"ServiceID": id,
	}

	buffer := &bytes.Buffer{}
	if err := json.NewEncoder(buffer).Encode(m); err != nil {
		return err
	}

	request, err := http.NewRequest("PUT", "http://consul/v1/catalog/deregister", buffer)
	if err != nil {
		return err
	}

	response, err := c.client.Do(request)
	if err != nil {
		return err
	}
	defer response.Body.Close()

	if response.StatusCode != http.StatusOK {
		return fmt.Errorf("expected status code 200, got %d", response.StatusCode)
	}

	return nil
}

// Service fetches all endpoints for the provided service name from the consul catalog.
func (c *Consul) Service(name string) (value.Endpoints, error) {
	response, err := c.client.Get(fmt.Sprintf("http://consul/v1/catalog/service/%s", name))
	if err != nil {
		return nil, err
	}
	defer response.Body.Close()

	body := []interface{}{}
	if err := json.NewDecoder(response.Body).Decode(&body); err != nil {
		return nil, err
	}

	endpoints := value.Endpoints{}
	for _, entry := range body {
		entryMap := entry.(map[string]interface{})
		network := "tcp"
		tags := entryMap["ServiceTags"].([]interface{})
		if len(tags) > 0 {
			network = tags[0].(string)
		}
		address := fmt.Sprintf("%s:%v", entryMap["ServiceAddress"], entryMap["ServicePort"])

		endpoints = append(endpoints, value.NewEndpoint(network, address))
	}

	return endpoints, nil
}

func makeID(prefix string) string {
	bytes := [8]byte{}
	if _, err := rand.Read(bytes[:]); err != nil {
		panic(err)
	}
	return fmt.Sprintf("%s%x", prefix, bytes[:])
}

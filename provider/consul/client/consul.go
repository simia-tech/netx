package client

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"math/rand"
	"net"
	"net/http"
	"strconv"

	"github.com/simia-tech/netx"
	"github.com/simia-tech/netx/provider/static"
	"github.com/simia-tech/netx/selector/roundrobin"
	"github.com/simia-tech/netx/value"
)

type Consul struct {
	client *http.Client
	node   string
}

func NewConsul(urls []string, options ...value.DialOption) (*Consul, error) {
	p := static.NewProvider()
	for _, url := range urls {
		dial, err := value.ParseDialURL(url, options...)
		if err != nil {
			return nil, fmt.Errorf("parsing dial url [%s]: %v", url, err)
		}
		p.Add("consul", dial)
	}

	md, err := netx.NewMultiDialer(p, roundrobin.NewSelector())
	if err != nil {
		return nil, err
	}

	client := &http.Client{
		Transport: netx.NewHTTPMultiTransport(md),
	}

	return &Consul{
		client: client,
	}, nil
}

func (c *Consul) Register(name string, addr net.Addr) (string, error) {
	host, p, err := net.SplitHostPort(addr.String())
	if err != nil {
		return "", err
	}
	port := 0
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
			"Address": host,
			"Port":    port,
		},
	}
	if port > 0 {
		m["Service"].(map[string]interface{})["Port"] = port
	}

	buffer := &bytes.Buffer{}
	if err = json.NewEncoder(buffer).Encode(m); err != nil {
		return "", err
	}

	request, err := http.NewRequest("PUT", "http://consul/v1/catalog/register", buffer)
	if err != nil {
		return "", err
	}

	response, err := http.DefaultClient.Do(request)
	if err != nil {
		return "", err
	}
	defer response.Body.Close()

	if response.StatusCode != http.StatusOK {
		return "", fmt.Errorf("expected status code 200, got %d", response.StatusCode)
	}

	return id, nil
}

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

	response, err := http.DefaultClient.Do(request)
	if err != nil {
		return err
	}
	defer response.Body.Close()

	if response.StatusCode != http.StatusOK {
		return fmt.Errorf("expected status code 200, got %d", response.StatusCode)
	}

	return nil
}

func (c *Consul) Service(name string) ([]net.Addr, error) {
	response, err := http.DefaultClient.Get(fmt.Sprintf("http://consul/v1/catalog/service/%s", name))
	if err != nil {
		return nil, err
	}
	defer response.Body.Close()

	body := []interface{}{}
	if err := json.NewDecoder(response.Body).Decode(&body); err != nil {
		return nil, err
	}

	addrs := []net.Addr{}
	for _, entry := range body {
		entryMap := entry.(map[string]interface{})
		address := fmt.Sprintf("%s:%v", entryMap["ServiceAddress"], entryMap["ServicePort"])
		addr, err := net.ResolveTCPAddr("tcp", address)
		if err != nil {
			log.Printf("could not resolve address [%s]: %v", address, err)
			continue
		}
		addrs = append(addrs, addr)
	}

	return addrs, nil
}

func makeID(prefix string) string {
	bytes := [8]byte{}
	if _, err := rand.Read(bytes[:]); err != nil {
		panic(err)
	}
	return fmt.Sprintf("%s%x", prefix, bytes[:])
}

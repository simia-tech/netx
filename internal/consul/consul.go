package consul

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"math/rand"
	"net"
	"net/http"
	"net/url"
	"os"
	"strconv"

	"github.com/pkg/errors"
)

type consul struct {
	address string
	node    string
}

func newConsulFrom(network string) (*consul, string, error) {
	url, err := url.Parse(network)
	if err != nil {
		return nil, "", errors.Wrapf(err, "parsing network url [%s] failed", network)
	}

	localAddress := ":0"
	if user := url.User; user != nil && user.Username() != "" {
		localAddress = user.Username()
		if p, ok := user.Password(); ok {
			localAddress = fmt.Sprintf("%s:%s", localAddress, p)
		} else {
			localAddress = fmt.Sprintf("%s:0", localAddress)
		}
	}

	node, _ := os.Hostname()
	if value := url.Query().Get("node"); value != "" {
		node = value
	}

	return &consul{
		address: url.Host,
		node:    node,
	}, localAddress, nil
}

func (c *consul) register(name string, addr net.Addr) (string, error) {
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

	request, err := http.NewRequest("PUT", fmt.Sprintf("http://%s/v1/catalog/register", c.address), buffer)
	if err != nil {
		return "", err
	}

	response, err := http.DefaultClient.Do(request)
	if err != nil {
		return "", err
	}
	defer response.Body.Close()

	if response.StatusCode != http.StatusOK {
		return "", errors.Errorf("expected status code 200, got %d", response.StatusCode)
	}

	return id, nil
}

func (c *consul) deregister(id string) error {
	m := map[string]interface{}{
		"Node":      c.node,
		"ServiceID": id,
	}

	buffer := &bytes.Buffer{}
	if err := json.NewEncoder(buffer).Encode(m); err != nil {
		return err
	}

	request, err := http.NewRequest("PUT", fmt.Sprintf("http://%s/v1/catalog/deregister", c.address), buffer)
	if err != nil {
		return err
	}

	response, err := http.DefaultClient.Do(request)
	if err != nil {
		return err
	}
	defer response.Body.Close()

	if response.StatusCode != http.StatusOK {
		return errors.Errorf("expected status code 200, got %d", response.StatusCode)
	}

	return nil
}

func (c *consul) service(name string) ([]net.Addr, error) {
	response, err := http.DefaultClient.Get(fmt.Sprintf("http://%s/v1/catalog/service/%s", c.address, name))
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

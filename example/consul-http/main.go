package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"time"

	"github.com/simia-tech/netx"
	"github.com/simia-tech/netx/filter/blacklist"
	"github.com/simia-tech/netx/provider/consul"
	"github.com/simia-tech/netx/selector/roundrobin"
)

// This example requires a consul node to run at localhost:8500. Running `consul agent -dev` should do the job.
func main() {
	listener, err := netx.Listen("tcp", "localhost:0")
	if err != nil {
		log.Fatal(err)
	}

	provider, err := consul.NewProvider([]string{"tcp://localhost:8500"})
	if err != nil {
		log.Fatal(err)
	}
	provider.AddEndpointAddr("greeter", listener.Addr())

	mux := &http.ServeMux{}
	mux.HandleFunc("/hello", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, "Hello")
	})

	server := &http.Server{Handler: mux}
	go func() {
		server.Serve(listener)
	}()

	md, err := netx.NewMultiDialer(provider, blacklist.NewFilter(blacklist.ConstantBackoff(100*time.Millisecond)), roundrobin.NewSelector(), time.Second)
	if err != nil {
		log.Fatal(err)
	}
	client := &http.Client{
		Transport: netx.NewHTTPMultiTransport(md),
	}
	response, err := client.Get("http://greeter/hello")
	if err != nil {
		log.Fatal(err)
	}
	defer response.Body.Close()

	body, err := ioutil.ReadAll(response.Body)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(string(body))
	// Output: Hello
}

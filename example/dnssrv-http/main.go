package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/http"

	"github.com/simia-tech/netx"
	_ "github.com/simia-tech/netx/network/consul"
	_ "github.com/simia-tech/netx/network/dnssrv"
)

// This example requires a consul node to run at localhost:8500 (http) and localhost:8600 (dns).
// Running `consul agent -dev` should do the job.
func main() {
	listener, err := netx.Listen("consul", "greeter", netx.Nodes("http://localhost:8500"), netx.PublicAddress("localhost:0"))
	if err != nil {
		log.Fatal(err)
	}

	mux := &http.ServeMux{}
	mux.HandleFunc("/hello", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, "Hello")
	})

	server := &http.Server{Handler: mux}
	go func() {
		server.Serve(listener)
	}()

	client := &http.Client{
		Transport: netx.NewHTTPTransport("dnssrv", netx.Nodes("localhost:8600")),
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

package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/http"

	"github.com/simia-tech/netx"
	_ "github.com/simia-tech/netx/network/nats"
	"github.com/simia-tech/netx/value"
)

// This example requires a NATS node to run at localhost:4222. Running `gnatsd -D` should do the job.
func main() {
	listener, err := netx.Listen("nats", "greeter", value.Nodes("nats://127.0.0.1:4222"))
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
		Transport: netx.NewHTTPTransport("nats", value.Nodes("nats://127.0.0.1:4222")),
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

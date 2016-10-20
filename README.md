# netx - Semantic addressing extention for go's net package

This package provides an extention of go stdlib's net package. It provides extended `Listen` and `Dial` methods
in order to enabled clients and servers for semantic addressing.

The following examples require a local NATS(http://nats.io) node on port 4222.

## TCP connection example

```go
import (
  "fmt"

  "github.com/simia-tech/netx"
)

func main() {
  listener, _ := netx.Listen("nats://localhost:4222", "echo")
	go func() {
		conn, _ := listener.Accept()
		defer conn.Close()

		buffer := make([]byte, 5)
		conn.Read(buffer)
		conn.Write(buffer)
	}()

	client, _ := netx.Dial("nats://localhost:4222", "echo")
	defer client.Close()

	fmt.Fprintf(client, "hello")

	buffer := make([]byte, 5)
	client.Read(buffer)

	fmt.Println(string(buffer))
	// Output: hello
}
```

## HTTP connection example

```go
import (
  "http"

  "github.com/simia-tech/netx"
  "github.com/simia-tech/netx/httpx"
)

func main() {
  listener, _ := netx.Listen("nats://localhost:4222", "greeter")

	mux := &http.ServeMux{}
	mux.HandleFunc("/hello", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, "Hello")
	})

	server := &http.Server{Handler: mux}
	go func() {
		server.Serve(listener)
	}()

	client := &http.Client{Transport: httpx.NewTransport("nats://localhost:4222")}
	response, _ := client.Get("http://greeter/hello")
	defer response.Body.Close()

	body, _ := ioutil.ReadAll(response.Body)
	fmt.Println(string(body))
	// Output: Hello
}
```

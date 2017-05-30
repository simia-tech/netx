# netx - Semantic addressing extention for go's net package

[![GoDoc](https://godoc.org/github.com/simia-tech/netx?status.svg)](https://godoc.org/github.com/simia-tech/netx) [![Build Status](https://travis-ci.org/simia-tech/netx.svg?branch=master)](https://travis-ci.org/simia-tech/netx)

This package provides an extention of go stdlib's net package. It provides extended `Listen` and `Dial` methods
in order to enabled clients and servers for semantic addressing. The returned structs implement `net.Listener` and
`net.Conn` and should seamlessly integrate with your existing application.

For transport/service organisation, [NATS](http://nats.io), [consul](https://consul.io) or DNSSRV can be used. An
implementation of quic is in development.

The following examples require a local [NATS](http://nats.io) node on port 4222.

## TCP connection example

```go
import (
  "fmt"

  "github.com/simia-tech/netx"
  _ "github.com/simia-tech/netx/network/nats"
)

func main() {
  listener, _ := netx.Listen("nats", "echo", netx.Nodes("nats://localhost:4222"))
  go func() {
    conn, _ := listener.Accept()
    defer conn.Close()

    buffer := make([]byte, 5)
    conn.Read(buffer)
    conn.Write(buffer)
  }()

  client, _ := netx.Dial("nats", "echo", netx.Nodes("nats://localhost:4222"))
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
  "net/http"
  "fmt"
  "io/ioutil"

  "github.com/simia-tech/netx"
  _ "github.com/simia-tech/netx/network/nats"
)

func main() {
  listener, _ := netx.Listen("nats", "greeter", netx.Nodes("nats://localhost:4222"))

  mux := &http.ServeMux{}
  mux.HandleFunc("/hello", func(w http.ResponseWriter, r *http.Request) {
    fmt.Fprintf(w, "Hello")
  })

  server := &http.Server{Handler: mux}
  go func() {
    server.Serve(listener)
  }()

  client := &http.Client{
    Transport: netx.NewHTTPTransport("nats", netx.Nodes("nats://localhost:4222")),
  }
  response, _ := client.Get("http://greeter/hello")
  defer response.Body.Close()

  body, _ := ioutil.ReadAll(response.Body)
  fmt.Println(string(body))
  // Output: Hello
}
```

## More examples

More example can be found in the [examples](https://github.com/simia-tech/netx/example) directory.

## Tests

In order to run the tests, type

    go test -v ./...

# netx - Semantic addressing extention for go's net package

This package provides an extention of go stdlib's net package. It provides extended `Listen` and `Dial` methods
in order to enabled clients and servers for semantic addressing.

## Code

```go
listener, _ := netx.Listen("nats://localhost:4222", "echo")
go func () {
  conn, _ := listener.Accept()

  buffer := make([]byte, 5)
  conn.Read(buffer)
}()

client, _ := netx.Dial("nats://localhost:4222", "echo")

fmt.Fprintf(client, "hello")
```

package main

import (
	"context"
	"fmt"
	"log"
	"time"

	"google.golang.org/grpc"

	"github.com/simia-tech/netx"
	_ "github.com/simia-tech/netx/network/nats"
	"github.com/simia-tech/netx/value"
)

type echoServer struct{}

func (e *echoServer) Echo(ctx context.Context, request *EchoRequest) (*EchoResponse, error) {
	return &EchoResponse{Text: request.Text}, nil
}

// This example requires a NATS node to run at localhost:4222. Running `gnatsd -D` should do the job.
func main() {
	listener, err := netx.Listen("nats", "echo", value.Nodes("nats://127.0.0.1:4222"))
	if err != nil {
		log.Fatal(err)
	}

	server := grpc.NewServer()
	RegisterEchoServiceServer(server, &echoServer{})

	go func() {
		server.Serve(listener)
	}()

	conn, err := grpc.Dial("echo",
		grpc.WithDialer(netx.NewGRPCDialer("nats", value.Nodes("nats://127.0.0.1:4222"))),
		grpc.WithInsecure())
	if err != nil {
		log.Fatal(err)
	}

	client := NewEchoServiceClient(conn)

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	response, err := client.Echo(ctx, &EchoRequest{Text: "Hello"})
	cancel()
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println(response.Text)
	// Output: Hello
}

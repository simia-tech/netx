package main

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/tls"
	"crypto/x509"
	"encoding/pem"
	"fmt"
	"log"
	"math/big"
	"net"
	"time"

	"github.com/simia-tech/netx"
	"github.com/simia-tech/netx/filter/blacklist"
	_ "github.com/simia-tech/netx/network/quic"
	"github.com/simia-tech/netx/provider/static"
	"github.com/simia-tech/netx/selector/roundrobin"
	"github.com/simia-tech/netx/value"
)

func main() {
	addrOneChan := runEchoServer("127.0.0.1:0")
	addrTwoChan := runEchoServer("127.0.0.1:0")

	provider := static.NewProvider()
	provider.Add("test", value.NewEndpointFromAddr(<-addrOneChan, value.TLS(&tls.Config{InsecureSkipVerify: true})))
	provider.Add("test", value.NewEndpointFromAddr(<-addrTwoChan, value.TLS(&tls.Config{InsecureSkipVerify: true})))

	md, err := netx.NewMultiDialer(provider, blacklist.NewFilter(blacklist.ConstantBackoff(100*time.Millisecond)), roundrobin.NewSelector())
	if err != nil {
		log.Fatal(err)
	}

	conn, err := md.Dial("test")
	if err != nil {
		log.Fatal(err)
	}

	fmt.Fprintln(conn, "hello")

	buffer := make([]byte, 6)
	conn.Read(buffer)
	fmt.Print(string(buffer))
}

func runEchoServer(address string) chan net.Addr {
	ch := make(chan net.Addr)
	go func() {
		listener, err := netx.Listen("quic", address, value.TLS(generateTLSConfig()))
		if err != nil {
			log.Fatal(err)
		}
		ch <- listener.Addr()

		conn, err := listener.Accept()
		if err != nil {
			log.Fatal(err)
		}

		buffer := make([]byte, 6)
		conn.Read(buffer)
		conn.Write(buffer)
	}()
	return ch
}

func generateTLSConfig() *tls.Config {
	key, err := rsa.GenerateKey(rand.Reader, 1024)
	if err != nil {
		panic(err)
	}
	template := x509.Certificate{SerialNumber: big.NewInt(1)}
	certDER, err := x509.CreateCertificate(rand.Reader, &template, &template, &key.PublicKey, key)
	if err != nil {
		panic(err)
	}
	keyPEM := pem.EncodeToMemory(&pem.Block{Type: "RSA PRIVATE KEY", Bytes: x509.MarshalPKCS1PrivateKey(key)})
	certPEM := pem.EncodeToMemory(&pem.Block{Type: "CERTIFICATE", Bytes: certDER})

	tlsCert, err := tls.X509KeyPair(certPEM, keyPEM)
	if err != nil {
		panic(err)
	}
	return &tls.Config{Certificates: []tls.Certificate{tlsCert}}
}

package quic_test

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/tls"
	"crypto/x509"
	"encoding/pem"
	"math/big"
	"testing"

	_ "github.com/simia-tech/netx/network/quic"
	"github.com/simia-tech/netx/test"
	"github.com/simia-tech/netx/value"
)

var options = &test.Options{
	ListenNetwork: "quic",
	ListenAddress: "localhost:0",
	DialNetwork:   "quic",
}

func init() {
	options.ListenOptions = []value.Option{value.TLS(generateTLSConfig())}
	options.DialOptions = []value.Option{value.TLS(&tls.Config{InsecureSkipVerify: true})}
}

func TestConnection(t *testing.T) {
	test.ConnectionTest(t, options)
}

func BenchmarkConnection(b *testing.B) {
	test.ConnectionBenchmark(b, options)
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

package quic_test

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/tls"
	"crypto/x509"
	"crypto/x509/pkix"
	"encoding/pem"
	"math/big"
	"net"
	"testing"
	"time"

	"github.com/simia-tech/netx"
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
	tlsConfig, err := generateTLSConfig([]string{"localhost"}, x509.ExtKeyUsageServerAuth)
	if err != nil {
		panic(err)
	}
	options.ListenOptions = []netx.Option{netx.TLS(tlsConfig)}

	tlsConfig, err = generateTLSConfig(nil, x509.ExtKeyUsageClientAuth)
	if err != nil {
		panic(err)
	}
	options.DialOptions = []value.DialOption{value.TLS(tlsConfig)}
}

func TestConnection(t *testing.T) {
	t.SkipNow()
	test.ConnectionTest(t, options)
}

func BenchmarkConnection(b *testing.B) {
	b.SkipNow()
	test.ConnectionBenchmark(b, options)
}

func generateTLSConfig(hosts []string, usages ...x509.ExtKeyUsage) (*tls.Config, error) {
	privateKey, err := rsa.GenerateKey(rand.Reader, 2048)
	if err != nil {
		return nil, err
	}

	serialNumberLimit := new(big.Int).Lsh(big.NewInt(1), 128)
	serialNumber, err := rand.Int(rand.Reader, serialNumberLimit)
	if err != nil {
		return nil, err
	}

	now := time.Now()
	template := x509.Certificate{
		SerialNumber:          serialNumber,
		Subject:               pkix.Name{Organization: []string{"Posteo.de e.K."}},
		NotBefore:             now,
		NotAfter:              now.Add(356 * 24 * time.Hour),
		KeyUsage:              x509.KeyUsageKeyEncipherment | x509.KeyUsageDigitalSignature,
		ExtKeyUsage:           usages,
		BasicConstraintsValid: true,
	}

	for _, h := range hosts {
		if ip := net.ParseIP(h); ip != nil {
			template.IPAddresses = append(template.IPAddresses, ip)
		} else {
			template.DNSNames = append(template.DNSNames, h)
		}
	}

	derBytes, err := x509.CreateCertificate(rand.Reader, &template, &template, &privateKey.PublicKey, privateKey)
	if err != nil {
		return nil, err
	}

	certificateBytes := pem.EncodeToMemory(&pem.Block{
		Type:  "CERTIFICATE",
		Bytes: derBytes,
	})
	privateKeyBytes := pem.EncodeToMemory(&pem.Block{
		Type:  "RSA PRIVATE KEY",
		Bytes: x509.MarshalPKCS1PrivateKey(privateKey),
	})

	certificate, err := tls.X509KeyPair(certificateBytes, privateKeyBytes)
	if err != nil {
		return nil, err
	}

	return &tls.Config{
		Certificates:       []tls.Certificate{certificate},
		InsecureSkipVerify: true,
	}, nil
}

var serverCertificate = `-----BEGIN CERTIFICATE-----
MIIDBzCCAe+gAwIBAgIQVJqVc3Dd8b0HgU5SL/u+6zANBgkqhkiG9w0BAQsFADAZ
MRcwFQYDVQQKEw5Qb3N0ZW8uZGUgZS5LLjAeFw0xNzAzMDMwOTM5MDVaFw0xODAy
MjIwOTM5MDVaMBkxFzAVBgNVBAoTDlBvc3Rlby5kZSBlLksuMIIBIjANBgkqhkiG
9w0BAQEFAAOCAQ8AMIIBCgKCAQEAwsM1oI676zFMCIqDqEGMmRITNaspg+/9PRIg
w1+8aUn4ucQZnnxhdPTFPaIRtrwVPcBudcRYT7vkM3I2Efv1NPIm87jZ4I8n/+pk
9Ve4rUZWC1OrjdcGqiPiY89VjwBjMmlaDs1fzcaw4XEzMg7Pqvs0iuEsvTEjK1pe
vxICYWt5iILSWz5tliljqUWA5pv0DOBMKOd7nEH6Ue6Ei1FJfYadowfV4Fhd7+OI
Wpqqiv3AVy4f/Nf3GSkV6GDx/DmFQiV2fPW/0AI1nthPXrl8e3fl5da3NJ+Gv7Yk
jIxS1vkMT1reGzvNkNFop2ahcaJu5pba5xzSHg+HfAP1m+48KwIDAQABo0swSTAO
BgNVHQ8BAf8EBAMCBaAwEwYDVR0lBAwwCgYIKwYBBQUHAwEwDAYDVR0TAQH/BAIw
ADAUBgNVHREEDTALgglsb2NhbGhvc3QwDQYJKoZIhvcNAQELBQADggEBAFQQ/lNX
x4kC1InBuhyPkvKx7YF5dTHAIRXSK+SyO1OOLLNZcPofGTbqDnZm5MyZR9lhEIvC
JK70fvjs0BQ2DY0xZaYrmEJwA0VdcBRXZXU/7GdbGbKmv7b9KsRNrjBWY8VXm1Q3
xVwKl4JU91fJ/6m8fdsixuRYGCBYXeYZU705jPyhrxPbXBXvmPV4ZeFxO694XWnY
WssOxR25Ketu3pJt819UoXo7ZeEYrhxpXUqIJikIj3j1L+gIrdH77GywbhbgfWA1
4JRRNpvMI/zvjw+BWygRCLga9sPF1AY4AJDYjJbMvSZrgdCGTiAOVtBvyMsdM5We
hqvJk0H7k+9xr2w=
-----END CERTIFICATE-----
`

var serverPrivateKey = `-----BEGIN RSA PRIVATE KEY-----
MIIEowIBAAKCAQEAwsM1oI676zFMCIqDqEGMmRITNaspg+/9PRIgw1+8aUn4ucQZ
nnxhdPTFPaIRtrwVPcBudcRYT7vkM3I2Efv1NPIm87jZ4I8n/+pk9Ve4rUZWC1Or
jdcGqiPiY89VjwBjMmlaDs1fzcaw4XEzMg7Pqvs0iuEsvTEjK1pevxICYWt5iILS
Wz5tliljqUWA5pv0DOBMKOd7nEH6Ue6Ei1FJfYadowfV4Fhd7+OIWpqqiv3AVy4f
/Nf3GSkV6GDx/DmFQiV2fPW/0AI1nthPXrl8e3fl5da3NJ+Gv7YkjIxS1vkMT1re
GzvNkNFop2ahcaJu5pba5xzSHg+HfAP1m+48KwIDAQABAoIBAQCTNHRN/FPwefwb
4IWOIK0h7NV0FQB15mgjRpZY1P1VH0HNNnienygR/AcwhFSbJyFu4QkcHeEZssvS
TpqrkSJOeFAUmjyjS9BRz1mrTjgZTzYtrXnW5RT2d9Te15C1Wvee3J3i3jtxDqIW
YzbsuOEhPGkEJxlnjcIlPFKsW/JMqltKj6h+q3lWEoRBnR1c5vtxiy9nN3HB9x5b
9S0c4m/gT1QehB9CLKeeM3MKvS6SLlvKB+ciiB/ojMQrSwdTnOFqQZRyTT+N8Nyl
q/IicfsPz3MhqEsATO1W/g6oU7/ro7R6cZFpixAaOoNhKIDzYNknnQMRo8nI8jbZ
FUhfu8kxAoGBANyXgFG24VYQP4Va1bCHUEe5ASjpr/44dn+4T4Jmvt5AU+I2S86A
IokXOVsmEDcoWU878ezi1Ar2hC7Sa7U6H6Kqke9YFeD8kiteFju71raToKIouHU0
fwwM8KR6nFRpgXKeoYY7gBw14JZJe1PluPmualFeY+GHHbLmalPHOLU9AoGBAOIG
VszZoyQlZ2lsJYrW8n1ODzJd7pUlUUf08e5mb/lr/ApVR7RyJO6+qdGX5HvX/wzf
eNjiQ9RmTWQ84N8eN6tMExXBS7sWSnmKoPsTOxl/JDldnQo2VLs32DX6LBE+Li5G
PSSByKVjsUFcMY70mXEKiccdVXSfRpbj8ujt5d2HAoGAe79t1+ltHde9mrTrXb17
FtE9SlNWTJomN94QlInJF2sk46kr+6s8NIXXwj9nJ1o2R9HMFOTmqUPDwXR/wcna
h2mCtq9GjtGBulxswpPMjt3gZjfLysxpXTxBHzQ9UMljOgatfF0SsEC1Pfn4+obL
z9HhRHMGUBkJKYIPmcjtMrs4U2fDB1ZU3SwnaWjICDn7BNTi1HFAjiGpl7ezxldR
LJeFhREZ7/CLr1c525QMelG6wTJW944tWA8mG+jPnxiZioSXjTpdH3pwHZ3dAwUm
3SFRewKBgA6ukz2eFrIZ2N7/Xy1bbxS70ww0BN56FwRH2D7DiN2XWxir3R2GyxXv
JpZBYxs8H5SQk6Viqk9d8TctZ1pWYkn5hCukdtDBUCnzxH9egwHNhRN8aMCh2N+q
2l9vveIb0/TdiqaFopmNEEYBtcxrwgOD+h0EAULbEx7u3/VUpek+
-----END RSA PRIVATE KEY-----
`

var clientCertificate = `-----BEGIN CERTIFICATE-----
MIIC8TCCAdmgAwIBAgIQIWMu6t4muX21wq4gBA/9PjANBgkqhkiG9w0BAQsFADAZ
MRcwFQYDVQQKEw5Qb3N0ZW8uZGUgZS5LLjAeFw0xNzAzMDMwOTM5MDVaFw0xODAy
MjIwOTM5MDVaMBkxFzAVBgNVBAoTDlBvc3Rlby5kZSBlLksuMIIBIjANBgkqhkiG
9w0BAQEFAAOCAQ8AMIIBCgKCAQEAyEV5zzX8qHxLMg+gvzSeeh+95qSy8K4br5na
S1asos7k1nLy8P81zgtSCHgNGuAApDnwtXOnlIIdY3BDLB2owJWL8hi/tnVg3oSr
0O9rxhvcCukhRWPJ79xdIOPkiGrDPVoB6PcgLdge29l61d3ctDY/OI0S1CAXA7ZU
TRmWAsK/s35HaMU+CvLsfMGWtyjXpiaHIsPKdOeFTpX3Rts83NbDaj3YH+r0zvhL
nXRZU/+rjSdeJyPlSBl29TB6xWndXLU/E3SBvV6UuY9wCESWU5Wa8ZA7+iIre098
FW00iAzYYv6s2cHtqlQXlzoS2ut4CIJtKCfg3k5wdVeXv+LX4QIDAQABozUwMzAO
BgNVHQ8BAf8EBAMCBaAwEwYDVR0lBAwwCgYIKwYBBQUHAwIwDAYDVR0TAQH/BAIw
ADANBgkqhkiG9w0BAQsFAAOCAQEAtsqHFONQmKV8QodSSQBjbFLXuIoLDY/OLlkS
16ozbllUIiE82ejQXRH0BRCGTwoGC2V+ZlX4HhGtaJkg6dGRzBgmIVTHDhEHAudr
fROESny+TVZyWQkzca1rRpvDFMrXehLgiAliQ5mpLq1yinwSt+lx/m53U7dOrLxQ
F8XGruOiulOWA+U16sr9GJqRmHz2fFhodsvyYczr2pC44CpX/ima2OoDTQ8gLVMC
xeQ5XdVWkz2zDkwtY3fmbxww4JYWkIphohCxjwxTAp895njaOyBoMCLY2F0OUA5Z
nUXJuDrIMjoeEqVfpdwT8aoGUj9t1ob6vhJ5/V6SDPx4SBHNWg==
-----END CERTIFICATE-----
`

var clientPrivateKey = `-----BEGIN RSA PRIVATE KEY-----
MIIEowIBAAKCAQEAyEV5zzX8qHxLMg+gvzSeeh+95qSy8K4br5naS1asos7k1nLy
8P81zgtSCHgNGuAApDnwtXOnlIIdY3BDLB2owJWL8hi/tnVg3oSr0O9rxhvcCukh
RWPJ79xdIOPkiGrDPVoB6PcgLdge29l61d3ctDY/OI0S1CAXA7ZUTRmWAsK/s35H
aMU+CvLsfMGWtyjXpiaHIsPKdOeFTpX3Rts83NbDaj3YH+r0zvhLnXRZU/+rjSde
JyPlSBl29TB6xWndXLU/E3SBvV6UuY9wCESWU5Wa8ZA7+iIre098FW00iAzYYv6s
2cHtqlQXlzoS2ut4CIJtKCfg3k5wdVeXv+LX4QIDAQABAoIBAQCCcS4blrBQFVTs
8Fzk7SLjrFUGKvQZ621A/NbAB14+VL+cJhayfafP4nO17Gusu4XvcYWkygW2puD5
aZ921oKZnU4fi5si/mTvyj4WwjrSOEckHCB9j7dxsXg++fBaIh+0kDF6Sav98Cx5
SbfGINGl10dqUAiZfaU/17wk06tvc7yJ2w/n3TZksbxBJnRQqnM4WSPWrda6qpRl
Ca/Cqu6ORWZsjSjI7qPk2wJpdU5m+y9YIg1SkI3U4PE3hNvvcWzZXS7P5/SN7LFi
0y6Rl91eSE0ZIdIAxFpT9PezRGh2HeiU4kAa8GMkU4qXPEI1MkTAGB4NMM+Le5bH
FqWM9p4BAoGBAMotdSgR8zfGfPZYc5v64iXkmEHSYkj599rlMQ4dVSoXB8XFb2PT
/nLpcTeCPTHrqndQydOYLwToAjSAXM/+xuAT/S25vUB2pxmf/1obtmTeUuMoj2GO
2XESdKT5HwMJmTnci1FF6i8/zceyUgEkKFO5n//W/zVnV4h2jgenOTnxAoGBAP2W
HFnkHBPPUjMzadumnAAZi7Pqvanf+HacBdM41i19RxaUrSAYHFxbS/9pspQvPb2g
Ke6QvpYxP1hIoFuMzX8SkDiEptsMGpJ0xzE7JCfNMdt13Q367cpxyHRk1f9XblyG
vth4Q6iOx7aD12ZsPLgUOdGMOhfuBxjqwoe9DAzxAoGAWAXEjSaLgswLGeHWq6Fm
FmNZGscy/Vy/WXERk3iX3JRcUPGtloP0sykJnsY4SGS3Oe1VgacvSW6NjzgXsILX
KTXqs567U7aU9+Yd8ahBF9dntPiyvCHKb50+ZZkEtHjYWkW37jGHTPz1Za0wYMjS
OemGTIfZYvHUPViIa7KVirECgYBziApOofBwzgmjLg9SdSupl/nf9FiIpnOqhhbZ
TpG1k9fpX78oWhPBuA59xQgJHyS/2dKA0A0knDdB34S/cPzGogx202i5b2BDzVRb
B5jHUWMfmyklD2d6zjAHZ1FfzdOH8BPOx6v7hWFTs+lUzoczTnOxFnP5JwawwXPz
J5vv4QKBgCDCNseU02FGRXExglXn64A1EZ6c9a6Z2npceqmwJp2UeZuU0SRflXRP
3qV9YZwnuR6blWXzEwKNZkvnPo6ss3pWaaoJxQSWKCGvYay/HUTB2t6f9Aohk0KT
bbAeYd6Dj6Yl4BTjC46AI9xEes1cXO3FUVH9WHdbDhuQrm5WzPT8
-----END RSA PRIVATE KEY-----
`

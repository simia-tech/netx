package httpx_test

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestHTTP(t *testing.T) {
	addr, counter, tearDown := setUpTestHTTPServer(t)
	defer tearDown()

	client := setUpTestHTTPClient(t)

	response, err := client.Get(fmt.Sprintf("http://%s/test", addr))
	require.NoError(t, err)
	defer response.Body.Close()
	require.Equal(t, http.StatusOK, response.StatusCode)

	body, err := ioutil.ReadAll(response.Body)
	require.NoError(t, err)
	assert.Equal(t, "test", string(body))

	assert.Equal(t, 1, counter())
}

func BenchmarkHTTPSimpleGet(b *testing.B) {
	addr, counter, tearDown := setUpTestHTTPServer(b)
	defer tearDown()

	client := setUpTestHTTPClient(b)
	b.ResetTimer()
	for index := 0; index < b.N; index++ {
		response, err := client.Get(fmt.Sprintf("http://%s/test", addr))
		require.NoError(b, err)
		require.Equal(b, http.StatusOK, response.StatusCode)

		body, err := ioutil.ReadAll(response.Body)
		require.NoError(b, err)
		require.Equal(b, "test", string(body))
		require.NoError(b, response.Body.Close())
	}
	b.StopTimer()

	assert.Equal(b, b.N, counter())
}

// func ExampleNewTransport() {
// 	listener, _ := netx.Listen("nats://localhost:4222", "greeter")
//
// 	mux := &http.ServeMux{}
// 	mux.HandleFunc("/hello", func(w http.ResponseWriter, r *http.Request) {
// 		fmt.Fprintf(w, "Hello")
// 	})
//
// 	server := &http.Server{Handler: mux}
// 	go func() {
// 		server.Serve(listener)
// 	}()
//
// 	client := &http.Client{Transport: httpx.NewTransport("nats://localhost:4222")}
// 	response, _ := client.Get("http://greeter/hello")
// 	defer response.Body.Close()
//
// 	body, _ := ioutil.ReadAll(response.Body)
// 	fmt.Println(string(body))
// 	// Output: Hello
// }

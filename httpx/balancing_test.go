package httpx_test

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/simia-tech/netx"
)

func TestHTTPBalancing(t *testing.T) {
	address := netx.RandomAddress("http-")

	_, counterOne, tearDownOne := setUpTestHTTPServer(t, address)
	defer tearDownOne()
	_, counterTwo, tearDownTwo := setUpTestHTTPServer(t, address)
	defer tearDownTwo()

	client := setUpTestHTTPClient(t)

	for index := 0; index < 4; index++ {
		response, err := client.Get(fmt.Sprintf("http://%s/test", address))
		require.NoError(t, err)
		require.Equal(t, http.StatusOK, response.StatusCode)

		body, err := ioutil.ReadAll(response.Body)
		require.NoError(t, err)
		require.NoError(t, response.Body.Close())

		assert.Equal(t, "test", string(body))
	}

	assert.Equal(t, 4, counterOne()+counterTwo())
}

func BenchmarkHTTPBalancing(b *testing.B) {
	address := netx.RandomAddress("http-")

	_, counterOne, tearDownOne := setUpTestHTTPServer(b, address)
	defer tearDownOne()
	_, counterTwo, tearDownTwo := setUpTestHTTPServer(b, address)
	defer tearDownTwo()

	client := setUpTestHTTPClient(b)

	b.ResetTimer()
	for index := 0; index < b.N; index++ {
		response, err := client.Get(fmt.Sprintf("http://%s/test", address))
		require.NoError(b, err)
		require.Equal(b, http.StatusOK, response.StatusCode)

		body, err := ioutil.ReadAll(response.Body)
		require.NoError(b, err)
		require.NoError(b, response.Body.Close())

		require.Equal(b, "test", string(body))
	}
	b.StopTimer()

	assert.Equal(b, b.N, counterOne()+counterTwo())
}

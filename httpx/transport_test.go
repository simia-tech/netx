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
	addr, tearDown := setUpTestHTTPServer(t)
	defer tearDown()

	client := setUpTestHTTPClient(t)

	response, err := client.Get(fmt.Sprintf("http://%s/test", addr))
	require.NoError(t, err)
	defer response.Body.Close()
	require.Equal(t, http.StatusOK, response.StatusCode)

	body, err := ioutil.ReadAll(response.Body)
	require.NoError(t, err)
	assert.Equal(t, "test", string(body))
}

func BenchmarkHTTPSimpleGet(b *testing.B) {
	addr, tearDown := setUpTestHTTPServer(b)
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
}

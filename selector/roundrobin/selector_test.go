package roundrobin_test

import (
	"sync"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/simia-tech/netx/selector"
	"github.com/simia-tech/netx/selector/roundrobin"
	"github.com/simia-tech/netx/value"
)

func TestRoundRobinSelect(t *testing.T) {
	testCases := []struct {
		name string
		urls []string

		expectFirstError    error
		expectFirstDialURL  string
		expectSecondDialURL string
		expectThirdDialURL  string
	}{
		{"Loop", []string{"tcp://localhost:1000", "tcp://localhost:2000"}, nil, "tcp://localhost:1000", "tcp://localhost:2000", "tcp://localhost:1000"},
		{"Empty", []string{}, selector.ErrNoEndpoint, "", "", ""},
	}

	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {
			endpoints := value.MustParseEndpointURLs(testCase.urls)
			s := roundrobin.NewSelector()

			endpoint, err := s.Select(endpoints)
			if testCase.expectFirstError == nil {
				require.NoError(t, err)
				assert.Equal(t, testCase.expectFirstDialURL, value.EndpointURL(endpoint))

				endpoint, err = s.Select(endpoints)
				require.NoError(t, err)
				assert.Equal(t, testCase.expectSecondDialURL, value.EndpointURL(endpoint))

				endpoint, err = s.Select(endpoints)
				require.NoError(t, err)
				assert.Equal(t, testCase.expectThirdDialURL, value.EndpointURL(endpoint))
			} else {
				assert.Equal(t, testCase.expectFirstError.Error(), err.Error())
			}
		})
	}
}

func TestRoundRobinConcurrentSelect(t *testing.T) {
	endpoints := value.MustParseEndpointURLs([]string{"tcp://localhost:1000", "tcp://localhost:2000"})
	s := roundrobin.NewSelector()

	wg := sync.WaitGroup{}
	for index := 0; index < 5; index++ {
		wg.Add(1)
		go func() {
			for i := 0; i < 20; i++ {
				_, err := s.Select(endpoints)
				require.NoError(t, err)
			}
			wg.Done()
		}()
	}
	wg.Wait()
}

package roundrobin_test

import (
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
		{"Empty", []string{}, selector.ErrNoDial, "", "", ""},
	}

	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {
			dials := value.MustParseDialURLs(testCase.urls)
			s := roundrobin.NewSelector()

			dial, err := s.Select(dials)
			if testCase.expectFirstError == nil {
				require.NoError(t, err)
				assert.Equal(t, testCase.expectFirstDialURL, value.DialURL(dial))

				dial, err = s.Select(dials)
				require.NoError(t, err)
				assert.Equal(t, testCase.expectSecondDialURL, value.DialURL(dial))

				dial, err = s.Select(dials)
				require.NoError(t, err)
				assert.Equal(t, testCase.expectThirdDialURL, value.DialURL(dial))
			} else {
				assert.Equal(t, testCase.expectFirstError.Error(), err.Error())
			}
		})
	}
}

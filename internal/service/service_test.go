package service

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestSHA1HashUrl(t *testing.T) {
	const SECRET = "secret"

	testCases := []struct {
		name     string
		url      string
		expected string
	}{
		{
			name:     "Success testing w plain string",
			url:      "thisshouldbehashed",
			expected: "ceb03b75323a3ae65a351130210396558ade157d",
		},
		{
			name:     "Success testing w url string",
			url:      "somenameurl?param1=val1&param2=val2",
			expected: "a81e64fc67c9db24d1c64c78dfefba035e64a741",
		},
	}

	for i, tc := range testCases {
		t.Run(fmt.Sprintf("%d# %s", i+1, tc.name), func(t *testing.T) {
			out := SHA1HashUrl(SECRET, tc.url)
			assert.Equal(t, tc.expected, out)
		})
	}
}

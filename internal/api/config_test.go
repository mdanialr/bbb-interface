package api

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestSanitization_AssertRequired(t *testing.T) {
	testCases := []struct {
		name   string
		sample Config
		isErr  bool
	}{
		{
			name:  "Empty sample should error because there are required values",
			isErr: true,
		},
		{
			name:   "Empty sample should error because host field is required",
			sample: Config{Secret: "sssttt"},
			isErr:  true,
		},
		{
			name:   "Should error because secret field is required",
			sample: Config{Host: "http://localhost"},
			isErr:  true,
		},
		{
			name:   "Should error because secret field is required",
			sample: Config{Host: "http://localhost"},
			isErr:  true,
		},

		{
			name:   "Should pass because all required fields are filled",
			sample: Config{Host: "http://localhost", Secret: "sstt"},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			err := tc.sample.Sanitization()

			switch tc.isErr {
			case true:
				require.Error(t, err)
			case false:
				require.NoError(t, err)
			}
		})
	}
}

func TestSanitization_AssertEqual(t *testing.T) {
	sample := Config{Secret: "sstt", Host: "http://localhost"}
	expected := Config{Secret: "sstt", Host: "http://localhost/", URL: "http://localhost/bigbluebutton/api"}

	testCases := []struct {
		name     string
		sample   Config
		expected Config
	}{
		{
			name:     "Host field should has trailing slash",
			sample:   sample,
			expected: expected,
		},
		{
			name:     "Url field should be combinations of host field and `bigbluebutton/api/`",
			sample:   sample,
			expected: expected,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			err := tc.sample.Sanitization()
			require.NoError(t, err)

			assert.Equal(t, tc.sample, tc.expected)
		})
	}
}

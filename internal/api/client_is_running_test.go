package api

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestParseIsRunning(t *testing.T) {
	testCases := []struct {
		name    string
		sample  IsRunning
		expect  string
		wantErr bool
	}{
		{
			name:    "Should error if `meeting_id` field is not provided",
			sample:  IsRunning{},
			wantErr: true,
		},
		{
			name:    "Should pass if `meeting_id` field is provided",
			sample:  IsRunning{MeetingId: "meet01"},
			expect:  "/isMeetingRunning?meetingID=meet01",
			wantErr: false,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			out, err := tc.sample.ParseIsRunning()

			switch tc.wantErr {
			case false:
				require.NoError(t, err)
				assert.Equal(t, tc.expect, out)
			case true:
				require.Error(t, err)
			}
		})
	}
}

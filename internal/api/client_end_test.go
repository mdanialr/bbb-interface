package api

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestParseEndMeeting(t *testing.T) {
	testCases := []struct {
		name    string
		sample  EndMeeting
		expect  string
		wantErr bool
	}{
		{
			name:    "Should error if `meeting_id` field is not provided",
			sample:  EndMeeting{Password: "pass"},
			wantErr: true,
		},
		{
			name:    "Should error if `password` field is not provided",
			sample:  EndMeeting{MeetingId: "meet01"},
			wantErr: true,
		},
		{
			name:    "Should error if both `meeting_id` `password` fields are not provided",
			sample:  EndMeeting{},
			wantErr: true,
		},
		{
			name:    "Should pass if all required fields are provided",
			sample:  EndMeeting{MeetingId: "meet02", Password: "pass"},
			expect:  "/end?meetingID=meet02&password=pass",
			wantErr: false,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			out, err := tc.sample.ParseEndMeeting()

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

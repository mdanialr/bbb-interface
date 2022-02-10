package api

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestParseJoinMeeting(t *testing.T) {
	t.Run("Error if required `name` field is not provided", func(t *testing.T) {
		sample := JoinMeeting{}
		_, err := sample.ParseJoinMeeting()
		require.Error(t, err)
	})

	t.Run("Error if required `meeting_id` field is not provided", func(t *testing.T) {
		sample := JoinMeeting{Name: "name"}
		_, err := sample.ParseJoinMeeting()
		require.Error(t, err)
	})

	t.Run("Error if required `password` field is not provided", func(t *testing.T) {
		sample := JoinMeeting{Name: "name", MeetingId: "meet01"}
		_, err := sample.ParseJoinMeeting()
		require.Error(t, err)
	})

	t.Run("Error if required `created_at` field is not provided", func(t *testing.T) {
		sample := JoinMeeting{Name: "name", MeetingId: "meet01", Password: "pass"}
		_, err := sample.ParseJoinMeeting()
		require.Error(t, err)
	})

	t.Run("Pass if all required fields provided", func(t *testing.T) {
		expect := "/join?meetingID=meet01&password=ap&fullName=Fake&createTime=273648"
		sample := JoinMeeting{MeetingId: "meet01", Password: "ap", Name: "Fake", CreateTime: "273648"}
		out, err := sample.ParseJoinMeeting()
		require.NoError(t, err)
		assert.Equal(t, expect, out)
	})
}

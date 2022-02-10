package api

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

type fakeRandStrGenerator struct {
	Length int
}

func (f fakeRandStrGenerator) RandString() (s string) {
	for i := 0; i < f.Length; i++ {
		s += "a"
	}

	return
}

func TestParseCreateMeeting(t *testing.T) {
	fake := fakeRandStrGenerator{Length: 8}

	t.Run("Error if required params not provided", func(t *testing.T) {
		sample := CreateMeeting{AttendeePass: "secret"}
		_, err := sample.ParseCreateMeeting(fake)
		require.Error(t, err)
	})

	t.Run("Should use random meeting ID if not provided by request", func(t *testing.T) {
		sample := CreateMeeting{Name: "test"}
		out, err := sample.ParseCreateMeeting(fake)
		require.NoError(t, err)
		assert.Equal(t, "/create?name=test&meetingID=aaaaaaaa&moderatorPW=aaaaaaaa&attendeePW=aaaaaaaa&logoutURL=", out)
	})

	t.Run("Should use given moderator password if provided by request", func(t *testing.T) {
		sample := CreateMeeting{Name: "te", ModeratorPass: "idm"}
		out, err := sample.ParseCreateMeeting(fake)
		require.NoError(t, err)
		assert.Equal(t, "/create?name=te&meetingID=aaaaaaaa&moderatorPW=idm&attendeePW=aaaaaaaa&logoutURL=", out)
	})

	t.Run("Should use given max participants if provided by request", func(t *testing.T) {
		sample := CreateMeeting{Name: "te", ModeratorPass: "mpw", AttendeePass: "apw", MaxParticipants: 255}
		out, err := sample.ParseCreateMeeting(fake)
		require.NoError(t, err)
		assert.Equal(t, "/create?name=te&meetingID=aaaaaaaa&moderatorPW=mpw&attendeePW=apw&logoutURL=&maxParticipants=255", out)
	})
}

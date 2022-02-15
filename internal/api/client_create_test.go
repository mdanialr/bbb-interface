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
		assert.Equal(t, "/create?name=test&meetingID=aaaaaaaa&moderatorPW=aaaaaaaa&attendeePW=aaaaaaaa", out)
	})

	t.Run("Should use given moderator password if provided by request", func(t *testing.T) {
		sample := CreateMeeting{Name: "te", ModeratorPass: "idm"}
		out, err := sample.ParseCreateMeeting(fake)
		require.NoError(t, err)
		assert.Equal(t, "/create?name=te&meetingID=aaaaaaaa&moderatorPW=idm&attendeePW=aaaaaaaa", out)
	})

	t.Run("Should use given max participants if provided by request", func(t *testing.T) {
		sample := CreateMeeting{Name: "te", ModeratorPass: "mpw", AttendeePass: "apw", MaxParticipants: 255}
		out, err := sample.ParseCreateMeeting(fake)
		require.NoError(t, err)
		assert.Equal(t, "/create?name=te&meetingID=aaaaaaaa&moderatorPW=mpw&attendeePW=apw&maxParticipants=255", out)
	})

	t.Run("Make sure name & welcome params got encoded to URL even using whitespace", func(t *testing.T) {
		sample := CreateMeeting{Name: "use blank", ModeratorPass: "mpw", AttendeePass: "apw", WelcomeMsg: "Hello from earth!!"}
		out, err := sample.ParseCreateMeeting(fake)
		require.NoError(t, err)
		assert.Equal(t, "/create?name=use+blank&meetingID=aaaaaaaa&moderatorPW=mpw&attendeePW=apw&welcome=Hello+from+earth%21%21", out)
	})

	t.Run("Should include `redirect_at_logout` field if provided then encoded to valid URL otherwise no need to include it in URL", func(t *testing.T) {
		sample := CreateMeeting{Name: "hello there", ModeratorPass: "mp", AttendeePass: "ap", RedirectAtLogout: "https://redirect.domain/dashboard user six"}
		out, err := sample.ParseCreateMeeting(fake)
		require.NoError(t, err)
		assert.Equal(t, "/create?name=hello+there&meetingID=aaaaaaaa&moderatorPW=mp&attendeePW=ap&logoutURL=https%3A%2F%2Fredirect.domain%2Fdashboard+user+six", out)
	})

	t.Run("Should include `record` boolean field if its true", func(t *testing.T) {
		sample := CreateMeeting{Name: "meet one", ModeratorPass: "mp", AttendeePass: "ap", IsRecording: true}
		out, err := sample.ParseCreateMeeting(fake)
		require.NoError(t, err)
		assert.Equal(t, "/create?name=meet+one&meetingID=aaaaaaaa&moderatorPW=mp&attendeePW=ap&record=true", out)
	})

	t.Run("Should not include `record` boolean field in URL if its false", func(t *testing.T) {
		sample := CreateMeeting{Name: "meet one", ModeratorPass: "mp", AttendeePass: "ap", IsRecording: false}
		out, err := sample.ParseCreateMeeting(fake)
		require.NoError(t, err)
		assert.Equal(t, "/create?name=meet+one&meetingID=aaaaaaaa&moderatorPW=mp&attendeePW=ap", out)
	})

	t.Run("If `record` boolean field is provided and its true then must also include `autoStartRecording=true` in URL", func(t *testing.T) {
		sample := CreateMeeting{Name: "meet one", ModeratorPass: "mp", AttendeePass: "ap", IsRecording: true}
		out, err := sample.ParseCreateMeeting(fake)
		require.NoError(t, err)
		assert.Equal(t, "/create?name=meet+one&meetingID=aaaaaaaa&moderatorPW=mp&attendeePW=ap&record=true", out)
	})
}

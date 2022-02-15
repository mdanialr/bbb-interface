package api

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestParseJoinMeeting(t *testing.T) {
	testCases := []struct {
		name    string
		sample  JoinMeeting
		expect  string
		wantErr bool
	}{
		{
			name:    "Error if required `name` field is not provided",
			sample:  JoinMeeting{},
			wantErr: true,
		},
		{
			name:    "Error if required `meeting_id` field is not provided",
			sample:  JoinMeeting{Name: "name"},
			wantErr: true,
		},
		{
			name:    "Error if required `password` field is not provided",
			sample:  JoinMeeting{Name: "name", MeetingId: "meet01"},
			wantErr: true,
		},
		{
			name:    "Error if required `created_at` field is not provided",
			sample:  JoinMeeting{Name: "name", MeetingId: "meet01", Password: "pass"},
			wantErr: true,
		},
		{
			name:   "Pass if all required fields provided",
			sample: JoinMeeting{MeetingId: "meet01", Password: "ap", Name: "Fake", CreateTime: "273648"},
			expect: "/join?meetingID=meet01&password=ap&fullName=Fake&createTime=273648",
		},
		{
			name:   "Should pass even the name has whitespace",
			sample: JoinMeeting{MeetingId: "meet01", Password: "ap", Name: "Fake Using Whitespace", CreateTime: "273648"},
			expect: "/join?meetingID=meet01&password=ap&fullName=Fake+Using+Whitespace&createTime=273648",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			out, err := tc.sample.ParseJoinMeeting()
			switch tc.wantErr {
			case true:
				require.Error(t, err)
			case false:
				require.NoError(t, err)
				assert.Equal(t, tc.expect, out)
			}
		})
	}
}

func TestParseJoinMeeting_OptionalParams(t *testing.T) {
	testCases := []struct {
		name   string
		sample JoinMeeting
		expect string
	}{
		{
			name:   "Pass even all optional fields are not provided",
			sample: JoinMeeting{MeetingId: "meet01", Password: "ap", Name: "Fake", CreateTime: "273648"},
			expect: "/join?meetingID=meet01&password=ap&fullName=Fake&createTime=273648",
		},
		{
			name:   "If user ID exist then make sure it exist in url and match",
			sample: JoinMeeting{MeetingId: "meet01", Password: "ap", Name: "Fake", CreateTime: "273648", UserId: "user01"},
			expect: "/join?meetingID=meet01&password=ap&fullName=Fake&createTime=273648&userID=user01",
		},
		{
			name:   "user ID should got encoded to URL even using whitespace if provided",
			sample: JoinMeeting{MeetingId: "meet01", Password: "ap", Name: "Fake", CreateTime: "273648", UserId: "user 02 whitespace"},
			expect: "/join?meetingID=meet01&password=ap&fullName=Fake&createTime=273648&userID=user+02+whitespace",
		},
		{
			name:   "If avatar exist then make sure it exist in url and match and already URL encoded",
			sample: JoinMeeting{MeetingId: "meet01", Password: "ap", Name: "Fake", CreateTime: "273648", Avatar: "https://site.com/avatar.png"},
			expect: "/join?meetingID=meet01&password=ap&fullName=Fake&createTime=273648&avatarURL=https%3A%2F%2Fsite.com%2Favatar.png",
		},
		{
			name:   "If guest exist then make sure it exist in url and match (bool in string: 'true' or 'false')",
			sample: JoinMeeting{MeetingId: "meet01", Password: "ap", Name: "Fake", CreateTime: "273648", IsGuest: true},
			expect: "/join?meetingID=meet01&password=ap&fullName=Fake&createTime=273648&guest=true",
		},
		{
			name:   "If avatar exist then make sure its encoded to URL",
			sample: JoinMeeting{MeetingId: "meet01", Password: "ap", Name: "Fake", CreateTime: "273648", Avatar: "https://site.com/path to avatar!!.png"},
			expect: "/join?meetingID=meet01&password=ap&fullName=Fake&createTime=273648&avatarURL=https%3A%2F%2Fsite.com%2Fpath+to+avatar%21%21.png",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			out, err := tc.sample.ParseJoinMeeting()
			require.NoError(t, err)
			assert.Equal(t, tc.expect, out)
		})
	}
}

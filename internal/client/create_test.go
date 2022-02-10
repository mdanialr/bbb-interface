package client

import (
	"encoding/xml"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gofiber/fiber/v2"
	"github.com/kurvaid/bbb-interface/internal/api"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

const fakeChecksum = "checksum"

type fakeRandStrGenerator struct{ Length int }

func (f fakeRandStrGenerator) RandString() (s string) {
	for i := 0; i < f.Length; i++ {
		s += "a"
	}

	return
}

func TestClientCreateMeeting_AssertUrl(t *testing.T) {
	// prepare fake server to mimic BBB Server
	server := httptest.NewServer(http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {
		expectSentUrl := "/bigbluebutton/api/create?name=meet-one&meetingID=meet01&moderatorPW=pass&attendeePW=pass&logoutURL=&checksum=" + fakeChecksum
		assert.Equal(t, expectSentUrl, req.URL.String())

		name := req.URL.Query().Get("name")
		meetId := req.URL.Query().Get("meetingID")
		attPass := req.URL.Query().Get("attendeePW")
		modPass := req.URL.Query().Get("moderatorPW")

		// assert name must not empty
		require.NotEqual(t, "", name)

		response := api.CreateMeetingResponse{
			MeetingId:     meetId,
			AttendeePass:  attPass,
			ModeratorPass: modPass,
		}

		resp, err := xml.Marshal(&response)
		require.NoError(t, err)

		rw.WriteHeader(fiber.StatusOK)
		_, err = rw.Write(resp)
		require.NoError(t, err)
	}))

	// sample data to mimic json request from client
	sample := api.CreateMeeting{
		Name:          "meet-one",
		MeetingId:     "meet01",
		AttendeePass:  "pass",
		ModeratorPass: "pass",
	}
	out, err := sample.ParseCreateMeeting(fakeRandStrGenerator{Length: 8})
	url := fmt.Sprintf("%s/%s%s", server.URL, api.EndPoint, out)

	fakeAPI := Create{server.Client(), url, fakeChecksum}
	resp, err := fakeAPI.CreateMeeting()
	require.NoError(t, err)

	var respModel api.CreateMeetingResponse
	err = xml.Unmarshal(resp, &respModel)
	require.NoError(t, err)

	assert.Equal(t, sample.MeetingId, respModel.MeetingId)
	assert.Equal(t, sample.AttendeePass, respModel.AttendeePass)
	assert.Equal(t, sample.ModeratorPass, respModel.ModeratorPass)
	t.Cleanup(func() {
		server.Close()
	})
}

func TestClientCreateMeeting_AssertWithoutModeratorAndAttendeePassword(t *testing.T) {
	// prepare fake server to mimic BBB Server
	server := httptest.NewServer(http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {
		name := req.URL.Query().Get("name")
		meetId := req.URL.Query().Get("meetingID")
		attPass := req.URL.Query().Get("attendeePW")
		modPass := req.URL.Query().Get("moderatorPW")

		// assert name must not empty
		require.NotEqual(t, "", name)

		response := api.CreateMeetingResponse{
			MeetingId:     meetId,
			AttendeePass:  attPass,
			ModeratorPass: modPass,
		}

		resp, err := xml.Marshal(&response)
		require.NoError(t, err)

		rw.WriteHeader(fiber.StatusOK)
		_, err = rw.Write(resp)
		require.NoError(t, err)
	}))

	// sample data to mimic json request from client
	sample := api.CreateMeeting{
		Name:      "meet-two",
		MeetingId: "meet02",
	}
	out, err := sample.ParseCreateMeeting(fakeRandStrGenerator{Length: 8})
	url := fmt.Sprintf("%s/%s%s", server.URL, api.EndPoint, out)

	fakeAPI := Create{server.Client(), url, fakeChecksum}
	resp, err := fakeAPI.CreateMeeting()
	require.NoError(t, err)

	var respModel api.CreateMeetingResponse
	err = xml.Unmarshal(resp, &respModel)
	require.NoError(t, err)

	assert.Equal(t, sample.MeetingId, respModel.MeetingId)
	assert.Equal(t, "aaaaaaaa", respModel.AttendeePass)
	assert.Equal(t, "aaaaaaaa", respModel.ModeratorPass)
	t.Cleanup(func() {
		server.Close()
	})
}

func TestClientCreateMeeting_IncreaseCoverage(t *testing.T) {
	// prepare fake server to mimic BBB Server
	server := httptest.NewServer(http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {
		rw.WriteHeader(fiber.StatusOK)
		rw.Write(nil)
	}))

	// sample data to mimic json request from client
	sample := api.CreateMeeting{MeetingId: "meet02"}

	out, err := sample.ParseCreateMeeting(fakeRandStrGenerator{Length: 8})
	url := fmt.Sprintf("%s/%s%s", "http://localhost", api.EndPoint, out)

	fakeAPI := Create{server.Client(), url, fakeChecksum}
	_, err = fakeAPI.CreateMeeting()
	require.Error(t, err)

	t.Cleanup(func() {
		server.Close()
	})
}

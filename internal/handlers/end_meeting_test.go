package handlers

import (
	"bytes"
	"encoding/json"
	"encoding/xml"
	"github.com/stretchr/testify/assert"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gofiber/fiber/v2"
	"github.com/kurvaid/bbb-interface/internal/api"
	"github.com/kurvaid/bbb-interface/internal/config"
	"github.com/stretchr/testify/require"
)

var sampleEndRequest = []string{
	`
{
	"meeting_id": "meet01"
}
`, `
{
	"password": "secret"
}
`, `
{
	"meeting_id": "meet01",
	"password": "pass"
}
`,
}

var sampleEndMeetStdResponse = []api.StdResponse{
	api.StdResponse{
		CodeString: "FAILED",
		MsgKey:     "invalidPassword",
		MsgDetail:  "The supplied moderator password is incorrect",
	},
	api.StdResponse{
		CodeString: "SUCCESS",
		MsgKey:     "sentEndMeetingRequest",
		MsgDetail:  "A request to end the meeting was sent. Please wait a few seconds, and then use the getMeetingInfo or isMeetingRunning API calls to verify that it was ended.",
	},
}

// prepare fake server to mimic BBB Server
var fakeEndMeetServer = func(t *testing.T) *httptest.Server {
	return httptest.NewServer(http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {
		resp := sampleEndMeetStdResponse[0]
		xm, err := xml.Marshal(&resp)
		if err != nil {
			t.Fatalf("failed to marshal object to xml: %s", err)
		}

		rw.WriteHeader(fiber.StatusOK)
		rw.Header().Set("Content-Type", fiber.MIMEApplicationXML)
		_, err = rw.Write(xm)
		require.NoError(t, err)
	}))
}

func TestEndMeeting(t *testing.T) {
	var fakeSuccessEndMeetServer = func(t *testing.T) *httptest.Server {
		return httptest.NewServer(http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {
			resp := sampleEndMeetStdResponse[1]
			xm, err := xml.Marshal(&resp)
			if err != nil {
				t.Fatalf("failed to marshal object to xml: %s", err)
			}

			rw.WriteHeader(fiber.StatusOK)
			rw.Header().Set("Content-Type", fiber.MIMEApplicationXML)
			_, err = rw.Write(xm)
			require.NoError(t, err)
		}))
	}

	conf, err := config.NewConfig(bytes.NewBufferString(sampleConfigFile[0]))
	require.NoError(t, err)
	require.NoError(t, conf.Sanitization())
	conf.BBB.Host = fakeSuccessEndMeetServer(t).URL
	require.NoError(t, conf.BBB.Sanitization())

	app := fiber.New()
	app.Post("/end", EndMeeting(conf, fakeSuccessEndMeetServer(t).Client()))

	t.Run("Success if all required fields are provided", func(t *testing.T) {
		buf := bytes.NewBufferString(sampleEndRequest[2])
		req := httptest.NewRequest(fiber.MethodPost, "/end", buf)
		req.Header.Set("Content-Type", fiber.MIMEApplicationJSON)
		res, err := app.Test(req)
		if err != nil {
			t.Fatalf("failed to initiate app test: %s", err)
		}
		assert.Equal(t, fiber.StatusOK, res.StatusCode)
		jsRes := struct {
			Message string `json:"message"`
		}{}
		if err := json.NewDecoder(res.Body).Decode(&jsRes); err != nil {
			t.Fatalf("failed to decode json in test : %v", err)
		}
		assert.Contains(t, jsRes.Message, "successfully")
	})
}

func TestEndMeeting_VariousErrorPaths(t *testing.T) {
	conf, err := config.NewConfig(bytes.NewBufferString(sampleConfigFile[0]))
	require.NoError(t, err)
	require.NoError(t, conf.Sanitization())

	app := fiber.New()
	app.Post("/end", EndMeeting(conf, fakeEndMeetServer(t).Client()))

	t.Run("Failed if `meeting_id` field is not provided", func(t *testing.T) {
		buf := bytes.NewBufferString(sampleEndRequest[0])
		req := httptest.NewRequest(fiber.MethodPost, "/end", buf)
		req.Header.Set("Content-Type", fiber.MIMEApplicationJSON)
		res, err := app.Test(req)
		if err != nil {
			t.Fatalf("failed to initiate app test: %s", err)
		}
		assert.Equal(t, fiber.StatusBadRequest, res.StatusCode)
	})

	t.Run("Failed if `password` field is not provided", func(t *testing.T) {
		buf := bytes.NewBufferString(sampleEndRequest[1])
		req := httptest.NewRequest(fiber.MethodPost, "/end", buf)
		req.Header.Set("Content-Type", fiber.MIMEApplicationJSON)
		res, err := app.Test(req)
		if err != nil {
			t.Fatalf("failed to initiate app test: %s", err)
		}
		assert.Equal(t, fiber.StatusBadRequest, res.StatusCode)
	})

	t.Run("Failed if request's content type is invalid", func(t *testing.T) {
		buf := bytes.NewBufferString(sampleEndRequest[2])
		req := httptest.NewRequest(fiber.MethodPost, "/end", buf)
		req.Header.Set("Content-Type", "should error")
		res, err := app.Test(req)
		if err != nil {
			t.Fatalf("failed to initiate app test: %s", err)
		}
		assert.Equal(t, fiber.StatusBadRequest, res.StatusCode)
	})

	t.Run("Failed if using invalid/empty BBB Server host, even if already has valid requirements", func(t *testing.T) {
		buf := bytes.NewBufferString(sampleEndRequest[2])
		req := httptest.NewRequest(fiber.MethodPost, "/end", buf)
		req.Header.Set("Content-Type", fiber.MIMEApplicationJSON)
		res, err := app.Test(req)
		if err != nil {
			t.Fatalf("failed to initiate app test: %s", err)
		}
		assert.Equal(t, fiber.StatusBadGateway, res.StatusCode)
	})
}

func TestEndMeeting_ErrorPathCausedByBBB(t *testing.T) {
	// prepare fake server to mimic BBB Server
	var fakeEndMeet = func(t *testing.T) *httptest.Server {
		return httptest.NewServer(http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {
			resp := struct {
				Message string `json:"message"`
			}{Message: "should error"}
			js, err := json.Marshal(&resp)
			if err != nil {
				t.Fatalf("failed to marshal object to xml: %s", err)
			}

			rw.WriteHeader(fiber.StatusOK)
			rw.Header().Set("Content-Type", fiber.MIMEApplicationJSON)
			_, err = rw.Write(js)
			require.NoError(t, err)
		}))
	}

	conf, err := config.NewConfig(bytes.NewBufferString(sampleConfigFile[0]))
	require.NoError(t, err)
	require.NoError(t, conf.Sanitization())
	conf.BBB.Host = fakeEndMeet(t).URL
	require.NoError(t, conf.BBB.Sanitization())

	app := fiber.New()
	app.Post("/end", EndMeeting(conf, fakeEndMeet(t).Client()))

	t.Run("Failed if BBB API send back response that has any other than XML content type", func(t *testing.T) {
		buf := bytes.NewBufferString(sampleEndRequest[2])
		req := httptest.NewRequest(fiber.MethodPost, "/end", buf)
		req.Header.Set("Content-Type", fiber.MIMEApplicationJSON)
		res, err := app.Test(req)
		if err != nil {
			t.Fatalf("failed to initiate app test: %s", err)
		}
		assert.Equal(t, fiber.StatusInternalServerError, res.StatusCode)
	})
}

func TestEndMeeting_ErrorPathCausedByBBBGiveNonSuccessResponse(t *testing.T) {
	conf, err := config.NewConfig(bytes.NewBufferString(sampleConfigFile[0]))
	require.NoError(t, err)
	require.NoError(t, conf.Sanitization())
	conf.BBB.Host = fakeEndMeetServer(t).URL
	require.NoError(t, conf.BBB.Sanitization())

	app := fiber.New()
	app.Post("/end", EndMeeting(conf, fakeEndMeetServer(t).Client()))

	t.Run("Failed if BBB API send back non SUCCESS response", func(t *testing.T) {
		buf := bytes.NewBufferString(sampleEndRequest[2])
		req := httptest.NewRequest(fiber.MethodPost, "/end", buf)
		req.Header.Set("Content-Type", fiber.MIMEApplicationJSON)
		res, err := app.Test(req)
		if err != nil {
			t.Fatalf("failed to initiate app test: %s", err)
		}
		assert.Equal(t, fiber.StatusBadGateway, res.StatusCode)
		jsRes := struct {
			Message string `json:"message"`
		}{}
		if err := json.NewDecoder(res.Body).Decode(&jsRes); err != nil {
			t.Fatalf("failed to decode json in test : %v", err)
		}
		assert.Contains(t, jsRes.Message, sampleEndMeetStdResponse[0].MsgKey)
	})
}

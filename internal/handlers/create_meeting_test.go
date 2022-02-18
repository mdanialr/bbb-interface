package handlers

import (
	"bytes"
	"encoding/json"
	"encoding/xml"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gofiber/fiber/v2"
	"github.com/kurvaid/bbb-interface/internal/api"
	"github.com/kurvaid/bbb-interface/internal/config"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

var sampleConfigFile = []string{
	`
env: prod
port: 7575
log: ./log
BBB:
  host:
  secret: secret
`,
}

var sampleRequestBody = []string{
	`
{
	"name": "test-meeting"
}
`,
}

// prepare fake server to mimic BBB Server
var fakeServerHelper = func(t *testing.T) *httptest.Server {
	return httptest.NewServer(http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {
		resp := api.CreateMeetingResponse{
			MeetingId:     "fake-id",
			ModeratorPass: "password",
			AttendeePass:  "secret",
			CreateTime:    "121212",
		}
		xm, err := xml.Marshal(&resp)
		if err != nil {
			t.Fatalf("failed to marshal object to xml: %s", err)
		}

		rw.WriteHeader(fiber.StatusCreated)
		rw.Header().Set("Content-Type", fiber.MIMEApplicationXML)
		_, err = rw.Write(xm)
		require.NoError(t, err)
	}))
}

func TestCreateMeeting(t *testing.T) {
	conf, err := config.NewConfig(bytes.NewBufferString(sampleConfigFile[0]))
	require.NoError(t, err)
	require.NoError(t, conf.Sanitization())
	conf.BBB.Host = fakeServerHelper(t).URL
	require.NoError(t, conf.BBB.Sanitization())

	app := fiber.New()
	app.Post("/meeting", CreateMeeting(conf, fakeServerHelper(t).Client()))

	t.Run("Success using minimum (required) json request", func(t *testing.T) {
		buf := bytes.NewBufferString(sampleRequestBody[0])
		req := httptest.NewRequest(fiber.MethodPost, "/meeting", buf)
		req.Header.Set("Content-Type", fiber.MIMEApplicationJSON)
		res, err := app.Test(req)
		if err != nil {
			t.Fatalf("failed to initiate app test: %s", err)
		}
		assert.Equal(t, fiber.StatusCreated, res.StatusCode)
		assert.Equal(t, fiber.MIMEApplicationJSON, res.Header.Get("Content-Type"))
		var rXML api.CreateMeetingResponse
		if err := json.NewDecoder(res.Body).Decode(&rXML); err != nil {
			t.Fatalf("failed to decode xml in test : %v", err)
		}
		assert.Equal(t, "fake-id", rXML.MeetingId)
		assert.Equal(t, "password", rXML.ModeratorPass)
		assert.Equal(t, "secret", rXML.AttendeePass)
		assert.Equal(t, "121212", rXML.CreateTime)
	})
}

func TestCreateMeeting_FailedUsingFakeBBBHost(t *testing.T) {
	conf, err := config.NewConfig(bytes.NewBufferString(sampleConfigFile[0]))
	require.NoError(t, err)
	require.NoError(t, conf.Sanitization())

	app := fiber.New()
	app.Post("/meeting", CreateMeeting(conf, fakeServerHelper(t).Client()))

	t.Run("Failed when sending using fake host for the BBB server", func(t *testing.T) {
		buf := bytes.NewBufferString(sampleRequestBody[0])
		req := httptest.NewRequest(fiber.MethodPost, "/meeting", buf)
		req.Header.Set("Content-Type", fiber.MIMEApplicationJSON)
		res, err := app.Test(req)
		if err != nil {
			t.Fatalf("failed to initiate app test: %s", err)
		}
		assert.Equal(t, fiber.StatusBadGateway, res.StatusCode)
	})
}

func TestCreateMeeting_FailedUsingEmptyJSONRequest(t *testing.T) {
	conf, err := config.NewConfig(bytes.NewBufferString(sampleConfigFile[0]))
	require.NoError(t, err)
	require.NoError(t, conf.Sanitization())
	conf.BBB.Host = fakeServerHelper(t).URL
	require.NoError(t, conf.BBB.Sanitization())

	app := fiber.New()
	app.Post("/meeting", CreateMeeting(conf, fakeServerHelper(t).Client()))

	t.Run("Failed when sending using fake host for the BBB server", func(t *testing.T) {
		buf := bytes.NewBufferString(`empty`)
		req := httptest.NewRequest(fiber.MethodPost, "/meeting", buf)
		req.Header.Set("Content-Type", fiber.MIMEApplicationJSON)
		res, err := app.Test(req)
		if err != nil {
			t.Fatalf("failed to initiate app test: %s", err)
		}
		assert.Equal(t, fiber.StatusBadRequest, res.StatusCode)
	})

	t.Run("Failed if `name` field not provided", func(t *testing.T) {
		buf := bytes.NewBufferString(`{"key": "value"}`)
		req := httptest.NewRequest(fiber.MethodPost, "/meeting", buf)
		req.Header.Set("Content-Type", fiber.MIMEApplicationJSON)
		res, err := app.Test(req)
		if err != nil {
			t.Fatalf("failed to initiate app test: %s", err)
		}
		assert.Equal(t, fiber.StatusInternalServerError, res.StatusCode)
	})
}

func TestCreateMeeting_FailedCausedByBBB(t *testing.T) {
	// prepare fake server to mimic BBB Server
	var fakeFailedServer = func(t *testing.T) *httptest.Server {
		return httptest.NewServer(http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {
			resp := struct {
				Name string `json:"name"`
			}{}
			js, err := json.Marshal(&resp)
			if err != nil {
				t.Fatalf("failed to marshal object to xml: %s", err)
			}

			rw.WriteHeader(fiber.StatusCreated)
			rw.Header().Set("Content-Type", fiber.MIMEApplicationJSON)
			_, err = rw.Write(js)
			require.NoError(t, err)
		}))
	}

	conf, err := config.NewConfig(bytes.NewBufferString(sampleConfigFile[0]))
	require.NoError(t, err)
	require.NoError(t, conf.Sanitization())
	conf.BBB.Host = fakeFailedServer(t).URL
	require.NoError(t, conf.BBB.Sanitization())

	app := fiber.New()
	app.Post("/meeting", CreateMeeting(conf, fakeFailedServer(t).Client()))

	t.Run("Failed if BBB server send response back using content type other than xml", func(t *testing.T) {
		buf := bytes.NewBufferString(sampleRequestBody[0])
		req := httptest.NewRequest(fiber.MethodPost, "/meeting", buf)
		req.Header.Set("Content-Type", fiber.MIMEApplicationJSON)
		res, err := app.Test(req)
		if err != nil {
			t.Fatalf("failed to initiate app test: %s", err)
		}
		assert.Equal(t, fiber.StatusBadGateway, res.StatusCode)
	})
}

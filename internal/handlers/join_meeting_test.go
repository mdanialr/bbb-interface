package handlers

import (
	"bytes"
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

var sampleJoinRequestBody = []string{
	`
{
	"name": "NzK",
	"meeting_id": "meet01",
	"password": "secret",
	"create_time": "121212"
}
`,
	`
{
	"name": "NzK",
	"meeting_id": "meet01",
	"password": "",
	"create_time": ""
}
`,
}

// prepare fake server to mimic BBB Server
var fakeJoinServerHelper = func(t *testing.T) *httptest.Server {
	return httptest.NewServer(http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {
		resp := api.JoinMeetingResponse{SessionToken: "secureSessionToken"}
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

func TestJoinMeeting(t *testing.T) {
	conf, err := config.NewConfig(bytes.NewBufferString(sampleConfigFile[0]))
	require.NoError(t, err)
	require.NoError(t, conf.Sanitization())
	conf.BBB.Host = fakeJoinServerHelper(t).URL
	require.NoError(t, conf.BBB.Sanitization())

	app := fiber.New()
	app.Post("/meeting", JoinMeeting(conf))

	t.Run("Success using minimum (required) json request", func(t *testing.T) {
		buf := bytes.NewBufferString(sampleJoinRequestBody[0])
		req := httptest.NewRequest(fiber.MethodPost, "/meeting", buf)
		req.Header.Set("Content-Type", fiber.MIMEApplicationJSON)
		res, err := app.Test(req)
		if err != nil {
			t.Fatalf("failed to initiate app test: %s", err)
		}
		assert.Equal(t, fiber.StatusOK, res.StatusCode)
	})

	t.Run("Failed when sending empty value on required fields", func(t *testing.T) {
		buf := bytes.NewBufferString(sampleJoinRequestBody[1])
		req := httptest.NewRequest(fiber.MethodPost, "/meeting", buf)
		req.Header.Set("Content-Type", fiber.MIMEApplicationJSON)
		res, err := app.Test(req)
		if err != nil {
			t.Fatalf("failed to initiate app test: %s", err)
		}
		assert.Equal(t, fiber.StatusInternalServerError, res.StatusCode)
	})
}

func TestJoinMeeting_VariousFailedCases(t *testing.T) {
	conf, err := config.NewConfig(bytes.NewBufferString(sampleConfigFile[0]))
	require.NoError(t, err)
	require.NoError(t, conf.Sanitization())

	app := fiber.New()
	app.Post("/meeting", JoinMeeting(conf))

	t.Run("Failed when sending wrong content type that should be json", func(t *testing.T) {
		buf := bytes.NewBufferString(sampleRequestBody[0])
		req := httptest.NewRequest(fiber.MethodPost, "/meeting", buf)
		req.Header.Set("Content-Type", fiber.MIMEApplicationXML)
		res, err := app.Test(req)
		if err != nil {
			t.Fatalf("failed to initiate app test: %s", err)
		}
		assert.Equal(t, fiber.StatusBadRequest, res.StatusCode)
	})
}

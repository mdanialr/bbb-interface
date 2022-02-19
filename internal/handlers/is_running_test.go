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

var sampleIsRunRequest = `
{
	"meeting_id": "meet01"
}
`

type isRunStdResponse struct {
	Status string `xml:"running"`
	api.StdResponse
}

var sampleIsRunStdResponse = []isRunStdResponse{
	{
		Status:      "true",
		StdResponse: api.StdResponse{CodeString: "SUCCESS"},
	},
	{
		Status:      "false",
		StdResponse: api.StdResponse{CodeString: "SUCCESS"},
	},
	{
		StdResponse: api.StdResponse{CodeString: "FAILED"},
	},
}

// prepare fake server to mimic BBB Server
var fakeIsRunServer = func(t *testing.T) *httptest.Server {
	return httptest.NewServer(http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {
		xm, err := xml.Marshal(&sampleIsRunStdResponse[2])
		if err != nil {
			t.Fatalf("failed to marshal object to xml: %s", err)
		}

		rw.WriteHeader(fiber.StatusOK)
		rw.Header().Set("Content-Type", fiber.MIMEApplicationXML)
		_, err = rw.Write(xm)
		require.NoError(t, err)
	}))
}

func TestIsRunning_VariousErrorPaths(t *testing.T) {
	conf, err := config.NewConfig(bytes.NewBufferString(sampleConfigFile[0]))
	require.NoError(t, err)
	require.NoError(t, conf.Sanitization())

	app := fiber.New()
	app.Post("/is_run", IsRunning(conf, fakeIsRunServer(t).Client()))

	t.Run("Failed if `meeting_id` field is not provided", func(t *testing.T) {
		buf := bytes.NewBufferString(`{"key": "value"}`)
		req := httptest.NewRequest(fiber.MethodPost, "/is_run", buf)
		req.Header.Set("Content-Type", fiber.MIMEApplicationJSON)
		res, err := app.Test(req)
		if err != nil {
			t.Fatalf("failed to initiate app test: %s", err)
		}
		assert.Equal(t, fiber.StatusBadRequest, res.StatusCode)
	})

	t.Run("Failed if request's content type is invalid", func(t *testing.T) {
		buf := bytes.NewBufferString(sampleIsRunRequest)
		req := httptest.NewRequest(fiber.MethodPost, "/is_run", buf)
		req.Header.Set("Content-Type", "should error")
		res, err := app.Test(req)
		if err != nil {
			t.Fatalf("failed to initiate app test: %s", err)
		}
		assert.Equal(t, fiber.StatusBadRequest, res.StatusCode)
	})

	t.Run("Failed if using invalid/empty BBB Server host, even if already has valid requirements", func(t *testing.T) {
		buf := bytes.NewBufferString(sampleIsRunRequest)
		req := httptest.NewRequest(fiber.MethodPost, "/is_run", buf)
		req.Header.Set("Content-Type", fiber.MIMEApplicationJSON)
		res, err := app.Test(req)
		if err != nil {
			t.Fatalf("failed to initiate app test: %s", err)
		}
		assert.Equal(t, fiber.StatusBadGateway, res.StatusCode)
	})
}

func TestIsRunning_ErrorPathCausedByBBB(t *testing.T) {
	// prepare fake server to mimic BBB Server
	var fakeIsRun = func(t *testing.T) *httptest.Server {
		return httptest.NewServer(http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {
			resp := struct {
				Message string `json:"message"`
			}{Message: "should error"}
			js, err := json.Marshal(&resp)
			if err != nil {
				t.Fatalf("failed to marshal object to json: %s", err)
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
	conf.BBB.Host = fakeIsRun(t).URL
	require.NoError(t, conf.BBB.Sanitization())

	app := fiber.New()
	app.Post("/is_run", IsRunning(conf, fakeIsRun(t).Client()))

	t.Run("Failed if BBB API send back response that has any other than XML content type", func(t *testing.T) {
		buf := bytes.NewBufferString(sampleIsRunRequest)
		req := httptest.NewRequest(fiber.MethodPost, "/is_run", buf)
		req.Header.Set("Content-Type", fiber.MIMEApplicationJSON)
		res, err := app.Test(req)
		if err != nil {
			t.Fatalf("failed to initiate app test: %s", err)
		}
		assert.Equal(t, fiber.StatusInternalServerError, res.StatusCode)
	})
}

func TestIsRunning_ErrorPathCausedByBBBGiveNonSuccessResponse(t *testing.T) {
	conf, err := config.NewConfig(bytes.NewBufferString(sampleConfigFile[0]))
	require.NoError(t, err)
	require.NoError(t, conf.Sanitization())
	conf.BBB.Host = fakeIsRunServer(t).URL
	require.NoError(t, conf.BBB.Sanitization())

	app := fiber.New()
	app.Post("/is_run", IsRunning(conf, fakeIsRunServer(t).Client()))

	t.Run("Failed if BBB API send back non SUCCESS response", func(t *testing.T) {
		buf := bytes.NewBufferString(sampleIsRunRequest)
		req := httptest.NewRequest(fiber.MethodPost, "/is_run", buf)
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
		assert.Contains(t, jsRes.Message, "BBB API")
	})
}

func TestIsRunning(t *testing.T) {
	var fakeSuccessIsRunServer = func(t *testing.T) *httptest.Server {
		return httptest.NewServer(http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {
			expect := "/bigbluebutton/api/isMeetingRunning?meetingID=meet01&checksum=38636bd23e45c9cc571054e9ab8d9c80170327b6"
			assert.Equal(t, expect, req.URL.String())

			xm, err := xml.Marshal(&sampleIsRunStdResponse[0])
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
	conf.BBB.Host = fakeSuccessIsRunServer(t).URL
	require.NoError(t, conf.BBB.Sanitization())

	app := fiber.New()
	app.Post("/is_run", IsRunning(conf, fakeSuccessIsRunServer(t).Client()))

	t.Run("Success if all required fields are provided", func(t *testing.T) {
		buf := bytes.NewBufferString(sampleIsRunRequest)
		req := httptest.NewRequest(fiber.MethodPost, "/is_run", buf)
		req.Header.Set("Content-Type", fiber.MIMEApplicationJSON)
		res, err := app.Test(req)
		if err != nil {
			t.Fatalf("failed to initiate app test: %s", err)
		}
		assert.Equal(t, fiber.StatusOK, res.StatusCode)
		jsRes := struct {
			Status bool `json:"status"`
		}{}
		if err := json.NewDecoder(res.Body).Decode(&jsRes); err != nil {
			t.Fatalf("failed to decode json in test : %v", err)
		}
		assert.Equal(t, true, jsRes.Status)
	})
}

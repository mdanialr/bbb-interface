package handlers

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gofiber/fiber/v2"
	"github.com/kurvaid/bbb-interface/internal/config"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestCallbackOnDestroy(t *testing.T) {
	// prepare fake server just to make this test pass.
	var fakeCallbackServerHelper = func() *httptest.Server {
		return httptest.NewServer(http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {
		}))
	}

	conf, err := config.NewConfig(bytes.NewBufferString(sampleConfigFile[0]))
	conf.CallbackOnDestroy = fakeCallbackServerHelper().URL
	require.NoError(t, err)
	require.NoError(t, conf.Sanitization())

	app := fiber.New()
	app.Get("/callback/destroy", CallbackOnDestroy(conf, fakeCallbackServerHelper().Client()))

	t.Run("Every GET request to this endpoint should pass", func(t *testing.T) {
		meetID := "meet01"
		uri := fmt.Sprintf("/callback/destroy?meetingID=%s", meetID)
		req := httptest.NewRequest(fiber.MethodGet, uri, nil)
		res, err := app.Test(req)
		require.NoError(t, err)
		assert.Equal(t, fiber.StatusOK, res.StatusCode)
	})
}

// prepare fake server to mimic lms app.
var fakeCallbackServerHelper = func(expectPayload string, t *testing.T) *httptest.Server {
	return httptest.NewServer(http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {
		// make sure its POST request
		assert.Equal(t, fiber.MethodPost, req.Method, "requests method must be `POST`")

		payload, err := io.ReadAll(req.Body)
		require.NoError(t, err, fmt.Sprintf("failed to read request body: %s", err))

		var modelReq DestroyCallbackModel
		err = json.Unmarshal(payload, &modelReq)
		require.NoError(t, err, fmt.Sprintf("failed to bind request body to model: %s", err))

		assert.Equal(t, expectPayload, modelReq.MeetingId, fmt.Sprintf("meeting id should be %s", expectPayload))

		js, err := json.Marshal(&fiber.Map{
			"code":    200,
			"message": "Success !",
			"data":    nil,
			"errors":  nil,
		})
		require.NoError(t, err, fmt.Sprintf("failed to marshal object to json: %s", err))

		rw.WriteHeader(fiber.StatusOK)
		rw.Header().Set("Content-Type", fiber.MIMEApplicationJSON)
		_, err = rw.Write(js)
		require.NoError(t, err, fmt.Sprintf("failed to write json bytes to response writer %s", err))
	}))
}

func TestCallbackOnDestroy_UsingFakeServerToReceivePOSTCallback(t *testing.T) {
	testCases := []struct {
		name   string
		sample string
		expect string
	}{
		{
			name:   "Should sent POST payload exactly as sent by BBB server",
			sample: "meet03",
			expect: "meet03",
		},
		{
			name:   "Should sent return 200 status after sent POST request to lms",
			sample: "meet01",
			expect: "meet01",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			conf, err := config.NewConfig(bytes.NewBufferString(sampleConfigFile[0]))
			conf.CallbackOnDestroy = fakeCallbackServerHelper(tc.expect, t).URL
			require.NoError(t, err)
			require.NoError(t, conf.Sanitization())

			app := fiber.New()
			app.Get("/callback/destroy", CallbackOnDestroy(conf, fakeCallbackServerHelper(tc.expect, t).Client()))

			uri := fmt.Sprintf("/callback/destroy?meetingID=%s", tc.sample)
			req := httptest.NewRequest(fiber.MethodGet, uri, nil)
			res, err := app.Test(req)
			require.NoError(t, err)
			assert.Equal(t, fiber.StatusOK, res.StatusCode)
		})
	}
}

func TestCallbackOnDestroy_ErrorPathToImproveCoverage(t *testing.T) {
	// prepare fake server just to make this test pass.
	var fakeCallbackServerHelper = func() *httptest.Server {
		return httptest.NewServer(http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {
		}))
	}

	conf, err := config.NewConfig(bytes.NewBufferString(sampleConfigFile[0]))
	conf.CallbackOnDestroy = "http://localhost"
	require.NoError(t, err)
	require.NoError(t, conf.Sanitization())

	app := fiber.New()
	app.Get("/callback/destroy", CallbackOnDestroy(conf, fakeCallbackServerHelper().Client()))

	t.Run("Using fake server url should error and return 500 status code", func(t *testing.T) {
		req := httptest.NewRequest(fiber.MethodGet, "/callback/destroy", nil)
		res, err := app.Test(req)
		require.NoError(t, err)
		assert.Equal(t, fiber.StatusInternalServerError, res.StatusCode)
	})
}

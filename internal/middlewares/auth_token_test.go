package middlewares

import (
	"bytes"
	"net/http/httptest"
	"testing"

	"github.com/gofiber/fiber/v2"
	"github.com/kurvaid/bbb-interface/internal/config"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

var sampleConfigFile = []string{
	`
env: prod
port: 7575
log: ./log
token: superSecret
`, `
env: prod
port: 7575
log: ./log
token:
`,
}

func TestAuthMiddleware(t *testing.T) {
	conf, err := config.NewConfig(bytes.NewBufferString(sampleConfigFile[0]))
	require.NoError(t, err)
	require.NoError(t, conf.Sanitization())

	app := fiber.New()
	app.Post("/auth",
		Auth(conf),
		func(c *fiber.Ctx) error {
			return c.SendStatus(fiber.StatusOK)
		},
	)

	t.Run("Failed if token in header is empty", func(t *testing.T) {
		req := httptest.NewRequest(fiber.MethodPost, "/auth", nil)
		res, err := app.Test(req)
		require.NoError(t, err, "failed to initiate app test: ", err)
		assert.Equal(t, fiber.StatusBadRequest, res.StatusCode)
	})

	t.Run("Failed if token in header doesn't match with config's token", func(t *testing.T) {
		req := httptest.NewRequest(fiber.MethodPost, "/auth", nil)
		req.Header.Set("Authorization", "secret")
		res, err := app.Test(req)
		require.NoError(t, err, "failed to initiate app test: ", err)
		assert.Equal(t, fiber.StatusBadRequest, res.StatusCode)
	})

	t.Run("Pass if token in header does match with config's token", func(t *testing.T) {
		req := httptest.NewRequest(fiber.MethodPost, "/auth", nil)
		req.Header.Set("Authorization", "superSecret")
		res, err := app.Test(req)
		require.NoError(t, err, "failed to initiate app test: ", err)
		assert.Equal(t, fiber.StatusOK, res.StatusCode)
	})

	conf, err = config.NewConfig(bytes.NewBufferString(sampleConfigFile[1]))
	require.NoError(t, err)
	require.NoError(t, conf.Sanitization())

	app = fiber.New()
	app.Post("/auth",
		Auth(conf),
		func(c *fiber.Ctx) error {
			return c.SendStatus(fiber.StatusOK)
		},
	)

	t.Run("Pass if token in config and header request is empty", func(t *testing.T) {
		req := httptest.NewRequest(fiber.MethodPost, "/auth", nil)
		req.Header.Set("Authorization", "")
		res, err := app.Test(req)
		require.NoError(t, err, "failed to initiate app test: ", err)
		assert.Equal(t, fiber.StatusOK, res.StatusCode)
	})
}

package routes

import (
	"testing"

	"github.com/gofiber/fiber/v2"
	"github.com/kurvaid/bbb-interface/internal/config"
)

func TestSetupRoutes(t *testing.T) {
	t.Run("1# Success test even there is no value supplied because there is no required value", func(t *testing.T) {
		conf := config.Model{}
		app := fiber.New()

		SetupRoutes(app, &conf)
	})

	t.Run("2# Success test with only one or more supplied value", func(t *testing.T) {
		conf := config.Model{EnvIsProd: true}
		app := fiber.New()

		SetupRoutes(app, &conf)
	})
}

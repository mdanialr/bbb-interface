package middlewares

import (
	"github.com/gofiber/fiber/v2"
	"github.com/kurvaid/bbb-interface/internal/config"
)

// Auth super simple middleware to check whether wanted Authorization token exist
// and match with token in config.
func Auth(conf *config.Model) func(ctx *fiber.Ctx) error {
	return func(c *fiber.Ctx) error {
		token := c.GetReqHeaders()["Authorization"]

		if conf.Token != token {
			c.Status(fiber.StatusBadRequest)
			return c.JSON(fiber.Map{
				"message": "token doesn't match",
			})
		}

		return c.Next()
	}
}

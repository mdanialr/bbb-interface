package handlers

import (
	"fmt"
	"github.com/gofiber/fiber/v2"
	"github.com/kurvaid/bbb-interface/internal/api"
	"github.com/kurvaid/bbb-interface/internal/config"
	"github.com/kurvaid/bbb-interface/internal/service"
)

// JoinMeeting handler that receive json request and proxy it to BBB API for joining meeting
// after convert to URL then send back response from API to the requester.
func JoinMeeting(conf *config.Model) func(c *fiber.Ctx) error {
	return func(c *fiber.Ctx) error {
		// bind incoming json request to predefined object.
		var jMeet api.JoinMeeting
		if err := c.BodyParser(&jMeet); err != nil {
			c.Status(fiber.StatusBadRequest)
			return c.JSON(fiber.Map{
				"message": fmt.Errorf("failed to bind request to join meeting object: %s", err),
			})
		}

		url, err := jMeet.ParseJoinMeeting()
		if err != nil {
			c.Status(fiber.StatusInternalServerError)
			return c.JSON(fiber.Map{
				"message": fmt.Sprintf("failed to parse join meeting url: %s", err),
			})
		}
		// prepare url and calculate their checksum.
		out := service.SHA1HashUrl(conf.BBB.Secret, url)
		url = fmt.Sprintf("%s%s%s", conf.BBB.Host, api.EndPoint, url)

		return c.JSON(fiber.Map{
			"url": fmt.Sprintf("%s&checksum=%s", url, out),
		})
	}
}

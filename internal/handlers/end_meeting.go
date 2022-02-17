package handlers

import (
	"encoding/xml"
	"fmt"
	"net/http"

	"github.com/gofiber/fiber/v2"
	"github.com/kurvaid/bbb-interface/internal/api"
	"github.com/kurvaid/bbb-interface/internal/client"
	"github.com/kurvaid/bbb-interface/internal/config"
	"github.com/kurvaid/bbb-interface/internal/service"
)

// EndMeeting handler that receive json request to end a meeting from client and transform it to xml
// request that match BBB API requirement and transform xml response to json before send it back to
// the client.
func EndMeeting(conf *config.Model, hCl *http.Client) func(*fiber.Ctx) error {
	return func(c *fiber.Ctx) error {
		// bind incoming json request to predefined object.
		var eMeet api.EndMeeting
		if err := c.BodyParser(&eMeet); err != nil {
			c.Status(fiber.StatusBadRequest)
			return c.JSON(fiber.Map{
				"message": fmt.Errorf("failed to parse request to end meeting object: %s", err),
			})
		}

		uri, err := eMeet.ParseEndMeeting()
		if err != nil {
			c.Status(fiber.StatusBadRequest)
			return c.JSON(fiber.Map{
				"message": fmt.Sprintf("failed to parse end meeting url: %s", err),
			})
		}

		// prepare url and calculate their checksum.
		out := service.SHA1HashUrl(conf.BBB.Secret, uri)
		uri = fmt.Sprintf("%s%s%s", conf.BBB.Host, api.EndPoint, uri)

		endMeetApi := client.Instance{Cl: hCl, Url: uri, Checksum: out}

		resp, err := endMeetApi.DispatchGET()
		if err != nil {
			c.Status(fiber.StatusBadGateway)
			return c.JSON(fiber.Map{
				"message": fmt.Sprintf("failed sending end meeting request to BBB API: %s", err),
			})
		}

		var res api.StdResponse
		if err := xml.Unmarshal(resp, &res); err != nil {
			c.Status(fiber.StatusInternalServerError)
			return c.JSON(fiber.Map{
				"message": fmt.Sprintf("failed parsing BBB API response to std response object: %s", err),
			})
		}

		// check if BBB API call success
		if res.CodeString != "SUCCESS" {
			c.Status(fiber.StatusBadGateway)
			return c.JSON(fiber.Map{
				"message": fmt.Sprintf("receiving error from BBB API: [%s] %s", res.MsgKey, res.MsgDetail),
			})
		}

		c.Status(fiber.StatusOK)
		return c.JSON(fiber.Map{
			"message": fmt.Sprintf("meeting %s successfully deleted", eMeet.MeetingId),
		})

	}
}

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

// CreateMeeting handler that receive json request and proxy it to BBB API after convert to URL
// then send back response from BBB API to the requester.
func CreateMeeting(conf *config.Model, httpClient *http.Client) func(c *fiber.Ctx) error {
	return func(c *fiber.Ctx) error {
		// bind incoming json request to predefined object.
		var cMeet api.CreateMeeting
		if err := c.BodyParser(&cMeet); err != nil {
			c.Status(fiber.StatusBadRequest)
			return c.JSON(fiber.Map{
				"message": fmt.Errorf("failed to bind request to create meeting object: %s", err),
			})
		}

		randNum := service.RandomString{Length: int(conf.RandomLen)}
		url, err := cMeet.ParseCreateMeeting(&randNum)
		if err != nil {
			c.Status(fiber.StatusInternalServerError)
			return c.JSON(fiber.Map{
				"message": fmt.Sprintf("failed to parse create meeting url: %s", err),
			})
		}
		// prepare url and calculate their checksum.
		out := service.SHA1HashUrl(conf.BBB.Secret, url)
		url = fmt.Sprintf("%s%s%s", conf.BBB.Host, api.EndPoint, url)

		createMeetApi := client.Create{Cl: httpClient, Url: url, CheckSum: out}

		resp, err := createMeetApi.CreateMeeting()
		if err != nil {
			c.Status(fiber.StatusBadGateway)
			return c.JSON(fiber.Map{
				"message": fmt.Sprintf("failed sending create meeting request to BBB API: %s", err),
			})
		}

		var jsonResp api.CreateMeetingResponse
		if err := xml.Unmarshal(resp, &jsonResp); err != nil {
			c.Status(fiber.StatusBadGateway)
			return c.JSON(fiber.Map{
				"message": fmt.Sprintf("failed binding BBB API response to json response: %s", err),
			})
		}

		c.Status(fiber.StatusCreated)
		return c.JSON(jsonResp)
	}
}

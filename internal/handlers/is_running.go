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

// IsRunning handler that receive json request to check whether a meeting is running or not from
// client and transform it to xml request that match BBB API requirement and transform xml
// response to json before send it back to the client.
func IsRunning(conf *config.Model, hCl *http.Client) func(*fiber.Ctx) error {
	return func(c *fiber.Ctx) error {
		// bind incoming json request to predefined object.
		var isRun api.IsRunning
		if err := c.BodyParser(&isRun); err != nil {
			c.Status(fiber.StatusBadRequest)
			return c.JSON(fiber.Map{
				"message": fmt.Errorf("failed to parse request to is running object: %s", err),
			})
		}

		uri, err := isRun.ParseIsRunning()
		if err != nil {
			c.Status(fiber.StatusBadRequest)
			return c.JSON(fiber.Map{
				"message": fmt.Sprintf("failed to parse is meeting running url: %s", err),
			})
		}

		// prepare url and calculate their checksum.
		out := service.SHA1HashUrl(conf.BBB.Secret, uri)
		uri = fmt.Sprintf("%s%s%s", conf.BBB.Host, api.EndPoint, uri)

		isRunApi := client.Instance{Cl: hCl, Url: uri, Checksum: out}

		resp, err := isRunApi.DispatchGET()
		if err != nil {
			c.Status(fiber.StatusBadGateway)
			return c.JSON(fiber.Map{
				"message": fmt.Sprintf("failed sending is meeting running request to BBB API: %s", err),
			})
		}

		var res struct {
			api.StdResponse
			Status bool `xml:"running" json:"status"`
		}
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
				"message": "something was wrong with BBB API",
			})
		}

		c.Status(fiber.StatusOK)
		return c.JSON(res)
	}
}

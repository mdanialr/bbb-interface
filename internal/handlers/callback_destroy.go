package handlers

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/gofiber/fiber/v2"
	"github.com/kurvaid/bbb-interface/internal/config"
)

// DestroyCallbackModel model that provided by lms app to notify that a meeting
// has destroyed or ended.
type DestroyCallbackModel struct {
	MeetingId string `json:"meeting_id"` // Meeting id that determine which meeting was destroyed.
}

// CallbackOnDestroy handler that will receive GET request from BBB server when a meeting was destroyed
// or ended, then sent POST request to designated lms endpoint complete with the body request that
// would determine which meeting was destroyed using meeting_id sent by BBB server's GET request.
func CallbackOnDestroy(conf *config.Model, htC *http.Client) func(*fiber.Ctx) error {
	return func(c *fiber.Ctx) error {
		// proses incoming URL from BBB server
		meetId := c.Query("meetingID")
		payload := &DestroyCallbackModel{MeetingId: meetId}

		jsonPayload, err := json.Marshal(&payload)
		if err != nil {
			return fmt.Errorf("failed to marshal model to json: %s", err)
		}

		res, err := htC.Post(conf.CallbackOnDestroy, fiber.MIMEApplicationJSON, bytes.NewReader(jsonPayload))
		if err != nil {
			return fmt.Errorf("failed to send request to lms endpoint: %s", err)
		}
		defer res.Body.Close()

		return c.SendStatus(fiber.StatusOK)
	}
}

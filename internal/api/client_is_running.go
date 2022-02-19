package api

import "fmt"

// IsRunning format that needed to check whether a meeting is running.
type IsRunning struct {
	MeetingId string `json:"meeting_id"` // The meeting ID that identifies the meeting you are attempting to check. Required.
}

// ParseIsRunning parse the given object that should be sent by client, sanitize it, then transform it to
// format that match BBB API requirement to check whether a meeting is running or not.
func (i *IsRunning) ParseIsRunning() (string, error) {
	if i.MeetingId == "" {
		return "", fmt.Errorf("`meeting_id` field is required")
	}

	str := fmt.Sprintf(
		"/%s?meetingID=%s",
		IsMeetingRun,
		i.MeetingId,
	)

	return str, nil
}

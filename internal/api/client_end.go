package api

import "fmt"

// EndMeeting format that needed to end a meeting.
type EndMeeting struct {
	MeetingId string `json:"meeting_id"` // The meeting ID that identifies the meeting you are attempting to destroy/end/remove. Required.
	Password  string `json:"password"`   // The password of a moderator for this particular meeting. Required.
}

// ParseEndMeeting parse the given object that should be sent by client, sanitize it, then transform it to
// that match BBB API requirement to end a meeting.
func (e *EndMeeting) ParseEndMeeting() (string, error) {
	if e.MeetingId == "" {
		return "", fmt.Errorf("`meeting_id` field is required")
	}

	if e.Password == "" {
		return "", fmt.Errorf("`password` field is required")
	}

	str := fmt.Sprintf(
		"/%s?meetingID=%s&password=%s",
		End,
		e.MeetingId,
		e.Password,
	)

	return str, nil
}

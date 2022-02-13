package api

import (
	"fmt"
	"net/url"

	"github.com/kurvaid/bbb-interface/internal/service"
)

// CreateMeeting format that needed to create meeting. This should be sent as URL.
type CreateMeeting struct {
	Name             string `json:"name"` // A name for the meeting. Required.
	MeetingId        string // A meeting ID that can be used to identify this meeting by the 3rd-party application. Required.
	AttendeePass     string `json:"attendee_pass"`      // Password that would be used by attendee to enter the meeting. Optional.
	ModeratorPass    string `json:"moderator_pass"`     // Password that would be used by moderator to enter the meeting. Optional.
	MaxParticipants  uint8  `json:"max_participant"`    // Set the maximum number of users allowed to join the conference at the same time.
	RedirectAtLogout string `json:"redirect_at_logout"` // The URL that the BigBlueButton client will go to after users click the OK button on the ‘You have been logged out message’.
}

// CreateMeetingResponse holds data from BBB API response after create meeting.
type CreateMeetingResponse struct {
	StatusCode    string `xml:"returncode"`
	MeetingId     string `xml:"meetingID" json:"meeting_id"`
	AttendeePass  string `xml:"attendeePW" json:"attendee_pass"`
	ModeratorPass string `xml:"moderatorPW" json:"moderator_pass"`
	CreateTime    string `xml:"createTime" json:"create_time"`
	CreatedAt     string `xml:"createDate" json:"created_at"`
	Duration      string `xml:"duration" json:"duration"`
}

// ParseCreateMeeting parse given request body binding from json and convert them
// to url string that meet BBB API requirements.
func (cm *CreateMeeting) ParseCreateMeeting(ran service.RandStringInterface) (string, error) {
	if cm.Name == "" {
		return "", fmt.Errorf("`name` field is required")
	}

	if cm.MeetingId == "" {
		cm.MeetingId = ran.RandString()
	}

	if cm.ModeratorPass == "" {
		cm.ModeratorPass = ran.RandString()
	}

	if cm.AttendeePass == "" {
		cm.AttendeePass = ran.RandString()
	}

	str := fmt.Sprintf(
		"/%s?name=%s&meetingID=%s&moderatorPW=%s&attendeePW=%s&logoutURL=%s",
		Create,
		url.QueryEscape(cm.Name),
		cm.MeetingId,
		cm.ModeratorPass,
		cm.AttendeePass,
		cm.RedirectAtLogout,
	)

	if cm.MaxParticipants != 0 {
		str += fmt.Sprintf("&maxParticipants=%d", cm.MaxParticipants)
	}

	return str, nil
}

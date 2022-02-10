package api

import "fmt"

// JoinMeeting format that needed to join a meeting. This should be sent to BBB API as URL.
type JoinMeeting struct {
	Name       string `json:"name"`        // The full name that is to be used to identify this user to other conference attendees. Required.
	MeetingId  string `json:"meeting_id"`  // The meeting ID that identifies the meeting you are attempting to join. Required.
	Password   string `json:"password"`    // The password that this attendee is using. Also, to determine whether this attendee is moderator or not based on the password given. Required.
	CreateTime string `json:"create_time"` // BigBlueButton will ensure it matches the ‘createTime’ for the session. If they differ, BigBlueButton will not proceed with the join request. This prevents a user from reusing their join URL for a subsequent session with the same meetingID. Required.
	UserId     string `json:"user_id"`     // An identifier for this user that will help your application to identify which person this is. This user ID will be returned for this user in the getMeetingInfo API call so that you can check.
	Avatar     string `json:"avatar"`      // The link for the user’s avatar to be displayed.
	IsGuest    bool   `json:"is_guest"`    // To indicate that the user is a guest.
}

// JoinMeetingResponse holds data from BBB API response after join the meeting.
type JoinMeetingResponse struct {
	SessionToken string `xml:"session_token"`
	Url          string `xml:"url"`
}

// ParseJoinMeeting parse given request body binding from json and convert them
// to url string that meet BBB API requirements.
func (j *JoinMeeting) ParseJoinMeeting() (string, error) {
	// /join?meetingID=test01&password=ap&fullName=Chris&createTime=273648
	if j.Name == "" {
		return "", fmt.Errorf("`name` is required")
	}

	if j.MeetingId == "" {
		return "", fmt.Errorf("`meeting_id` is required")
	}

	if j.Password == "" {
		return "", fmt.Errorf("`password` is required")
	}

	if j.CreateTime == "" {
		return "", fmt.Errorf("`created_time` is required")
	}

	str := fmt.Sprintf(
		"/%s?meetingID=%s&password=%s&fullName=%s&createTime=%s",
		Join,
		j.MeetingId,
		j.Password,
		j.Name,
		j.CreateTime,
	)

	return str, nil
}

package api

const (
	Create              = "create"              // Creates a new meeting.
	GetDefaultConfig    = "getDefaultConfigXML" // Gets the default config.xml (these settings configure the BigBlueButton client for each user).
	SetConfig           = "setConfigXML"        // Add a custom config.xml to an existing meeting.
	Join                = "join"                // Join a new user to an existing meeting.
	End                 = "end"                 // Ends meeting.
	IsMeetingRun        = "isMeetingRunning"    // Checks whether if a specified meeting is running.
	GetAllMeetings      = "getMeetings"         // Get the list of Meetings.
	GetMeetingDetail    = "getMeetingInfo"      // Get the details of a Meeting.
	GetAllRecordings    = "getRecordings"       // Get a list of recordings.
	DeleteRecording     = "deleteRecordings"    // Deletes an existing recording.
	UpdateRecordingMeta = "updateRecordings"    // Updates metadata in a recording.
	EndPoint            = "bigbluebutton/api"   // BBB API endpoint.
)

// StdResponse standard response from BBB API either one field or all fields would always be populated when
// receiving response from BBB API after make an API call.
type StdResponse struct {
	CodeString string `xml:"returncode" json:"-"` // May contain only either SUCCESS or FAILED.
	MsgKey     string `xml:"messageKey" json:"-"` // A unique key defined by BBB API to identify which error are thrown.
	MsgDetail  string `xml:"message" json:"-"`    // Detail message about the occurred error.
}

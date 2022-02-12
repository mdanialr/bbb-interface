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

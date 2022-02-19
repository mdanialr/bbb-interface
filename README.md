# BBB Interface
BigBlueButton API wrapper to make the other apps interact with BBB server easier.
## Features
* Create Meeting.
* Join Meeting.
* End Meeting. [*__forcibly end meeting__*]
* Is Meeting Running. [*__check whether a meeting is currently running or not__*]

# How to Use
## Auth/Token
Include the token in your request header.
```
key: Authorization
value: TOKEN
```
Example in GO
```go
import "net/http"

req, _ := http.NewRequest("POST", "http://url.example/endpoint", nil)
req.Header.Set("Authorization", "theTOKEN")
```
Example in Python
```python
import requests

headers = {"Authorization": "theTOKEN"}
requests.post("http://url.example/endpoint", headers=headers)
```

## Error (*if any*)
#### All error response return either *4xx* or *5xx* status code. 

Example Error Response. _in case there are errors._
```json
{
    "message": "the detail of the error"
}
```

## Create Meeting
> `POST` /create

Example Request
```json
{
    "name": "meeting with earth",
    "attendee_pass": "password-for-attendee",
    "moderator_pass": "password-for-moderator",
    "max_participant": 100,
    "redirect_at_logout": "https://maybe-back-to-lms.com/dashboard",
    "welcome_msg": "Hello from earth!!",
    "is_recording": true
}
```
Example Response
```json
{
    "meeting_id": "someRandomString",
    "attendee_pass": "password-for-attendee",
    "moderator_pass": "password-for-moderator",
    "create_time": "1531155809613",
    "created_at": "Mon Jul 09 17:03:29 UTC 2040"
}
```
### Parameters
> Request

`name` `string` `required`: A name for the meeting.

`attendee_pass` `string`: Password that would be used by attendee to enter the meeting. Optional. Would be generated by this service if not provided.

`moderator_pass` `string`: Password that would be used by moderator to enter the meeting. Optional. Would be generated by this service if not provided.

`max_participant` `int`: Set the maximum number of users allowed to join the conference at the same time.

`redirect_at_logout` `string`: The URL that the BigBlueButton client will go to after users click the OK button on the ‘You have been logged out message’.

`welcome_msg` `string`: A welcome message that gets displayed on the chat window when the participant joins.

`is_recording` `boolean`: Enable button to start/pause/stop recording the meeting.

> Response

`meeting_id` `string`: A meeting ID that can be used to identify this meeting by the 3rd-party application.

`create_time` `string`: Berguna saat join meeting. Third-party apps using the API can now pass createTime parameter (which was created in the create call), BigBlueButton will ensure it matches the ‘createTime’ for the session. If they differ, BigBlueButton will not proceed with the join request. This prevents a user from reusing their join URL for a subsequent session with the same meetingID.

`created_at` `string`: You know this ..right?.

## Join Meeting
> `POST` /join

Example Request
```json
{
    "name": "nama Mahasiswa Atau Dosen",
    "meeting_id": "someRandomStringFromCreateCall",
    "password": "passwordThatWillDecideThisUserIsAttendeeOrModerator",
    "create_time": "1531155809613",
    "user_id": "mhs 01",
    "avatar": "https://maybe-back-to-lms.com/assets/avatar/mhs01.png",
    "is_guest": false
}
```

Example Response. It's simply just parse the incoming request then append the calculated checksum. 
```json
{
    "url": "http://bbb-server.test/bigbluebutton/api/join?meetingID=someRandomStringFromCreateCall..."
}
```
### Parameters
> Request

`name` `string` `required`: The full name that is to be used to identify this user to other conference attendees.

`meeting_id` `string` `required`: The meeting ID that identifies the meeting you are attempting to join.

`password` `string` `required`: The password that this attendee is using. Also, to determine whether this attendee is moderator or not based on the password given.

`create_time` `string` `required`: BigBlueButton will ensure it matches the ‘createTime’ for the session. If they differ, BigBlueButton will not proceed with the join request. This prevents a user from reusing their join URL for a subsequent session with the same meetingID.

`user_id` `string`: An identifier for this user that will help your application to identify which person this is. This user ID will be returned for this user in the getMeetingInfo API call so that you can check.

`avatar` `string`: The link for the user’s avatar to be displayed.

`is_guest` `boolean`: To indicate that the user is a guest.

> Response

`url` `string`: The url that need to be open up in browser otherwise would end up got 401 error when joining.

## End Meeting
> `POST` /end

Example Request
```json
{
    "meeting_id": "someRandomStringFromCreateCall",
    "password": "passwordOfTheModeratorOfThisParticularMeeting"
}
```

Example Response.
```json
{
    "message": "meeting someRandomStringFromCreateCall successfully deleted"
}
```
### Parameters
> Request

`meeting_id` `string` `required`: The meeting ID that identifies the meeting you are attempting to forcibly ended.

`password` `string` `required`: The password of the moderator of this meeting.

> Response

`message` `string`: A message to indicate that this api call is success and should also return with http status code `200`.

## Is a Meeting Running
> `POST` /is_run

Example Request
```json
{
    "meeting_id": "someRandomStringFromCreateCall"
}
```

Example Response.
```json
{
    "status": true
}
```
### Parameters
> Request

`meeting_id` `string` `required`: The meeting ID that identifies the meeting you are attempting to check whether its currently running or not.

> Response

`status` `boolean`: The status of the meeting's running state and should also return with http status code `200`.
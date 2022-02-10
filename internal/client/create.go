package client

import (
	"fmt"
	"io"
	"net/http"

	"github.com/kurvaid/bbb-interface/internal/api"
)

// CreateInterface signature that implement this client to create meeting.
type CreateInterface interface {
	CreateMeeting(meeting api.CreateMeeting) ([]byte, error)
}

// Create holds data to send request to BBB API.
type Create struct {
	Cl       *http.Client
	Url      string
	CheckSum string
}

// CreateMeeting take json and transform it to url. Send GET request to BBB API using it.
// Then return response from BBB API.
func (c *Create) CreateMeeting() ([]byte, error) {
	// append checksum at the end of url.
	url := fmt.Sprintf("%s&checksum=%s", c.Url, c.CheckSum)

	// append parsed url with given base url.
	res, err := c.Cl.Get(url)
	if err != nil {
		return nil, fmt.Errorf("failed to send request to BBB API: %s", err)
	}
	defer res.Body.Close()

	body, err := io.ReadAll(res.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response body from BBB API: %s", err)
	}

	return body, nil
}

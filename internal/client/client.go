package client

import (
	"fmt"
	"io"
	"net/http"
)

// Instance holds data to send request to BBB API.
type Instance struct {
	Cl       *http.Client
	Url      string
	Checksum string
}

// DispatchGET take json and transform it to url. Send GET request to BBB API using it.
// Then return response from BBB API.
func (i *Instance) DispatchGET() ([]byte, error) {
	// append checksum at the end of url.
	url := fmt.Sprintf("%s&checksum=%s", i.Url, i.Checksum)

	res, err := i.Cl.Get(url)
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

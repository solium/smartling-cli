// methods that wrap around http communication

package smartling

import (
	"bytes"
	"net/http"
	"io/ioutil"
)

func (c *Client) doPostRequest(apiCall string, data []byte) (response []byte, statusCode int, err error) {

	request, err := http.NewRequest("POST", apiCall, bytes.NewBuffer(data))
	if err != nil {
		return
	}
	request.Header.Set("Content-Type", "application/json")

	resp, err := c.httpClient.Do(request)
	if err != nil {
		return
	}
	defer resp.Body.Close()

	response, err = ioutil.ReadAll(resp.Body)
	if err != nil {
		return
	}

	statusCode = resp.StatusCode
	return
}

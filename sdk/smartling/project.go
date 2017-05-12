package smartling

import (
	"encoding/json"
	"fmt"
	"log"
)


// API endpoints
const (
	projectApiList = "/accounts-api/v2/accounts/%v/projects"
)


func (c *Client) ListProjects(accountId string) (error) {

	header := c.auth.AccessHeader()
	bytes, statusCode, err := c.doGetRequest(c.baseUrl + fmt.Sprintf(projectApiList, accountId), header)

	if err != nil {
		return err
	}

	log.Printf("%v", string(bytes))

	if statusCode != 200 {
		return fmt.Errorf("ListProjects call returned unexpected status code: %v", statusCode)
	}

	// unmarshal transport header
	apiResponse := SmartlingApiResponse{}

	err = json.Unmarshal(bytes, &apiResponse)
	if err != nil {
		return err
	}
	log.Printf("%#v", apiResponse)

	log.Printf("List proijects - received %v status code", statusCode)

	return nil
}

package smartling

import (
	"encoding/json"
	"fmt"
	"log"
)


// API endpoints
const (
	projectApiList    = "/accounts-api/v2/accounts/%v/projects"
	projectApiDetails = "/projects-api/v2/projects/%v"
)


// list project api call response data
type ProjectsApiResponse struct {
	Items []Project
}

type Project struct {
	PprojectId string
	ProjectName string
	AccountUid string
	SourceLocaleId string
	SourceLocaleDescription string
	Archived bool
}

type ProjectDetails struct {
	Project
	TargetLocales []Locale
}

type Locale struct {
	LocaleId string
	Description string
}


func (c *Client) ListProjects(accountId string) (projects []Project, err error) {

	header, err := c.auth.AccessHeader(c)
	if err != nil {
		return
	}

	bytes, statusCode, err := c.doGetRequest(c.baseUrl + fmt.Sprintf(projectApiList, accountId), header)
	if err != nil {
		return
	}

	if statusCode != 200 {
		err = fmt.Errorf("ListProjects call returned unexpected status code: %v", statusCode)
		return
	}

	// unmarshal transport header
	apiResponse := SmartlingApiResponse{}
	err = json.Unmarshal(bytes, &apiResponse)
	if err != nil {
		return
	}

	// unmarshal projects array
	projectsApiResponse := ProjectsApiResponse{}
	err = json.Unmarshal(apiResponse.Response.Data, &projectsApiResponse)
	if err != nil {
		return
	}

	log.Printf("List proijects - received %v status code", statusCode)

	return projectsApiResponse.Items, nil
}

func (c *Client) ProjectDetails(projectId string) (projectDetails ProjectDetails, err error) {

	header, err := c.auth.AccessHeader(c)
	if err != nil {
		return
	}

	bytes, statusCode, err := c.doGetRequest(c.baseUrl + fmt.Sprintf(projectApiDetails, projectId), header)
	if err != nil {
		return
	}

	if statusCode != 200 {
		err = fmt.Errorf("ProjectDetailss call returned unexpected status code: %v", statusCode)
		return
	}

	// unmarshal transport header
	apiResponse := SmartlingApiResponse{}
	err = json.Unmarshal(bytes, &apiResponse)
	if err != nil {
		return
	}

	log.Printf("Project Details - received %v status code", statusCode)

	// unmarshal project details
	err = json.Unmarshal(apiResponse.Response.Data, &projectDetails)
	if err != nil {
		return
	}

	return
}


/*
	Smartling SDK v2 project api sample

	Sample shows usage of smartling project api
	http://docs.smartling.com/pages/API/v2/Projects/
	Replace userIdentifier and tokenSecret with your credentials
	
*/

package main

import (
	"log"

	// should be replaced with "github.com/Smartling/smartling-cli" 
	// relative path is here only for making it easy to run out of the box
	"./../../smartling"
)

const (
	userIdentifier = "" // put your user identifier here
	tokenSecret    = "" // put your token secret here
	accountId      = "" // put your account id here
)

func main() {

	log.Printf("Initializing smartling client and performing autorization")
	client := smartling.NewClient(userIdentifier, tokenSecret)

	log.Printf("Listing projects for accountId %v:", accountId)

	listRequest := smartling.ProjectListRequest {
		ProjectNameFilter : "",
		IncludeArchived : false,
	}

	projects, err := client.ListProjects(accountId, listRequest)
	if err != nil {
		log.Printf(err.Error())
		return
	}

	log.Printf("Found %v project(s) belonging to user account", projects.TotalCount)

	// now iterate over every project and request it's details
	for _, project := range projects.Items {
		projectDetails, err := client.ProjectDetails(project.ProjectId)
		if err != nil {
			log.Printf(err.Error())
			return
		}
		// print project details struct
		log.Printf("%+v", projectDetails)
	}
}

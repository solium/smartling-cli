// Simple smartling sdk usage example

package main


import (
	"time"
	"log"

	"./smartling"
)

const (
	userIdentifier = ""
	tokenSecret    = ""
	accountId      = ""
	projectId      = ""
)

func main() {

	log.Printf("Initializing smartling client and performing autorization")

	client := smartling.NewClient(userIdentifier, tokenSecret)

	log.Printf("Listing projects:")

	projects, err := client.ListProjects(accountId)
	if err != nil {
		log.Printf("%v", err.Error())
		return
	}
	log.Printf("Success")

	log.Printf("Projects belonging to user account:")
	log.Printf("%+v", projects)

	projectDetails, err := client.ProjectDetails(projectId)
	if err != nil {
		log.Printf("%v", err.Error())
		return
	}
	log.Printf("Success")
	log.Printf("Projects details are")
	log.Printf("%+v", projectDetails)

	for {
		// sleep 6 minutes to issue reauth call
		time.Sleep(time.Minute * 6)
		_, err = client.ListProjects(accountId)
		if err != nil {
			log.Printf("%v", err.Error())
			return
		}
		log.Printf("Success")
	}
}

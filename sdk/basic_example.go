// Simple smartling sdk usage example

package main


import (
	"log"

	"./smartling"
)

const (
	userIdentifier = ""
	tokenSecret    = ""
)

func main() {

	log.Printf("Initializing smartling client and performing autorization")

	client := smartling.NewClient(userIdentifier, tokenSecret)
	err := client.AuthenticationTest()

	if err != nil {
		log.Printf("Authentication Test failed: %v", err.Error())
		return
	}

	log.Printf("Success")

	log.Printf("Listing projects:")

	err = client.ListProjects("")
	if err != nil {
		log.Printf("%v", err.Error())
		return
	}
	log.Printf("Success")
}

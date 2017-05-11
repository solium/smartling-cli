// Simple smartling sdk usage example

package main


import (
	"log"

	"./smartling"
)

const (
	userIdentifier = "gwadsptgphgkdzhqvwmwvithbusrup"
	tokenSecret = "n53mvfl4hct6c357c9eamntgdYt}ddu9frj5h32fdjfrq5jmrrbfg4"
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
}

// class responsible for maintaining user auth valid

package smartling

import (
	"encoding/json"
	"fmt"
	"time"
	"log"
)

// API endpoints
const (
	authApiAuth    = "/auth-api/v2/authenticate"
	authApiRefresh = "/auth-api/v2/authenticate/refresh"
)


// auth api call response data
type AuthApiResponse struct {
	AccessToken      string
	ExpiresIn        int32
	RefreshExpiresIn int32
	RefreshToken     string
}

type Token struct {
	Token string
	ExpirationDate time.Time
}

func (t *Token) IsValid() bool {
	if len(t.Token) == 0 {
		return false
	}

	return time.Now().Before(t.ExpirationDate)
}


type Auth struct {
	userIdentifier  string
	tokenSecret     string
	accessToken     Token
	reauthToken     Token
}

// get access token for request
func (a *Auth) AccessHeader() string {
	return fmt.Sprintf("Bearer %v", a.accessToken.Token)
}

func (a *Auth) authData() ([]byte, error) {

	// prepare the auth map
	authInfo := make(map[string]string)
	authInfo["userIdentifier"] = a.userIdentifier
	authInfo["userSecret"] = a.tokenSecret

	// marshall into bytes
	return json.Marshal(&authInfo)
}

// actually performs auth call
func (a *Auth) doAuthCall(c *Client) error {
	authBytes, err := a.authData()
	if err != nil {
		return err
	}
	// use empty auth header
	bytes, statusCode, err := c.doPostRequest(c.baseUrl + authApiAuth, "", authBytes)

	if err != nil {
		return err
	}

	if statusCode != 200 {
		return fmt.Errorf("Auth call returned unexpected status code: %v", statusCode)
	}

	// unmarshal transport header
	apiResponse := SmartlingApiResponse{}

	err = json.Unmarshal(bytes, &apiResponse)
	if err != nil {
		return err
	}

	// check status
	if apiResponse.Response.Code != "SUCCESS" {
		return fmt.Errorf("Auth call returned unexpected response code: %v", apiResponse.Response.Code)
	}
	log.Printf("auth status %v", statusCode)

	// unmarshal auth body
	authResponse := AuthApiResponse{}

	err = json.Unmarshal(apiResponse.Response.Data, &authResponse)
	if err != nil {
		return err
	}

	// fill tokens
	a.accessToken.Token = authResponse.AccessToken
	a.accessToken.ExpirationDate = time.Now().Add(time.Duration(authResponse.ExpiresIn) * time.Second)

	a.reauthToken.Token = authResponse.RefreshToken
	a.reauthToken.ExpirationDate = time.Now().Add(time.Duration(authResponse.RefreshExpiresIn) * time.Second)

	log.Printf("%#v", authResponse)

	return nil
}

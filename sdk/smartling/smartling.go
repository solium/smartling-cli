// Package smartling is a client implementation of the Smartling Translation API v2 as documented at
// http://docs.smartling.com/pages/API/v2/
package smartling

import (
	"net/http"
	"time"

	"log"
)

// API endpoints
const (
	apiBaseUrl  = "https://api.smartling.com"

	authApiAuth    = "/auth-api/v2/authenticate"
	authApiRefresh = "/auth-api/v2/authenticate/refresh"
)

// Smartling SDK clinet
type Client struct {
	baseUrl         string
	auth            *Auth
	httpClient      *http.Client
}

// Smartling client initialization
func NewClient(userIdentifier string, tokenSecret string) *Client {
	return &Client{
		baseUrl:    apiBaseUrl,
		auth:       &Auth { userIdentifier, tokenSecret },
		httpClient: &http.Client{ Timeout: 60 * time.Second },
	}
}

// custom initialization with overrided base url
func NewClientWithBaseUrl(userIdentifier string, tokenSecret string, baseUrl string) *Client {
	client := NewClient(userIdentifier, tokenSecret)
	client.baseUrl = baseUrl
	return client
}

// attempts to authenticate with smnartling
// not a required call
func (c *Client) AuthenticationTest() error {
	authBytes, err := c.auth.authData()
	if err != nil {
		return err
	}
	bytes, status, err := c.doPostRequest(c.baseUrl + authApiAuth, authBytes)

	if err != nil {
		return err
	}

	log.Printf("auth status %v, response: %#v", status, string(bytes))

	return nil
}

func (c *Client) SetHttpTimeout(t time.Duration) {
	c.httpClient.Timeout = t
}

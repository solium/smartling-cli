// class responsible for maintaining user auth valid

package smartling

import (
	"encoding/json"
)


type Auth struct {
	userIdentifier  string
	tokenSecret     string
}

func (a *Auth) authData() ([]byte, error) {

	// prepare the auth map
	authInfo := make(map[string]string)
	authInfo["userIdentifier"] = a.userIdentifier
	authInfo["userSecret"] = a.tokenSecret

	// marshall into bytes
	return json.Marshal(&authInfo)
}


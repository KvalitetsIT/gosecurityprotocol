package securityprotocol

import (
	"time"

	"encoding/xml"

	saml2 "github.com/russellhaering/gosaml2/types"
)

type SamlSessionDataCreator struct {

	tokenAsString	string

	samlAssertion	*saml2.Assertion
}

func NewSamlSessionDataCreator(assertionString string) (*SamlSessionDataCreator, error) {

	// Parse assertionString to Assertion
    	var assertion saml2.Assertion
        err := xml.Unmarshal([]byte(assertionString), &assertion)
	if (err != nil) {
		return nil, err
	}

	return &SamlSessionDataCreator { tokenAsString: assertionString, samlAssertion: &assertion }, nil
}


func (creator SamlSessionDataCreator) CreateSessionData() (*SessionData, error) {

	userAttributes := make(map[string][]string)

	expiry := time.Now()

	sessionData, err := CreateSessionData(creator.tokenAsString, userAttributes, expiry)

	return sessionData, err
}

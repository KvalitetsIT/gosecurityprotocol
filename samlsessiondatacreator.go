package securityprotocol

import (
	"time"

	"encoding/xml"

	saml2 "github.com/russellhaering/gosaml2/types"
)

type SamlSessionDataCreator struct {
	id string

	samlAssertion *saml2.Assertion

	tokenAsString string

	clientCertificateHash string
}

func NewSamlSessionDataCreator(assertionString string) (*SamlSessionDataCreator, error) {

	return NewSamlSessionDataCreatorWithId("", assertionString)
}

func NewSamlSessionDataCreatorWithId(id string, assertionString string) (*SamlSessionDataCreator, error) {

	// Parse assertionString to Assertion
	var assertion saml2.Assertion
	err := xml.Unmarshal([]byte(assertionString), &assertion)
	if err != nil {
		return nil, err
	}

	return NewSamlSessionDataCreatorWithAssertionAndClientCert(id, assertionString, &assertion, "")
}

func NewSamlSessionDataCreatorWithAssertionAndClientCert(id string, assertionString string, samlAssertion *saml2.Assertion, clientCertificateHash string) (*SamlSessionDataCreator, error) {

	return &SamlSessionDataCreator{id: id, tokenAsString: assertionString, samlAssertion: samlAssertion, clientCertificateHash: clientCertificateHash}, nil
}

func (creator SamlSessionDataCreator) CreateSessionData() (*SessionData, error) {

	userAttributes := make(map[string][]string)
	if creator.samlAssertion != nil && creator.samlAssertion.AttributeStatement != nil {

		for _, samlAttribute := range creator.samlAssertion.AttributeStatement.Attributes {
			userAttributeKey := samlAttribute.Name
			var userAttributeValues []string
			if userAttributes[userAttributeKey] != nil {
				userAttributeValues = userAttributes[userAttributeKey]
			} else {
				userAttributeValues = make([]string, 0)
			}

			for _, samlAttributeValue := range samlAttribute.Values {

				userAttributeValues = append(userAttributeValues, samlAttributeValue.Value)
			}
			userAttributes[userAttributeKey] = userAttributeValues
		}
	}

	expiry, err := time.Parse(time.RFC3339, creator.samlAssertion.Conditions.NotOnOrAfter)
	if err != nil {
		return nil, err
	}

	if creator.id == "" {
		return CreateSessionData(creator.tokenAsString, userAttributes, expiry, creator.clientCertificateHash)
	} else {
		return CreateSessionDataWithId(creator.id, creator.tokenAsString, userAttributes, expiry, creator.clientCertificateHash)
	}
}

package securityprotocol

import (
        "testing"
        "gotest.tools/assert"
	"io/ioutil"
)

func TestSamlSessionDataCreatorWithStringNotSamlAssertionFails(t *testing.T) {

	// Given
	samlassertionAsString := "<not-saml></not-saml>"

        // When
	_, err := NewSamlSessionDataCreator(samlassertionAsString)

	// Then
	assert.Error(t, err, "expected element type <Assertion> but have <not-saml>")
}


func TestSamlSessionDataCreatorWithSamlAssertionSucceedsAndReturnsFullyInitializedSessionData(t *testing.T) {

        // Given
        samlassertionAsBytes, _ := ioutil.ReadFile("./testdata/saml-assertion-test.xml")

        // When
        samlSessionDataCreator, errSessionDataCreator := NewSamlSessionDataCreator(string(samlassertionAsBytes))
	samlSessionData, errSessionData := samlSessionDataCreator.CreateSessionData()

        // Then
        assert.NilError(t, errSessionDataCreator)
	assert.NilError(t, errSessionData)

	assert.Equal(t, string(samlassertionAsBytes), samlSessionData.Token)

	assert.Equal(t, len(samlSessionData.UserAttributes), 3)

	assert.Equal(t, len(samlSessionData.UserAttributes["eduPersonAffiliation"]), 2)
	assert.Equal(t, samlSessionData.UserAttributes["eduPersonAffiliation"][0], "users")
	assert.Equal(t, samlSessionData.UserAttributes["eduPersonAffiliation"][1], "examplerole1")

        assert.Equal(t, len(samlSessionData.UserAttributes["mail"]), 1)
        assert.Equal(t, samlSessionData.UserAttributes["mail"][0], "test@example.com")

        assert.Equal(t, len(samlSessionData.UserAttributes["uid"]), 1)
        assert.Equal(t, samlSessionData.UserAttributes["uid"][0], "test")
}


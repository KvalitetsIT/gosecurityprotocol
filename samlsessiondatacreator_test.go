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
}


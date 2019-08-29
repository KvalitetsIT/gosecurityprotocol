package securityprotocol

import (
        "testing"
	"net/http"

        "gotest.tools/assert"
)

type MockTokenCache struct {

}

func (tokenCache *MockTokenCache) FindTokenDataForSessionId(sessionId string) (*TokenData, error) {
        result := TokenData{ Sessionid: sessionId, Hash: "hash" }
        return &result, nil
}

func (tokenCache *MockTokenCache) SaveAuthenticationKeysForSessionId(sessionId string, authenticationToken string, expires_in int64, hash string) (*TokenData, error) {

	expiryTime := GetExpiryDate(expires_in)
	tokenData := &TokenData{ Sessionid: sessionId, Authenticationtoken: authenticationToken, Timestamp: expiryTime, Hash: hash  }
	return tokenData, nil
}


type MockService struct {

}

func (mock *MockService) Handle(http.ResponseWriter, *http.Request) (int, error) {
	return http.StatusOK, nil
}


func TestExampleProtocolAnswersUnautorizedWhenNoSessionIdCanBeFound(t *testing.T) {

	// Given
	service := new(MockService)
	tokenCache := new(MockTokenCache)
	exampleClientProtocol := NewExampleClientProtocol(tokenCache, service)
        req, _ := http.NewRequest("GET", "/someurl", nil)

        // When
	httpCode, err := exampleClientProtocol.Handle(nil, req)

	// Then
	assert.NilError(t, err)
	assert.Equal(t, http.StatusUnauthorized, httpCode)
}

func TestExampleProtocolStartsAutorizationIfNoTokenMatchingTheSessionIdCanBeFound(t *testing.T) {

        // Given
        service := new(MockService)
        tokenCache := new(MockTokenCache)

	authenticationCalled := false
	authenticationHook := func(sessionData *SessionData) (*ClientAuthenticationInfo, error) {
		authenticationCalled = true
		return ExampleDoAuthenticationHook(sessionData)
	}

        exampleClientProtocol := NewExampleClientProtocolWithHooks(tokenCache, service, authenticationHook)
        req, _ := http.NewRequest("GET", "/someurl", nil)


        // When
        httpCode, err := exampleClientProtocol.Handle(nil, req)

        // Then
        assert.NilError(t, err)
        assert.Equal(t, http.StatusUnauthorized, httpCode)
	assert.Equal(t, true, authenticationCalled)
}


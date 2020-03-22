package securityprotocol

import (
        "testing"
	"net/http"
	"fmt"
	"time"
        "gotest.tools/assert"
)

type MockTokenCache struct {

}

func (tokenCache *MockTokenCache) FindTokenDataForSessionId(sessionId string) (*TokenData, error) {
        result := TokenData{ Sessionid: sessionId, Hash: "hash" }
        return &result, nil
}

func (tokenCache *MockTokenCache) SaveAuthenticationKeysForSessionId(sessionId string, authenticationToken string, expires_in int64, hash string) (*TokenData, error) {

	expiryTime := time.Now() //GetExpiryDate(expires_in)
	return tokenCache.SaveAuthenticationKeysForSessionIdWithExpiry(sessionId, authenticationToken, expiryTime, hash)
}


func (tokenCache *MockTokenCache) SaveAuthenticationKeysForSessionIdWithExpiry(sessionId string, authenticationToken string, expiryTime time.Time, hash string) (*TokenData, error) {
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

	authenticationHook := func(writer http.ResponseWriter, request *http.Request, sessionData *SessionData) (*ClientAuthenticationInfo, int, error) {
		authenticationCalled = true
		return ExampleDoAuthenticationHook(writer, request, sessionData)
	}

        exampleClientProtocol := NewExampleClientProtocolWithHooks(ExampleMatchHandler, tokenCache, service, ExamplePreAuthentication, authenticationHook)
        req, _ := http.NewRequest("GET", "/someurl", nil)
	sessionId := "session-123-xyz-999999"
	ExampleAddSessionIdToRequest(req, sessionId)

        // When
        httpCode, err := exampleClientProtocol.Handle(nil, req)

        // Then
        assert.NilError(t, err)
        assert.Equal(t, http.StatusOK, httpCode)
	assert.Equal(t, true, authenticationCalled)
}

func TestExampleProtocolSkipsAuthenticationWhenPreAuthenticationCausesRedirect(t *testing.T) {

        // Given
        service := new(MockService)
        tokenCache := new(MockTokenCache)

        authenticationCalled := false
        authenticationHook := func(w http.ResponseWriter, r *http.Request, sessionData *SessionData) (*ClientAuthenticationInfo, int, error) {
                authenticationCalled = true
                return ExampleDoAuthenticationHook(w, r, sessionData)
        }

	preAuthErrMsg := "Redirecting because we want to"
	examplePreAuthentication := func(w http.ResponseWriter, r *http.Request, sessionData *SessionData) (int, error) {
		return http.StatusTemporaryRedirect, fmt.Errorf(preAuthErrMsg)
	}

        exampleClientProtocol := NewExampleClientProtocolWithHooks(ExampleMatchHandler, tokenCache, service, examplePreAuthentication, authenticationHook)
        req, _ := http.NewRequest("GET", "/someurl", nil)
        sessionId := "session-123-xyz-999999"
        ExampleAddSessionIdToRequest(req, sessionId)

        // When
        httpCode, err := exampleClientProtocol.Handle(nil, req)

        // Then
        assert.Error(t, err, preAuthErrMsg)
        assert.Equal(t, http.StatusTemporaryRedirect, httpCode)
        assert.Equal(t, false, authenticationCalled)
}

func TestExampleProtocolSkipsUrlsThatShouldBeIgnored(t *testing.T) {

        // Given
        service := new(MockService)
        tokenCache := new(MockTokenCache)
        exampleClientProtocol := NewExampleClientProtocol(tokenCache, service)
        req, _ := http.NewRequest("GET", "/someurl?skip=yes", nil)

        // When
        httpCode, err := exampleClientProtocol.Handle(nil, req)

        // Then
        assert.NilError(t, err)
        assert.Equal(t, http.StatusOK, httpCode)
}

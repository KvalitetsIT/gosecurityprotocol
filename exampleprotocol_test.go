package securityprotocol

import (
        "testing"
	"net/http"

//        "gotest.tools/assert"
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


func TestExampleProtocolWhenNoTokenIsFound(t *testing.T) {

	// Given
	service := new(MockService)
	tokenCache := new(MockTokenCache)
	exampleClientProtocol := NewExampleClientProtocol(tokenCache, service)

        // When
	exampleClientProtocol.Handle(nil, nil)

	// Then
}

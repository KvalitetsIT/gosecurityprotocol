package securityprotocol

import (
        "net/http"
)

const EXAMPLEPROTOCOL_HEADER_NAME = "examplesession"

type ExampleSessionDataFetcher struct {

	authorizationHook func(sessionData *SessionData) (*ClientAuthenticationInfo, error)

}

func (fetcher ExampleSessionDataFetcher) GetSessionData(sessionId string, sessionIdHandler SessionIdHandler)  (*SessionData, error) {

	sessionData := SessionData{ SessionId: sessionId, Hash: sessionId }

	return &sessionData, nil
}


func NewExampleClientProtocol(tokenCache TokenCache, service HttpHandler) (*HttpProtocolClient) {

	return NewExampleClientProtocolWithHooks(ExampleMatchHandler, tokenCache, service, ExamplePreAuthentication, ExampleDoAuthenticationHook)
}

func NewExampleClientProtocolWithHooks(matchHandler MatchHandler, tokenCache TokenCache, service HttpHandler, preAuthentication func(w http.ResponseWriter, r *http.Request, sessionData *SessionData) (int, error),clientAuthenticationInfo func(sessionData *SessionData) (*ClientAuthenticationInfo, error)) (*HttpProtocolClient) {

        sessionIdHandler := &HttpHeaderSessionIdHandler{ HttpHeaderName: EXAMPLEPROTOCOL_HEADER_NAME }

        sessionDataFetcher := new(ExampleSessionDataFetcher)

        protocolClient := NewHttpProtocolClient(matchHandler, tokenCache, sessionIdHandler, sessionDataFetcher, preAuthentication, clientAuthenticationInfo, ExampleDecorateRequestWithAuthenticationToken, service)

        return protocolClient
}


func ExamplePreAuthentication(w http.ResponseWriter, r *http.Request, sessionData *SessionData) (int, error) {

	return http.StatusTeapot, nil
}


func ExampleDoAuthenticationHook(sessionData *SessionData) (*ClientAuthenticationInfo, error) {

        // Default implementation
        mock := new(ClientAuthenticationInfo)
        mock.Token = "mock"
        mock.ExpiresIn = 2000

        return mock, nil
}

func ExampleDecorateRequestWithAuthenticationToken(tokenData *TokenData, r *http.Request) error {

        // TODO
        return nil
}

func ExampleAddSessionIdToRequest(r *http.Request, sessionId string) {
	r.Header.Add(EXAMPLEPROTOCOL_HEADER_NAME, sessionId)
}

// Example matcher matcher alt undtagen request med en 'skip' query parameter
func ExampleMatchHandler (r *http.Request) bool {

	keys, ok := r.URL.Query()["skip"]
	if (ok && len(keys) > 0) {
		return false
	}

        return true
}


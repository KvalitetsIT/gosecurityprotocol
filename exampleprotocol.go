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

	return NewExampleClientProtocolWithHooks(tokenCache, service, ExampleDoAuthenticationHook)
}

func NewExampleClientProtocolWithHooks(tokenCache TokenCache, service HttpHandler, clientAuthenticationInfo func(sessionData *SessionData) (*ClientAuthenticationInfo, error)) (*HttpProtocolClient) {

        sessionIdHandler := &HttpHeaderSessionIdHandler{ HttpHeaderName: EXAMPLEPROTOCOL_HEADER_NAME }

        sessionDataFetcher := new(ExampleSessionDataFetcher)

        protocolClient := NewHttpProtocolClient(tokenCache, sessionIdHandler, sessionDataFetcher, clientAuthenticationInfo, ExampleDecorateRequestWithAuthenticationToken, service)

        return protocolClient
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


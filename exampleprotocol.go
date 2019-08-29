package securityprotocol

import (
        "net/http"
)

type ExampleSessionIdHandler struct {

}

func (handler ExampleSessionIdHandler) GetSessionIdFromHttpRequest(request *http.Request) string {

        return "test"
}

func (handler ExampleSessionIdHandler) SetSessionIdOnHttpRequest(sessionId string, request *http.Request)  {

}


type ExampleSessionDataFetcher struct {

}

func (fetcher ExampleSessionDataFetcher) GetSessionData(sessionId string, sessionIdHandler SessionIdHandler)  (*SessionData, error) {

	sessionData := SessionData{ SessionId: sessionId, Hash: sessionId }

	return &sessionData, nil
}


func NewExampleClientProtocol(tokenCache TokenCache, service HttpHandler) (*HttpProtocolClient) {

        sessionIdHandler := new(ExampleSessionIdHandler)

        sessionDataFetcher := new(ExampleSessionDataFetcher)

        protocolClient := NewHttpProtocolClient(tokenCache, sessionIdHandler, sessionDataFetcher, ExampleDoAuthenticationHook, ExampleDecorateRequestWithAuthenticationToken, service)

        return protocolClient
}


func ExampleDoAuthenticationHook(sessionData *SessionData) (*ClientAuthenticationInfo, error) {

        // TODO
        mock := new(ClientAuthenticationInfo)
        mock.Token = "mock"
        mock.ExpiresIn = 2000

        return mock, nil
}

func ExampleDecorateRequestWithAuthenticationToken(tokenData *TokenData, r *http.Request) error {

        // TODO
        return nil
}


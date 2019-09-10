package securityprotocol

import (
        "net/http"
	"fmt"
)


type ClientAuthenticationInfo struct {

        Token           string
        ExpiresIn       int64
}

type DoClientAuthentification func(*SessionData) (*ClientAuthenticationInfo, error)

type DecorateRequestWithAuthenticationToken func(tokenData *TokenData, r *http.Request) error

type PreAuthentication func(w http.ResponseWriter, r *http.Request, sessionData *SessionData) (int, error)

type HttpProtocolClient struct {

	tokenCache      	TokenCache

	sessionIdHandler	SessionIdHandler
	sessionDataFetcher	SessionDataFetcher

	preAuthentication	PreAuthentication
	doClientAuthentication	DoClientAuthentification
	decorateRequest		DecorateRequestWithAuthenticationToken

	service			HttpHandler
}

func NewHttpProtocolClient(tokenCache TokenCache, sessionIdHandler SessionIdHandler, sessionDataFetcher SessionDataFetcher, preAuthentication PreAuthentication, doClientAuthentication DoClientAuthentification, decorateRequest DecorateRequestWithAuthenticationToken, service HttpHandler) (*HttpProtocolClient) {

	httpProtocolClient := new (HttpProtocolClient)
	httpProtocolClient.tokenCache = tokenCache
	httpProtocolClient.sessionIdHandler = sessionIdHandler
	httpProtocolClient.sessionDataFetcher = sessionDataFetcher
	httpProtocolClient.preAuthentication = preAuthentication
	httpProtocolClient.doClientAuthentication = doClientAuthentication
	httpProtocolClient.decorateRequest = decorateRequest
	httpProtocolClient.service = service

	return httpProtocolClient
}

func (client HttpProtocolClient) Handle(w http.ResponseWriter, r *http.Request) (int, error) {

	// Check for session id
	sessionId := client.sessionIdHandler.GetSessionIdFromHttpRequest(r)
	if (sessionId == "") {
		return http.StatusUnauthorized, nil
        }

	// Check if we have a token cached
	tokenData, err := client.tokenCache.FindTokenDataForSessionId(sessionId)
	if (err != nil) {
		fmt.Println(fmt.Sprintf("Error in FindTokenDataForSessionId: %s (error:%v)", sessionId, err))
		return http.StatusInternalServerError, err
	}

	// Get sessiondata
	sessionData, err := client.sessionDataFetcher.GetSessionData(sessionId, client.sessionIdHandler)
        if (err != nil) {
		fmt.Println(fmt.Sprintf("Error in GetSessionData: %s (error:%v)", sessionId, err))
                return http.StatusInternalServerError, err
        }

	if (tokenData == nil || (tokenData.Hash != sessionData.Hash)) {

		if (client.preAuthentication != nil) {
			httpCode, err := client.preAuthentication(w, r, sessionData)
			if (err != nil) {
				return httpCode, err
			}
		}

		// No token - or sessiondata has changed since issueing - run authentication (again)
		authentication, err := client.doClientAuthentication(sessionData)
		if (err != nil) {
			return http.StatusUnauthorized, nil
		}

		tokenData, err = client.tokenCache.SaveAuthenticationKeysForSessionId(sessionId, authentication.Token, authentication.ExpiresIn, sessionData.Hash)
		if (err != nil) {
                        return http.StatusUnauthorized, nil
                }
	}

	// Add the authentication token to the request
	client.decorateRequest(tokenData, r)

	// Let the service do its work
        return client.service.Handle(w, r)
}

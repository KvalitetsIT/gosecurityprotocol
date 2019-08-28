package securityprotocol

import (
        "net/http"
)


type ClientAuthenticationInfo struct {

        Token           string
        ExpiresIn       int64
}

type DoClientAuthentification func(*SessionData) (*ClientAuthenticationInfo, error)

type DecorateRequestWithAuthenticationToken func(tokenData *TokenData, r *http.Request) error

type HttpProtocolClient struct {

	tokenCache      	*TokenCache

	sessionIdHandler	SessionIdHandler
	sessionDataFetcher	SessionDataFetcher

	doClientAuthentication	DoClientAuthentification
	decorateRequest		DecorateRequestWithAuthenticationToken

	service			HttpHandler
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
		return http.StatusInternalServerError, err
	}

	// Get sessiondata
	sessionData, err := client.sessionDataFetcher(sessionId, client.sessionIdHandler)
        if (err != nil) {
                return http.StatusInternalServerError, err
        }

	if (tokenData == nil || (tokenData.Hash != sessionData.Hash)) {
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

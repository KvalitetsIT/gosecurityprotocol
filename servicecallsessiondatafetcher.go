package securityprotocol

import (
	"fmt"
	"net/http"
	"bytes"
	"encoding/json"
)

type SessionDataDto struct {

	SessionData

	SamlToken	string
}

type ServiceCallSessionDataFetcher struct {
	sessionDataServiceEndpoint string
	client *http.Client
}

func NewServiceCallSessionDataFetcher(sessionDataServiceEndpoint string, client *http.Client) *ServiceCallSessionDataFetcher {

	return &ServiceCallSessionDataFetcher{ sessionDataServiceEndpoint: sessionDataServiceEndpoint, client: client}
}

func (fetcher ServiceCallSessionDataFetcher) GetSessionData(sessionId string, sessionIdHandler SessionIdHandler)  (*SessionData, error) {

	// Create request
        req, err := http.NewRequest("GET", fmt.Sprintf("%s/getsessiondata", fetcher.sessionDataServiceEndpoint), nil)
        if (err != nil) {
                return nil, err
        }
	sessionIdHandler.SetSessionIdOnHttpRequest(sessionId, req)

	// Make call
        resp, err := fetcher.client.Do(req)
        if (err != nil) {
                return nil, err
        }

	// Parse response
        buffer := new(bytes.Buffer)
        buffer.ReadFrom(resp.Body)
        var result SessionDataDto
        err = json.Unmarshal(buffer.Bytes(), &result)
        if (err != nil) {
                return nil, err
        }
	if (len(result.SamlToken) > 0) {
		result.SessionData.Authenticationtoken =  result.SamlToken
	}

        return &(result.SessionData), nil
}

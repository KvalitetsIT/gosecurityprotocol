package securityprotocol

import (
	"fmt"
	"net/http"
	"bytes"
	"encoding/json"
)

type ServiceCallSessionDataFetcher struct {
	SessionDataServiceEndpoint string
}

func (fetcher ServiceCallSessionDataFetcher) GetSessionData(sessionId string, sessionIdHandler SessionIdHandler)  (*SessionData, error) {

	client := &http.Client{}

	// Create request
        req, err := http.NewRequest("GET", fmt.Sprintf("%s/getsessiondata", fetcher.SessionDataServiceEndpoint), nil)
        if (err != nil) {
                return nil, err
        }
	sessionIdHandler.SetSessionIdOnHttpRequest(sessionId, req)

	// Make call
        resp, err := client.Do(req)
        if (err != nil) {
                return nil, err
        }

	// Parse response
        buffer := new(bytes.Buffer)
        buffer.ReadFrom(resp.Body)
        var result SessionData
        err = json.Unmarshal(buffer.Bytes(), &result)
        if (err != nil) {
                return nil, err
        }

        return &result, nil
}

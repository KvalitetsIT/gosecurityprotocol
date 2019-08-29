package securityprotocol

import (
	"fmt"
	"time"
	"net/http"
	"bytes"
	"encoding/json"
)

type SessionData struct {

	SessionId         string
	Token         	  string
	UserAttributes    map[string][]string
	SessionAttributes map[string]string
	Timestamp         time.Time
	Hash              string
}

type SessionDataFetcher interface {
	GetSessionData(string, SessionIdHandler) (*SessionData, error)
}


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

package securityprotocol

import (
        "testing"
        "gotest.tools/assert"
	"io/ioutil"
//	uuid "github.com/google/uuid"
	httptest "net/http/httptest"
	http "net/http"
//	"fmt"
)


func TestSessionDataHandlerGetSessionDataOldInterface(t *testing.T) {

        // Given
        sessionId := "sessionId1"
        hash := "hash1"
        sessionData := SessionData{ Sessionid: sessionId, Hash: hash }
	getSessionDataRequest := httptest.NewRequest(http.MethodGet, "/getsessiondata", nil)
	responseRecorder := httptest.NewRecorder()

	// When
	getHandlerFunction := IsRequestForSessionData(&sessionData, nil, responseRecorder, getSessionDataRequest)
	var getStatus int
	var getError error
	if (getHandlerFunction != nil) {
		getStatus, getError = getHandlerFunction()
	}

        // Then
	assert.Assert(t, getHandlerFunction != nil)
	assert.NilError(t, getError)
	assert.Equal(t, getStatus, http.StatusOK)
	assert.Assert(t, responseRecorder.Result() != nil)
	assert.Equal(t, "{\"ID\":\"\",\"Sessionid\":\"sessionId1\",\"Authenticationtoken\":\"\",\"Timestamp\":\"0001-01-01T00:00:00Z\",\"Hash\":\"hash1\",\"UserAttributes\":null,\"SessionAttributes\":null,\"ClientCertHash\":\"\"}", readStringFromResponseBody(responseRecorder.Result()))
}

func TestSessionDataHandlerGetSessionData(t *testing.T) {

        // Given
        sessionId := "sessionId1"
        hash := "hash1"
        sessionData := SessionData{ Sessionid: sessionId, Hash: hash }
	getSessionDataRequest := httptest.NewRequest(http.MethodGet, "/sessiondata", nil)
	responseRecorder := httptest.NewRecorder()

	// When
	getHandlerFunction := IsRequestForSessionData(&sessionData, nil, responseRecorder, getSessionDataRequest)
	var getStatus int
	var getError error
	if (getHandlerFunction != nil) {
		getStatus, getError = getHandlerFunction()
	}

        // Then
	assert.Assert(t, getHandlerFunction != nil)
	assert.NilError(t, getError)
	assert.Equal(t, getStatus, http.StatusOK)
	assert.Assert(t, responseRecorder.Result() != nil)
	assert.Equal(t, "{\"ID\":\"\",\"Sessionid\":\"sessionId1\",\"Authenticationtoken\":\"\",\"Timestamp\":\"0001-01-01T00:00:00Z\",\"Hash\":\"hash1\",\"UserAttributes\":null,\"SessionAttributes\":null,\"ClientCertHash\":\"\"}", readStringFromResponseBody(responseRecorder.Result()))

}

func TestSessionDataHandlerSetSessionAttributeOldInterface(t *testing.T) {

        // Given
        mongoSessionCache, createErr := NewMongoSessionCache("mongo", "testsessionapi", "session")

        sessionId := "session2"
        hash1 := "hash1"
	sessionKey := "key_my.kuk.æøå"
	sessionValue := "value1_./&blaÅÆØ"
        sessionData := SessionData{ Sessionid: sessionId, Hash: hash1 }
        saveErr := mongoSessionCache.SaveSessionData(&sessionData)
        setSessionAttributeRequest := httptest.NewRequest(http.MethodGet, "/setsessionattribute", nil)
	q := setSessionAttributeRequest.URL.Query()
    	q.Add("key", sessionKey)
	q.Add("value", sessionValue)
	setSessionAttributeRequest.URL.RawQuery = q.Encode()
	responseRecorder := httptest.NewRecorder()

        // When
	setHandlerFunction := IsRequestForSessionData(&sessionData, mongoSessionCache, responseRecorder, setSessionAttributeRequest)
        var setStatus int
        var setError error
        if (setHandlerFunction != nil) {
                setStatus, setError = setHandlerFunction()
        }
	gotSessionData, _ := mongoSessionCache.FindSessionDataForSessionId(sessionId)

        // Then
        assert.NilError(t, createErr)
        assert.NilError(t, saveErr)
        assert.Assert(t, setHandlerFunction != nil)
        assert.NilError(t, setError)
        assert.Equal(t, setStatus, http.StatusOK)
        assert.Assert(t, responseRecorder.Result() != nil)
	assert.Assert(t, gotSessionData != nil)
	assert.Equal(t, len(gotSessionData.SessionAttributes), 1)
	assert.Equal(t, gotSessionData.SessionAttributes[sessionKey], sessionValue)
//	jsonSessionData, _ := gotSessionData.ToString()
//	assert.Equal(t, jsonSessionData, "kuk")
 	assert.Equal(t, gotSessionData.Hash, "dNVCwG36HBY+n2vCoAw1+g==")
}

func readStringFromResponseBody(resp *http.Response) string {

	bodyBytes, _ := ioutil.ReadAll(resp.Body)
	bodyString := string(bodyBytes)
	return bodyString
}

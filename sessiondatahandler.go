package securityprotocol

import (
	"fmt"
	"net/http"
	"encoding/json"
)


type SessionDataHandlerFunction = func() (int, error)


func IsRequestForSessionData(sessionData *SessionData, sessionCache SessionCache, w http.ResponseWriter, r *http.Request) (SessionDataHandlerFunction) {

        path := r.URL.Path
        if ((path == "/getsessiondata") && (http.MethodGet == r.Method)) {
                // Old Interface ... will be removed in the future
                return func() (int,error) {
                        return handleRequestForSessionData(sessionData, w, r)
                }
        }
        if (path == "/setsessionattribute" && (http.MethodGet == r.Method)) {
                // Old Interface ... will be removed in the future
                return func() (int, error) {
                        return handleUpdateSessionDataFromQueryParameters(sessionData, sessionCache, w, r)
                }
        }

        if (path == "/sessiondata" && (http.MethodGet == r.Method)) {
                return func() (int, error) {
                        return handleRequestForSessionData(sessionData, w, r)
                }

        }
        return nil
}

func handleUpdateSessionDataFromQueryParameters(sessionData *SessionData, sessionCache SessionCache, w http.ResponseWriter, r *http.Request) (int, error) {

        keys, keyOk := r.URL.Query()["key"]
        values, valOk := r.URL.Query()["value"]

        if (!keyOk || !valOk || !(len(keys) == len(values))) {
                return http.StatusBadRequest, fmt.Errorf("Legal request for update sessiondata includes matching sets of key and value query parameters")
        }

        for index, key := range keys {

                value := values[index]
                sessionData.AddSessionAttribute(key, value)
        }
        err := sessionCache.SaveSessionData(sessionData)
        if (err != nil) {
                return http.StatusInternalServerError, err
        }
        return http.StatusOK, nil
}

func handleRequestForSessionData(sessionData *SessionData, w http.ResponseWriter, r *http.Request) (int, error) {

        sessionDataBytes, marshalErr := json.Marshal(sessionData)
        if (marshalErr != nil) {
                return http.StatusInternalServerError, marshalErr
        }
	w.Header().Set("Content-Type", "application/json")
        w.Write(sessionDataBytes)

        return http.StatusOK, nil
}

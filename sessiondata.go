package securityprotocol

import (
	"time"
	"sort"
	"crypto/md5"
	"io"
	"encoding/base64"
	uuid "github.com/google/uuid"

	"encoding/json"
	"gopkg.in/mgo.v2/bson"
)

type SessionData struct {

        ID                      bson.ObjectId `bson:"_id,omitempty"`
        Sessionid               string `bson:"sessionid"`
        Authenticationtoken     string
        Timestamp               time.Time `bson:"timestamp"`
        Hash                    string

	UserAttributes    map[string][]string
	SessionAttributes map[string]string

	ClientCertHash string
}

type SessionDataFetcher interface {
	GetSessionData(string, SessionIdHandler) (*SessionData, error)
}

type SessionDataCreator interface {
	CreateSessionData() (*SessionData, error)
}

type SessionCache interface {
        SaveSessionData(*SessionData) error
        FindSessionDataForSessionId(sessionId string) (*SessionData, error)
}

func CreateSessionDataWithId(id string, token string, userAttributes map[string][]string, expiry time.Time, clientCertHash string) (*SessionData, error) {

        sessionData := SessionData { Sessionid: id, Authenticationtoken: token, Timestamp: expiry, UserAttributes: userAttributes, SessionAttributes: make(map[string]string) }
        sessionData.Hash = sessionData.CalculateHash()
	sessionData.ClientCertHash = clientCertHash

        return &sessionData, nil
}

func CreateSessionData(token string, userAttributes map[string][]string, expiry time.Time, clientCertHash string) (*SessionData, error) {

	id := uuid.New().String()
	return CreateSessionDataWithId(id, token, userAttributes, expiry, clientCertHash)
}

func (data *SessionData) AddSessionAttribute(key string, value string) {

	if (data.SessionAttributes == nil) {
		data.SessionAttributes = make(map[string]string)
	}
	data.SessionAttributes[key] = value
	data.Hash = data.CalculateHash()
}

func (data *SessionData) CalculateHash() string {

        s := data.Sessionid
        s = s + data.Authenticationtoken
        s = s + data.Timestamp.Format(time.UnixDate)

	userAttributeKeys := []string{}
	for k, _ := range data.UserAttributes {
		userAttributeKeys = append(userAttributeKeys, k)
	}
	sort.Strings(userAttributeKeys)
	for _, k := range userAttributeKeys {
		s = s + k
		for _, v := range data.UserAttributes[k] {
			s = s + v
		}
	}

	sessionAttributeKeys := []string{}
	for k, _ := range data.SessionAttributes {
		sessionAttributeKeys = append(sessionAttributeKeys, k)
	}
	sort.Strings(sessionAttributeKeys)
	for _, k := range sessionAttributeKeys {
		s = s + k + data.SessionAttributes[k]
	}

	h := md5.New()
	io.WriteString(h, s)
	hash := base64.StdEncoding.EncodeToString(h.Sum(nil))

	return hash
}

func (data *SessionData) ToString() (string, error) {
	sessionDataBytes, marshalErr := json.Marshal(data)
	if (marshalErr != nil) {
		return "", marshalErr
	}
	return string(sessionDataBytes), nil
}

type NilSessionDataFetcher struct {
}

func (fetcher NilSessionDataFetcher) GetSessionData(sessionId string, sessionIdHandler SessionIdHandler)  (*SessionData, error) {
	return nil, nil
}

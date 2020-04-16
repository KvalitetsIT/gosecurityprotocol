package securityprotocol

import (
	"crypto/md5"
	"encoding/base64"
	"encoding/json"
	uuid "github.com/google/uuid"
	primitive "go.mongodb.org/mongo-driver/bson/primitive"
	"io"
	"sort"
	"time"
)


type Identifiable interface {
	GetID() *primitive.ObjectID
	GetKey() string
}

type SessionData struct {
	identifiable Identifiable
	ID                  *primitive.ObjectID `bson:"_id,omitempty"`
	Sessionid           string        `bson:"sessionid"`
	Authenticationtoken string
	Timestamp           time.Time `bson:"timestamp"`
	Hash                string

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
	DeleteSessionData(sessionId string) error
}

func CreateSessionDataWithId(id string, token string, userAttributes map[string][]string, expiry time.Time, clientCertHash string) (*SessionData, error) {
	encodedToken := base64.StdEncoding.EncodeToString([]byte(token))
	sessionData := SessionData{Sessionid: id, Authenticationtoken: encodedToken, Timestamp: expiry, UserAttributes: userAttributes, SessionAttributes: make(map[string]string)}
	sessionData.Hash = sessionData.CalculateHash()
	sessionData.ClientCertHash = clientCertHash

	return &sessionData, nil
}

func CreateSessionData(token string, userAttributes map[string][]string, expiry time.Time, clientCertHash string) (*SessionData, error) {

	id := uuid.New().String()
	return CreateSessionDataWithId(id, token, userAttributes, expiry, clientCertHash)
}

func (data SessionData) GetID() *primitive.ObjectID {
	return data.ID
}

func (data SessionData) GetKey() string {
	return data.Sessionid
}

func (data *SessionData) AddSessionAttribute(key string, value string) {

	if data.SessionAttributes == nil {
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
	if marshalErr != nil {
		return "", marshalErr
	}
	return string(sessionDataBytes), nil
}

type NilSessionDataFetcher struct {
}

func (fetcher NilSessionDataFetcher) GetSessionData(sessionId string, sessionIdHandler SessionIdHandler) (*SessionData, error) {
	return nil, nil
}

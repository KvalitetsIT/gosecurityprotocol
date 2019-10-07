package securityprotocol

import (
	"time"
	"sort"
	"crypto/md5"
	"io"
	"encoding/base64"
	uuid "github.com/google/uuid"
)

type SessionData struct {

	TokenData

	UserAttributes    map[string][]string
	SessionAttributes map[string]string
}

type SessionDataFetcher interface {
	GetSessionData(string, SessionIdHandler) (*SessionData, error)
}

type SessionDataCreator interface {
	CreateSessionData() (*SessionData, error)
}

type SessionCache interface {
        SaveAuthenticationKeysForSessionId(sessionId string, authenticationToken string, expires_in int64, hash string) (*SessionData, error)
        FindSessionDataForSessionId(sessionId string) (*SessionData, error)
}



func CreateSessionDataWithId(id string, token string, userAttributes map[string][]string, expiry time.Time) (*SessionData, error) {

        sessionData := SessionData { TokenData: TokenData { Sessionid: id, Authenticationtoken: token, Timestamp: expiry}, UserAttributes: userAttributes, SessionAttributes: make(map[string]string) }
        sessionData.Hash = sessionData.CalculateHash()

        return &sessionData, nil
}

func CreateSessionData(token string, userAttributes map[string][]string, expiry time.Time) (*SessionData, error) {

	id := uuid.New().String()
	return CreateSessionDataWithId(id, token, userAttributes, expiry)
}

func (data *SessionData) AddSessionAttribute(key string, value string) {

	data.SessionAttributes[key] = value
	data.Hash = data.CalculateHash()
}

func (data *SessionData) CalculateHash() string {

	s := data.TokenData.CalculateHash()

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

type NilSessionDataFetcher struct {
}

func (fetcher NilSessionDataFetcher) GetSessionData(sessionId string, sessionIdHandler SessionIdHandler)  (*SessionData, error) {
	return nil, nil
}


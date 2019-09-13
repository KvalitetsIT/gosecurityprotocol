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

type SessionDataCreator interface {
	CreateSessionData() (*SessionData, error)
}



func CreateSessionDataWithId(id string, token string, userAttributes map[string][]string, expiry time.Time) (*SessionData, error) {

        sessionData := SessionData { SessionId: id, Token: token, UserAttributes: userAttributes, SessionAttributes: make(map[string]string), Timestamp: expiry }
        sessionData.recalculateHash()

        return &sessionData, nil
}



func CreateSessionData(token string, userAttributes map[string][]string, expiry time.Time) (*SessionData, error) {

	id := uuid.New().String()
	return CreateSessionDataWithId(id, token, userAttributes, expiry)
}

func (data *SessionData) AddSessionAttribute(key string, value string) {

	data.SessionAttributes[key] = value
	data.recalculateHash()
}

func (data *SessionData) recalculateHash() string {

	s := data.SessionId
	s = s + data.Token
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
	data.Hash = base64.StdEncoding.EncodeToString(h.Sum(nil))

	return data.Hash
}

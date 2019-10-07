package securityprotocol

import "fmt"

type MongoSessionCache struct {

	MongoCache	*MongoCache
}

func NewMongoSessionCache(mongodb string, mongodb_database string, mongodb_collection string) (*MongoSessionCache, error) {
	mongoCache, err := NewMongoCache(mongodb, mongodb_database, mongodb_collection, "token")
	if (err != nil) {
		return nil, err
	}
	return &MongoSessionCache{ MongoCache: mongoCache }, nil
}

func (sessionCache *MongoSessionCache) FindSessionDataForSessionId(sessionId string) (*SessionData, error) {
	if (sessionId == "") {
		return nil, fmt.Errorf("Session id cannot be empty")
	}

	result := SessionData{}
	err := sessionCache.MongoCache.FindDataForSessionId("sessionid", sessionId, &result)
	if (err != nil) {
		return nil, err
	}

	return &result, nil
}

func (sessionCache *MongoSessionCache) SaveAuthenticationKeysForSessionId(sessionId string, authenticationToken string, expires_in int64, hash string) (*SessionData, error) {
	if (sessionId != "") {
		existing, _ := sessionCache.FindSessionDataForSessionId(sessionId)
		if (existing != nil) {
			sessionCache.MongoCache.Delete(existing)
		}

               	expiryTime := GetExpiryDate(expires_in)
		sessionData := &SessionData{ TokenData: TokenData{ Sessionid: sessionId, Authenticationtoken: authenticationToken, Timestamp: expiryTime, Hash: hash  } }
		err := sessionCache.MongoCache.Save(sessionData)
		if (err != nil) {
			return nil, err
		}
		return sessionData, nil
	}
	return nil, fmt.Errorf("sessionId cannot be empty")
}

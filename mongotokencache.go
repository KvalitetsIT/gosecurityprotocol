package securityprotocol

import "fmt"

type MongoTokenCache struct {
	MongoCache	*MongoCache
}

func NewMongoTokenCache(mongodb string, mongodb_database string, mongodb_collection string) (*MongoTokenCache, error) {
	mongoCache, err := NewMongoCache(mongodb, mongodb_database, mongodb_collection, "sessionid")
	if (err != nil) {
		return nil, err
	}
	return &MongoTokenCache{ MongoCache: mongoCache }, nil
}

func (tokenCache *MongoTokenCache) FindTokenDataForSessionId(sessionId string) (*TokenData, error) {
	if (sessionId == "") {
		return nil, fmt.Errorf("Session id cannot be empty")
	}

	result := TokenData{}
	err := tokenCache.MongoCache.FindDataForSessionId("sessionid", sessionId, &result)
	if (err != nil) {
		return nil, err
	}

	return &result, nil
}

func (tokenCache *MongoTokenCache) SaveAuthenticationKeysForSessionId(sessionId string, authenticationToken string, expires_in int64, hash string) (*TokenData, error) {
	if (sessionId != "") {
		existing, _ := tokenCache.FindTokenDataForSessionId(sessionId)
		if (existing != nil) {
			tokenCache.MongoCache.Delete(existing)
		}

               	expiryTime := GetExpiryDate(expires_in)
		tokenData := &TokenData{ Sessionid: sessionId, Authenticationtoken: authenticationToken, Timestamp: expiryTime, Hash: hash  }
		err := tokenCache.MongoCache.Save(tokenData)
		if (err != nil) {
			return nil, err
		}
		return tokenData, nil
	}
	return nil, fmt.Errorf("sessionId cannot be empty")
}

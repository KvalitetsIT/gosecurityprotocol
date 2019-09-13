package securityprotocol

import "fmt"
import "time"

type MongoTokenCache struct {
	MongoCache	*MongoCache
}

func NewMongoTokenCache(mongodb string, mongodb_database string, mongodb_collection string) (*MongoTokenCache, error) {
	mongoCache, err := NewMongoCache(mongodb, mongodb_database, mongodb_collection, "token")
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
	_, err := tokenCache.MongoCache.FindTokenDataForSessionId("sessionid", sessionId, result)
	if (err != nil) {
		return nil, err
	}

	return &result, nil
}

func getExpiryDate(expiresIn int64) time.Time {

        expiryTime := time.Now().Add(time.Duration(expiresIn) * time.Millisecond)
        return expiryTime
}

func (tokenCache *MongoTokenCache) SaveAuthenticationKeysForSessionId(sessionId string, authenticationToken string, expires_in int64, hash string) (*TokenData, error) {
	if (sessionId != "") {
               	expiryTime := getExpiryDate(expires_in)
		tokenData := &TokenData{ Sessionid: sessionId, Authenticationtoken: authenticationToken, Timestamp: expiryTime, Hash: hash  }
		err := tokenCache.MongoCache.Save(tokenData)
		if (err != nil) {
			return nil, err
		}
		return tokenData, nil
	}
	return nil, fmt.Errorf("sessionId cannot be empty")
}

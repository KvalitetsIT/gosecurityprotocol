package securityprotocol

import "fmt"


type MongoTokenCache struct {
	MongoCache	*MongoCache
}

func NewMongoTokenCache(mongodb string, mongodb_database string, mongodb_collection string) (*MongoTokenCache, error) {
	mongoCache, err := NewMongoCache(mongodb, mongodb_database, mongodb_collection, SESSIONID_COLUMN)
	if (err != nil) {
		return nil, err
	}
	return &MongoTokenCache{ MongoCache: mongoCache }, nil
}

func (tokenCache *MongoTokenCache) FindTokenDataForSessionId(sessionId string) (*TokenData, error) {
	if (sessionId == "") {
		return nil, fmt.Errorf("Session id cannot be empty")
	}

	// Query Mongo
	queryTokenData := TokenData{}
	found, err := tokenCache.MongoCache.FindDataForSessionId(SESSIONID_COLUMN, sessionId, &queryTokenData)
	if (err != nil || found == nil) {
		return nil, err
	}

	// Safely cast to TokenData
	result, ok := found.(*TokenData)
	if (ok) {
		return result, nil
	}
	return nil, nil
}

func (tokenCache *MongoTokenCache) SaveAuthenticationKeysForSessionId(sessionId string, authenticationToken string, expires_in int64, hash string) (*TokenData, error) {
	if (sessionId != "") {
		existing, findErr := tokenCache.FindTokenDataForSessionId(sessionId)
		if (findErr != nil) {
			return nil, findErr
		}
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

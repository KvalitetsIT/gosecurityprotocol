package securityprotocol

import "fmt"
import "time"
import primitive "go.mongodb.org/mongo-driver/bson/primitive"

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

func (tokenCache *MongoTokenCache) MaintainCache() error {

        _, err := tokenCache.MongoCache.EnsureIndexes()
        return err
}


func (tokenCache *MongoTokenCache) FindTokenDataForSessionId(sessionId string) (*TokenData, error) {
	if (sessionId == "") {
		return nil, fmt.Errorf("Session id cannot be empty")
	}

	// Query Mongo
	queryTokenData := TokenData{}
	found, err := tokenCache.MongoCache.FindDataForSessionId(SESSIONID_COLUMN, sessionId, &queryTokenData)
	if (err != nil || !found) {
		return nil, err
	}
	return &queryTokenData, nil
}

func (tokenCache *MongoTokenCache) DeleteTokenDataWithId(id primitive.ObjectID) error {
	err := tokenCache.MongoCache.DeleteWithId(id)
	return err
}

func (tokenCache *MongoTokenCache) SaveAuthenticationKeysForSessionId(sessionId string, authenticationToken string, expires_in int64, hash string) (*TokenData, error) {
        expiryTime := GetExpiryDate(expires_in)
	return tokenCache.SaveAuthenticationKeysForSessionIdWithExpiry(sessionId, authenticationToken, expiryTime, hash)
}

func (tokenCache *MongoTokenCache) SaveAuthenticationKeysForSessionIdWithExpiry(sessionId string, authenticationToken string, expiryTime time.Time, hash string) (*TokenData, error) {
        if (sessionId != "") {
                tokenData := &TokenData{ Sessionid: sessionId, Authenticationtoken: authenticationToken, Timestamp: expiryTime, Hash: hash  }
                err := tokenCache.MongoCache.Save(tokenData, tokenData)
                if (err != nil) {
                        return nil, err
                }
                return tokenData, nil
        }
        return nil, fmt.Errorf("sessionId cannot be empty")
}

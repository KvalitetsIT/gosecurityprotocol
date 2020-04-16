package securityprotocol

import "fmt"

type MongoSessionCache struct {
	MongoCache *MongoCache
}

func NewMongoSessionCache(mongodb string, mongodb_database string, mongodb_collection string) (*MongoSessionCache, error) {
	mongoCache, err := NewMongoCache(mongodb, mongodb_database, mongodb_collection, SESSIONID_COLUMN)
	if err != nil {
		return nil, err
	}
	return &MongoSessionCache{MongoCache: mongoCache}, nil
}

func (sessionCache *MongoSessionCache) FindSessionDataForSessionId(sessionId string) (*SessionData, error) {
	if sessionId == "" {
		return nil, fmt.Errorf("Session id cannot be empty")
	}

	// Query Mongo
	sessionData := SessionData{}
	existing, err := sessionCache.MongoCache.FindDataForSessionId(SESSIONID_COLUMN, sessionId, &sessionData)
	if (err != nil) {
		return nil, err
	}
	if (existing) {
		return &sessionData, nil
	}
	return nil, nil
}

func (sessionCache *MongoSessionCache) SaveSessionData(sessionData *SessionData) error {
	sessionId := sessionData.Sessionid
	if sessionId != "" {
		err := sessionCache.MongoCache.Save(sessionData, sessionData)
		if (err != nil) {
			return err
		}
		return nil
	}
	return fmt.Errorf("sessionId cannot be empty")
}

func (sessionCache *MongoSessionCache) DeleteSessionData(sessionId string) error {
	if sessionId != "" {
		sessionData, err := sessionCache.FindSessionDataForSessionId(sessionId)
		if (err != nil) {
			return err
		}
		if (sessionData != nil) {
			err = sessionCache.MongoCache.Delete(sessionData)
			if (err != nil) {
				return err
			}
		}
	}
	return nil
}

package securityprotocol

import "fmt"


type MongoSessionCache struct {

	MongoCache	*MongoCache
}

func NewMongoSessionCache(mongodb string, mongodb_database string, mongodb_collection string) (*MongoSessionCache, error) {
	mongoCache, err := NewMongoCache(mongodb, mongodb_database, mongodb_collection, SESSIONID_COLUMN)
	if (err != nil) {
		return nil, err
	}
	return &MongoSessionCache{ MongoCache: mongoCache }, nil
}

func (sessionCache *MongoSessionCache) FindSessionDataForSessionId(sessionId string) (*SessionData, error) {
	if (sessionId == "") {
		return nil, fmt.Errorf("Session id cannot be empty")
	}

	// Query Mongo
	querySessionData := SessionData{}
	found, err := sessionCache.MongoCache.FindDataForSessionId(SESSIONID_COLUMN, sessionId, &querySessionData)
	if (err != nil) {
		return nil, err
	}
	if (found == nil) {
		return nil, nil
	}

	// Safely cast to SessionData
 	result, ok := found.(*SessionData)
        if (ok) {
                return result, nil
        }
        return nil, nil
}

func (sessionCache *MongoSessionCache) SaveSessionData(sessionData *SessionData) error {
	sessionId := sessionData.Sessionid
	if (sessionId != "") {
		existing, _ := sessionCache.FindSessionDataForSessionId(sessionId)
		if (existing != nil) {
			sessionCache.MongoCache.Delete(existing)
		}
		err := sessionCache.MongoCache.Save(sessionData)
		return err
	}
	return fmt.Errorf("sessionId cannot be empty")
}

package securityprotocol

import "gopkg.in/mgo.v2"
import "gopkg.in/mgo.v2/bson"
import "time"

type TokenCache interface {
	SaveAuthenticationKeysForSessionId(sessionId string, authenticationToken string, expires_in int64, hash string) (*TokenData, error)
	FindTokenDataForSessionId(sessionId string) (*TokenData, error)
}


type MongoTokenCache struct {
	mongoSession    *mgo.Session
	tokenCollection *mgo.Collection
	keyColumn	string
	dbName        	string
	collectionName 	string
}

type TokenData struct {
	ID       		bson.ObjectId `bson:"_id,omitempty"`
	Sessionid      		string
	Authenticationtoken	string
	Timestamp 		time.Time `bson:"timestamp"`
	Hash			string
}

func NewTokenCache(mongodb string, mongodb_database string, mongodb_collection string) (*MongoTokenCache, error) {
	tokenCache, err := newTokenCache(mongodb, mongodb_database, mongodb_collection, "token")
	return tokenCache, err
}

func (tokenCache *MongoTokenCache) ReConnect() {
	tokenCache.mongoSession.Refresh()
	tokenCache.tokenCollection = tokenCache.mongoSession.DB(tokenCache.dbName).C(tokenCache.collectionName)	
}

func newTokenCache(mongodb string, dbName string, collectionName string, keyColumn string) (*MongoTokenCache, error) {

	session, err := mgo.Dial(mongodb)
	if err != nil {
		return nil, err
	}

	session.SetMode(mgo.Monotonic, true)

	// Collection Sessions
	c := session.DB(dbName).C(collectionName)

	// Index
	index := mgo.Index{
		Key:        []string{ keyColumn },
		Unique:     true,
		DropDups:   true,
		Background: true,
		Sparse:     true,
	}
	err = c.EnsureIndex(index)
	if (err != nil) {
		return nil, err
	}

	ttlIndex := mgo.Index {
              	Key:		[]string{ "timestamp" },
  		Unique:		false,
    		Background:	true,
    		ExpireAfter:	time.Second,
	}
  	err = c.EnsureIndex(ttlIndex)
        if (err != nil) {
		return nil, err
        }

        return &MongoTokenCache{ mongoSession: session, tokenCollection: c, keyColumn: keyColumn, dbName: dbName, collectionName: collectionName }, nil
}

func (tokenCache *MongoTokenCache) FindTokenDataForSessionId(sessionId string) (*TokenData, error) {
	if (sessionId == "") {
		return nil, nil
	}
	result := TokenData{}
	err := tokenCache.tokenCollection.Find(bson.M{"sessionid": sessionId}).One(&result)
	if err != nil {
		tokenCache.ReConnect()
		return nil, err
	}
	return &result, nil
}

func GetExpiryDate(expiresIn int64) time.Time {

	expiryTime := time.Now().Add(time.Duration(expiresIn) * time.Millisecond)
	return expiryTime
}

func (tokenCache *MongoTokenCache) SaveAuthenticationKeysForSessionId(sessionId string, authenticationToken string, expires_in int64, hash string) (*TokenData, error) {
	if (sessionId != "") {
               	expiryTime := GetExpiryDate(expires_in)
		tokenData := &TokenData{ Sessionid: sessionId, Authenticationtoken: authenticationToken, Timestamp: expiryTime, Hash: hash  }
		err := tokenCache.tokenCollection.Insert(tokenData)
		if (err != nil) {
			tokenCache.ReConnect()
			return nil, err
		}
		return tokenData, nil
	}
	return nil, nil
}

func (tokenCache *MongoTokenCache) Close() {
	tokenCache.mongoSession.Close()
}

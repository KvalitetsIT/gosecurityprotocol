package securityprotocol

import "gopkg.in/mgo.v2"
import "gopkg.in/mgo.v2/bson"
import "time"

type TokenCache struct {
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

func NewTokenCache(mongodb string, mongodb_database string, mongodb_collection string) (*TokenCache, error) {
	tokenCache, err := newTokenCache(mongodb, mongodb_database, mongodb_collection, "token")
	return tokenCache, err
}

func (tokenCache *TokenCache) ReConnect() {
	tokenCache.mongoSession.Refresh()
	tokenCache.tokenCollection = tokenCache.mongoSession.DB(tokenCache.dbName).C(tokenCache.collectionName)	
}

func newTokenCache(mongodb string, dbName string, collectionName string, keyColumn string) (*TokenCache, error) {

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

        return &TokenCache{ mongoSession: session, tokenCollection: c, keyColumn: keyColumn, dbName: dbName, collectionName: collectionName }, nil
}

func (tokenCache *TokenCache) FindTokenDataForSessionId(sessionId string) (*TokenData, error) {
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

func (tokenCache *TokenCache) SaveAuthenticationKeysForSessionId(sessionId string, authenticationToken string, expires_in int64, hash string) error {
	if (sessionId != "") {
               	expiryTime := GetExpiryDate(expires_in)
		err := tokenCache.tokenCollection.Insert(&TokenData{Sessionid: sessionId, Authenticationtoken: authenticationToken, Timestamp: expiryTime, Hash: hash  })
		if err != nil {
			tokenCache.ReConnect()
			return err
		}
	}
	return nil
}

func (tokenCache *TokenCache) Close() {
	tokenCache.mongoSession.Close()
}

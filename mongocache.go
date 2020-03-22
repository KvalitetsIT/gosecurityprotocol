package securityprotocol

import "gopkg.in/mgo.v2"
import "gopkg.in/mgo.v2/bson"
import "time"

type MongoCache struct {
	mongoSession    *mgo.Session
	collection 	*mgo.Collection
	keyColumn	string
	dbName        	string
	collectionName 	string
}

func NewMongoCache(mongodb string, mongodb_database string, mongodb_collection string, keyColumn string) (*MongoCache, error) {
	tokenCache, err := newCache(mongodb, mongodb_database, mongodb_collection, keyColumn)
	return tokenCache, err
}

func (mongoCache *MongoCache) ReConnect() {
	mongoCache.mongoSession.Refresh()
	mongoCache.collection = mongoCache.mongoSession.DB(mongoCache.dbName).C(mongoCache.collectionName)
}

func newCache(mongodb string, dbName string, collectionName string, keyColumn string) (*MongoCache, error) {

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

        return &MongoCache{ mongoSession: session, collection: c, keyColumn: keyColumn, dbName: dbName, collectionName: collectionName }, nil
}

func (mongoCache *MongoCache) FindDataForSessionId(sessionKey string, sessionId string, object interface{}) (interface{}, error) {
	if (sessionId == "") {
		return nil, nil
	}

	query := mongoCache.collection.Find(bson.M{sessionKey: sessionId})
	count, err := query.Count()
	if (err != nil) {
		return nil, err
	}

	if (count == 1) {
		err = query.One(object)
		if (err != nil) {
                	mongoCache.ReConnect()
                	return nil, err
        	}
		return object, nil
	}
	return nil, nil
}

func (mongoCache *MongoCache) Delete(object interface{}) error {
        err := mongoCache.collection.Remove(object)
        if (err != nil) {
                mongoCache.ReConnect()
                return err
        }
        return nil
}

func (mongoCache *MongoCache) Save(object interface{}) error {
	err := mongoCache.collection.Insert(object)
	if (err != nil) {
		mongoCache.ReConnect()
		return err
	}
	return nil
}

func (mongoCache *MongoCache) Close() {
	mongoCache.mongoSession.Close()
}

func GetExpiryDate(expiresIn int64) time.Time {

        expiryTime := time.Now().Add(time.Duration(expiresIn) * time.Second)
        return expiryTime
}


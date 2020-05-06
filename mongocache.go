package securityprotocol


import (
    "go.mongodb.org/mongo-driver/mongo"
    options "go.mongodb.org/mongo-driver/mongo/options"
    bson "go.mongodb.org/mongo-driver/bson"
    primitive "go.mongodb.org/mongo-driver/bson/primitive"
    "context"
    "time"
    "fmt"
)

type MongoCache struct {
     	mongoClient     *mongo.Client
	keyColumn	string
	dbName        	string
	collectionName 	string
}

func NewMongoCache(mongodb string, mongodb_database string, mongodb_collection string, keyColumn string) (*MongoCache, error) {
	tokenCache, err := newCache(mongodb, mongodb_database, mongodb_collection, keyColumn)
	return tokenCache, err
}

/*func (mongoCache *MongoCache) ReConnect() {
	mongoCache.mongoSession.Refresh()
	mongoCache.collection = mongoCache.mongoSession.DB(mongoCache.dbName).C(mongoCache.collectionName)
}
*/
func (mongoCache *MongoCache) getCollection() *mongo.Collection {

	collection := mongoCache.mongoClient.Database(mongoCache.dbName).Collection(mongoCache.collectionName)
	return collection
}

func getContext() context.Context {

	ctx, _ := context.WithTimeout(context.Background(), 5*time.Second)
	return ctx
}

func newCache(mongodb string, dbName string, collectionName string, keyColumn string) (*MongoCache, error) {

	mongoClient, err := mongo.NewClient(options.Client().ApplyURI(fmt.Sprintf("mongodb://%s", mongodb)))
	if (err != nil) {
		return nil, err
	}
	ctx := getContext()
	err = mongoClient.Connect(ctx)
	if (err != nil) {
		return nil, err
	}

/*	// Index
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
*/
        mongoCache := MongoCache{ mongoClient: mongoClient, keyColumn: keyColumn, dbName: dbName, collectionName: collectionName }
	mongoCache.EnsureIndexes()

	return &mongoCache, nil
}

func (mongoCache *MongoCache) EnsureIndexes() ([]string, error) {

	ttlName := "idx_session_ttl"
	ukName := "idx_session_uk"

        t := true
	ctx := getContext()
        ukIndexOptions := options.IndexOptions {
		Name: &ukName,
                Unique: &t,
        }
        uniqueKeyIndexModel := mongo.IndexModel{
                Keys: bson.M{ mongoCache.keyColumn: 1 },
                Options: &ukIndexOptions,
        }
        _, err := mongoCache.getCollection().Indexes().CreateOne(ctx, uniqueKeyIndexModel)
        if (err != nil) {
                return nil, err
        }

        var expiryAfterSeconds int32
        expiryAfterSeconds = 0
        ttlIndexOptions := options.IndexOptions {
		Name: &ttlName,
                ExpireAfterSeconds: &expiryAfterSeconds,
        }
        ttlIndexModel := mongo.IndexModel{
                Keys: bson.M{ "timestamp": 1 },
                Options: &ttlIndexOptions,
        }

	opts := options.CreateIndexes().SetMaxTime(2 * time.Second)
	return mongoCache.getCollection().Indexes().CreateMany(ctx, []mongo.IndexModel{ uniqueKeyIndexModel, ttlIndexModel }, opts)
}

func (mongoCache *MongoCache) FindDataForSessionId(sessionKey string, sessionId string, result interface{}) (bool, error) {
	if (sessionId == "") {
		return false, nil
	}

	collection := mongoCache.getCollection()
	ctx := getContext()

	cur, err := collection.Find(ctx, bson.M{sessionKey: sessionId})
	if (err != nil) {
		return false, err
	}
	defer cur.Close(ctx)
	for cur.Next(ctx) {
   		err := cur.Decode(result)
   		if (err != nil) {
			return false, err
		}
		return true, nil
	}
	return false, nil
}

func (mongoCache *MongoCache) DeleteWithId(id primitive.ObjectID) error {

	collection := mongoCache.getCollection()
        ctx := getContext()

	_, err := collection.DeleteOne(ctx, bson.M{"_id": id})

        return err
}


func (mongoCache *MongoCache) Delete(object Identifiable) error {
	collection := mongoCache.getCollection()
        ctx := getContext()

        _, err := collection.DeleteOne(ctx, object)
        return err
}

func (mongoCache *MongoCache) Save(object interface{}, filterId Identifiable) error {

	collection := mongoCache.getCollection()
        ctx := getContext()

	filter := bson.M{ mongoCache.keyColumn : filterId.GetKey() }

	upsert := true
	after := options.After
	opt := options.FindOneAndReplaceOptions{
		ReturnDocument: &after,
		Upsert:         &upsert,
	}

	result := collection.FindOneAndReplace(ctx, filter, object, &opt)
	if result.Err() != nil {
		return result.Err()
	}

	decodeErr := result.Decode(object)
	if (decodeErr != nil) {
		return decodeErr
	}

	return nil
}

func GetExpiryDate(expiresIn int64) time.Time {

        expiryTime := time.Now().Add(time.Duration(expiresIn) * time.Second)
        return expiryTime
}

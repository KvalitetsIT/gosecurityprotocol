package securityprotocol

import "time"
import "gopkg.in/mgo.v2/bson"

type TokenCache interface {
	SaveAuthenticationKeysForSessionId(sessionId string, authenticationToken string, expires_in int64, hash string) (*TokenData, error)
	FindTokenDataForSessionId(sessionId string) (*TokenData, error)
}


type TokenData struct {
	ID       		bson.ObjectId `bson:"_id,omitempty"`
	Sessionid      		string
	Authenticationtoken	string
	Timestamp 		time.Time `bson:"timestamp"`
	Hash			string
}

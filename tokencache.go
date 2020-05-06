package securityprotocol

import "time"
import  "crypto/md5"
import  "io"
import  "encoding/base64"
import primitive "go.mongodb.org/mongo-driver/bson/primitive"

const SESSIONID_COLUMN = "sessionid"

type TokenCache interface {
	SaveAuthenticationKeysForSessionId(sessionId string, authenticationToken string, expires_in int64, hash string) (*TokenData, error)
	SaveAuthenticationKeysForSessionIdWithExpiry(sessionId string, authenticationToken string, expiryTime time.Time, hash string) (*TokenData, error)
	FindTokenDataForSessionId(sessionId string) (*TokenData, error)
	DeleteTokenDataWithId(id primitive.ObjectID) error
}


type TokenData struct {
	ID       		*primitive.ObjectID `bson:"_id,omitempty"`
	Sessionid      		string `bson:"sessionid"`
	Authenticationtoken	string
	Timestamp 		time.Time `bson:"timestamp"`
	Hash			string
}

func (data TokenData) GetID() *primitive.ObjectID {
        return data.ID
}

func (data TokenData) GetKey() string {
        return data.Sessionid
}


func (data *TokenData) CalculateHash() string {

	s := data.Sessionid
	s = s + data.Authenticationtoken
	s = s + data.Timestamp.Format(time.UnixDate)

	h := md5.New()
	io.WriteString(h, s)
	hash := base64.StdEncoding.EncodeToString(h.Sum(nil))

	return hash
}


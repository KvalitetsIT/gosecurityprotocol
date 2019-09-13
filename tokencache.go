package securityprotocol

import "time"
import "gopkg.in/mgo.v2/bson"
import  "crypto/md5"
import  "io"
import  "encoding/base64"


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

func (data *TokenData) CalculateHash() string {

	s := data.Sessionid
	s = s + data.Authenticationtoken
	s = s + data.Timestamp.Format(time.UnixDate)

	h := md5.New()
	io.WriteString(h, s)
	hash := base64.StdEncoding.EncodeToString(h.Sum(nil))

	return hash
}


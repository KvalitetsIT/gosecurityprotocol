package securityprotocol

import (
	uuid "github.com/google/uuid"
	"gotest.tools/assert"
	"testing"
)

func TestMongoSessionCache(t *testing.T) {

	// Given
	mongoSessionCache, createErr := NewMongoSessionCache("mongo", "testsession", "session")
	assert.NilError(t, createErr)

	sessionId := uuid.New().String()
	hash1 := "hash1"
	hash2 := "hash2"

	sessionDataFirst := SessionData{Sessionid: sessionId, Hash: hash1}
	sessionDataSecond := SessionData{Sessionid: sessionId, Hash: hash2}

	// When
	saveFirstErr := mongoSessionCache.SaveSessionData(&sessionDataFirst)
	assert.NilError(t, saveFirstErr)
	assert.Assert(t, sessionDataFirst.GetID() != nil)

	saveSecondErr := mongoSessionCache.SaveSessionData(&sessionDataSecond)
	assert.NilError(t, saveSecondErr)
	assert.Assert(t, sessionDataSecond.GetID() != nil)

	assert.Equal(t, sessionDataFirst.GetID().Hex(), sessionDataSecond.GetID().Hex())

	sessionDataGet, getErr := mongoSessionCache.FindSessionDataForSessionId(sessionId)
	assert.NilError(t, getErr)
	assert.Assert(t, sessionDataGet != nil)
	assert.Equal(t, sessionDataFirst.GetID().Hex(), sessionDataGet.GetID().Hex())

	mongoSessionCache.DeleteSessionData(sessionId)
	sessionDataGetDeleted, err := mongoSessionCache.FindSessionDataForSessionId(sessionId)
	assert.NilError(t, err)
	assert.Assert(t, sessionDataGetDeleted == nil)

	// Then
	assert.Equal(t, sessionId, sessionDataGet.Sessionid)
	assert.Equal(t, hash2, sessionDataGet.Hash)
}

/* TODO

func TestMongoSessionCacheSaveAndGet(t *testing.T) {

        // Given
	mongoSessionCache, createErr := NewMongoSessionCache("mongo", "testsession", "session")

        // When
        tokenDataFirst, saveErr := mongoTokenCache.SaveAuthenticationKeysForSessionId(sessionId, firstToken, 2000, firstHash)
        tokenDataGet, getErr := mongoTokenCache.FindTokenDataForSessionId(sessionId)

        // Then
        assert.NilError(t, saveErr)
        assert.NilError(t, getErr)

        assert.Equal(t, sessionId, tokenDataFirst.Sessionid)
        assert.Equal(t, firstToken, tokenDataFirst.Authenticationtoken)
        assert.Equal(t, firstHash, tokenDataFirst.Hash)

        assert.Equal(t, sessionId, tokenDataGet.Sessionid)
        assert.Equal(t, firstToken, tokenDataGet.Authenticationtoken)
        assert.Equal(t, firstHash, tokenDataGet.Hash)
}


func IgnoreTestFindNonExistingReturnsNil (t *testing.T) {

        // Given
        mongoTokenCache, _ := NewMongoTokenCache("mongo", "testdb", "testcoll")
        sessionId := fmt.Sprintf("nonexisting-%s", uuid.New().String())

	// When
	tokenDataGet, getErr := mongoTokenCache.FindTokenDataForSessionId(sessionId)

	// Then
	assert.NilError(t, getErr)
	assert.Assert(t, tokenDataGet == nil)
}
*/

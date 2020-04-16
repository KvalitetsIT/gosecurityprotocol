package securityprotocol

import (
        "testing"
        "gotest.tools/assert"
	"fmt"
	uuid "github.com/google/uuid"
)

func TestMongoTokenCache(t *testing.T) {

	// Given
	mongoTokenCache, createErr := NewMongoTokenCache("mongo", "testdb", "testcoll")
	assert.NilError(t, createErr)
	sessionId := fmt.Sprintf("sessionid-%s", uuid.New().String())
	testToken := "test-token"
	firstToken:= "first-token"
	testHash  := fmt.Sprintf("hash-xyz-%s", uuid.New().String())

        // When
        tokenDataFirst, saveErr := mongoTokenCache.SaveAuthenticationKeysForSessionId(sessionId, firstToken, 2000, "first hash")
	assert.NilError(t, saveErr)
	tokenDataSaved, saveErr := mongoTokenCache.SaveAuthenticationKeysForSessionId(sessionId, testToken, 1000, testHash)
	assert.NilError(t, saveErr)
	fmt.Println(fmt.Sprintf("leder efter: %s", sessionId))
	tokenDataGet, getErr := mongoTokenCache.FindTokenDataForSessionId(sessionId)

	// Then
	assert.NilError(t, getErr)

	assert.Equal(t, firstToken, tokenDataFirst.Authenticationtoken)

	assert.Equal(t, sessionId, tokenDataSaved.Sessionid)
	assert.Equal(t, testToken, tokenDataSaved.Authenticationtoken)
	assert.Equal(t, testHash, tokenDataSaved.Hash)

        assert.Equal(t, sessionId, tokenDataGet.Sessionid)
        assert.Equal(t, testToken, tokenDataGet.Authenticationtoken)
        assert.Equal(t, testHash, tokenDataGet.Hash)
}


func TestMongoTokenCacheSaveAndGet(t *testing.T) {

        // Given
        mongoTokenCache, _ := NewMongoTokenCache("mongo", "testdb", "testcoll")
        sessionId := fmt.Sprintf("sessionid-%s", uuid.New().String())
        firstToken:= "first-token"
        firstHash  := fmt.Sprintf("hash-xyz-%s", uuid.New().String())

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

func TestMongoTokenCacheSaveAndDelete(t *testing.T) {
	fmt.Println("DUM TEST")

        // Given
        mongoTokenCache, _ := NewMongoTokenCache("mongo", "testdb", "testcoll")
        sessionId := fmt.Sprintf("sessionid-%s123", uuid.New().String())
        firstToken:= "first-token"
        firstHash  := fmt.Sprintf("hash-xyz-%s", uuid.New().String())

        // When
        tokenDataFirst, saveErr := mongoTokenCache.SaveAuthenticationKeysForSessionId(sessionId, firstToken, 2000, firstHash)

	deleteErr := mongoTokenCache.DeleteTokenDataWithId(*tokenDataFirst.ID)

        tokenDataGetAfterDelete, getErr := mongoTokenCache.FindTokenDataForSessionId(sessionId)

        // Then
        assert.NilError(t, saveErr)
	assert.NilError(t, deleteErr)
        assert.NilError(t, getErr)

        assert.Equal(t, sessionId, tokenDataFirst.Sessionid)
        assert.Equal(t, firstToken, tokenDataFirst.Authenticationtoken)
        assert.Equal(t, firstHash, tokenDataFirst.Hash)

        assert.Assert(t, tokenDataGetAfterDelete == nil)
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

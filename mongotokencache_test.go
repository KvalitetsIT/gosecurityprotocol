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
	sessionId := fmt.Sprintf("sessionid-%s", uuid.New().String())
	testToken := "test-token"
	testHash  := fmt.Sprintf("hash-xyz-%s", uuid.New().String())

        // When
	tokenDataSaved, saveErr := mongoTokenCache.SaveAuthenticationKeysForSessionId(sessionId, testToken, 1000, testHash)
	tokenDataGet, getErr := mongoTokenCache.FindTokenDataForSessionId(sessionId)

	// Then
	assert.NilError(t, createErr)
	assert.NilError(t, saveErr)
	assert.NilError(t, getErr)

	assert.Equal(t, sessionId, tokenDataSaved.Sessionid)
	assert.Equal(t, testToken, tokenDataSaved.Authenticationtoken)
	assert.Equal(t, testHash, tokenDataSaved.Hash)

        assert.Equal(t, sessionId, tokenDataGet.Sessionid)
        assert.Equal(t, testToken, tokenDataGet.Authenticationtoken)
        assert.Equal(t, testHash, tokenDataGet.Hash)
}


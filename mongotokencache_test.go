package securityprotocol

import (
        "testing"
	"net/http"
	"fmt"
	"time"
        "gotest.tools/assert"
)

func TestMongoTokenCache(t *testing.T) {

	// Given
	mongoTokenCache, createErr := NewMongoTokenCache("mongo", "testdb", "testcoll")
	sessionId := "session-1234"

        // When
	tokenData, saveErr := mongoTokenCache.SaveAuthenticationKeysForSessionId(sessionId, "test-token", 1000, "hash-xyz")

	// Then
	assert.NilError(t, createErr)
	assert.NilError(t, saveErr)
}


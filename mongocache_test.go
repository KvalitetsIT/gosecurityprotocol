package securityprotocol

import (
        "testing"
        "gotest.tools/assert"
)

func TestMongoCache(t *testing.T) {

	mongoCache, err := NewMongoCache("mongo", "testmc", "xyz", "keycol")
	assert.NilError(t, err)
	assert.Assert(t, mongoCache != nil)

	createdIndexNames, err := mongoCache.EnsureIndexes()
	assert.NilError(t, err)
        assert.Assert(t, createdIndexNames != nil)
	assert.Equal(t, 0, len(createdIndexNames))

}

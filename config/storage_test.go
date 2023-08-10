package config

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestStorageLoad(t *testing.T) {
	var s Storage

	assert.ErrorIs(t, s.load("storage"), ErrNoDbURI)

	require.NoError(t, os.Setenv("STORAGE_URI", "mongo://localhost:28017"))
	assert.ErrorIs(t, s.load("storage"), ErrNoDbName)

	require.NoError(t, os.Setenv("STORAGE_DB_NAME", "test-db"))
	require.NoError(t, s.load("storage"))

	assert.Equal(t, "mongo://localhost:28017", s.URI)
	assert.Equal(t, "test-db", s.DbName)
}

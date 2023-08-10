package config

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestFeedLoad(t *testing.T) {
	t.Parallel()

	var f Feed
	assert.ErrorIs(t, f.load("feed"), ErrNoApiURL)

	require.NoError(t, os.Setenv("FEED_API_URL", "bad URL format"))
	assert.ErrorIs(t, f.load("feed"), ErrInvalidUrlFormat)

	require.NoError(t, os.Setenv("FEED_API_URL", "google.com"))
	require.NoError(t, f.load("feed"))

	assert.Equal(t, "google.com", f.ApiURL.String())
}

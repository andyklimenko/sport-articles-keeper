package config

import (
	"os"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestPollerLoad(t *testing.T) {
	t.Parallel()

	var p Poller
	assert.ErrorIs(t, p.load("poll"), ErrZeroInterval)

	require.NoError(t, os.Setenv("POLL_INTERVAL", "invalid duration"))
	assert.ErrorIs(t, p.load("poll"), ErrZeroInterval)

	require.NoError(t, os.Setenv("POLL_INTERVAL", "5m"))
	require.NoError(t, p.load("poll"))
	assert.Equal(t, 5*time.Minute, p.Interval)
}

package config

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestServerLoad(t *testing.T) {
	var s Server
	assert.ErrorIs(t, s.load("server"), ErrNoServerAddr)

	require.NoError(t, os.Setenv("SERVER_ADDR", "0.0.0.0:8080"))
	require.NoError(t, s.load("server"))

	assert.Equal(t, "0.0.0.0:8080", s.Address)
}

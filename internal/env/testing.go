package env

import (
	"os"
	"testing"

	"github.com/stretchr/testify/require"
)

// NewTestConfig is a helper for setting up a Config object and error checking
func NewTestConfig(t *testing.T) *Config {
	conf, err := NewConfig(os.Stdout, os.Stderr)
	require.NoError(t, err)

	return conf
}

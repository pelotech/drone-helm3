package env

import (
	"os"
	"testing"

	"github.com/stretchr/testify/require"
)

func NewTestConfig(t *testing.T) *Config {
	conf, err := NewConfig(os.Stdout, os.Stderr)
	require.NoError(t, err)

	return conf
}

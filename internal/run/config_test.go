package run

import (
	"github.com/pelotech/drone-helm3/internal/env"
	"github.com/stretchr/testify/suite"
	"strings"
	"testing"
)

type ConfigTestSuite struct {
	suite.Suite
}

func TestConfigTestSuite(t *testing.T) {
	suite.Run(t, new(ConfigTestSuite))
}

func (suite *ConfigTestSuite) TestNewConfig() {
	stdout := &strings.Builder{}
	stderr := &strings.Builder{}
	envCfg := env.Config{
		Namespace: "private",
		Debug:     true,
		Stdout:    stdout,
		Stderr:    stderr,
	}
	cfg := newConfig(envCfg)
	suite.Require().NotNil(cfg)
	suite.Equal(&config{
		namespace: "private",
		debug:     true,
		stdout:    stdout,
		stderr:    stderr,
	}, cfg)
}

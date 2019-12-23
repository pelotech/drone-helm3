package helm

import (
	"github.com/stretchr/testify/suite"
	"os"
	"testing"
)

type ConfigTestSuite struct {
	suite.Suite
	// These tests need to mutate the environment, so the suite.setenv and .unsetenv functions store the original contents of the
	// relevant variable in this map. Its use of *string is so they can distinguish between "not set" and "set to empty string"
	envBackup map[string]*string
}

func TestConfigTestSuite(t *testing.T) {
	suite.Run(t, new(ConfigTestSuite))
}

func (suite *ConfigTestSuite) TestPopulateWithPluginPrefix() {
	suite.unsetenv("HELM_COMMAND")
	suite.unsetenv("UPDATE_DEPENDENCIES")
	suite.unsetenv("DEBUG")

	suite.setenv("PLUGIN_HELM_COMMAND", "execute order 66")
	suite.setenv("PLUGIN_UPDATE_DEPENDENCIES", "true")
	suite.setenv("PLUGIN_DEBUG", "true")

	cfg := Config{}
	cfg.Populate()

	suite.Equal("execute order 66", cfg.Command)
	suite.True(cfg.UpdateDependencies)
	suite.True(cfg.Debug)
}

func (suite *ConfigTestSuite) TestPopulateWithNoPrefix() {
	suite.unsetenv("PLUGIN_HELM_COMMAND")
	suite.unsetenv("PLUGIN_UPDATE_DEPENDENCIES")
	suite.unsetenv("PLUGIN_DEBUG")

	suite.setenv("HELM_COMMAND", "execute order 66")
	suite.setenv("UPDATE_DEPENDENCIES", "true")
	suite.setenv("DEBUG", "true")

	cfg := Config{}
	cfg.Populate()

	suite.Equal("execute order 66", cfg.Command)
	suite.True(cfg.UpdateDependencies)
	suite.True(cfg.Debug)
}

func (suite *ConfigTestSuite) TestPopulateWithConflictingVariables() {
	suite.setenv("PLUGIN_HELM_COMMAND", "execute order 66")
	suite.setenv("HELM_COMMAND", "defend the jedi")

	cfg := Config{}
	cfg.Populate()

	suite.Equal("defend the jedi", cfg.Command)
}

func (suite *ConfigTestSuite) setenv(key, val string) {
	orig, ok := os.LookupEnv(key)
	if ok {
		suite.envBackup[key] = &orig
	} else {
		suite.envBackup[key] = nil
	}
	os.Setenv(key, val)
}

func (suite *ConfigTestSuite) unsetenv(key string) {
	orig, ok := os.LookupEnv(key)
	if ok {
		suite.envBackup[key] = &orig
	} else {
		suite.envBackup[key] = nil
	}
	os.Unsetenv(key)
}

func (suite *ConfigTestSuite) BeforeTest(_, _ string) {
	suite.envBackup = make(map[string]*string)
}

func (suite *ConfigTestSuite) AfterTest(_, _ string) {
	for key, val := range suite.envBackup {
		if val == nil {
			os.Unsetenv(key)
		} else {
			os.Setenv(key, *val)
		}
	}
}

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
	suite.unsetenv("PLUGIN_PREFIX")
	suite.unsetenv("HELM_COMMAND")
	suite.unsetenv("UPDATE_DEPENDENCIES")
	suite.unsetenv("DEBUG")

	suite.setenv("PLUGIN_HELM_COMMAND", "execute order 66")
	suite.setenv("PLUGIN_UPDATE_DEPENDENCIES", "true")
	suite.setenv("PLUGIN_DEBUG", "true")

	cfg := Config{}
	suite.Require().NoError(cfg.Populate())

	suite.Equal("execute order 66", cfg.Command)
	suite.True(cfg.UpdateDependencies)
	suite.True(cfg.Debug)
}

func (suite *ConfigTestSuite) TestPopulateWithNoPrefix() {
	suite.unsetenv("PLUGIN_PREFIX")
	suite.unsetenv("PLUGIN_HELM_COMMAND")
	suite.unsetenv("PLUGIN_UPDATE_DEPENDENCIES")
	suite.unsetenv("PLUGIN_DEBUG")

	suite.setenv("HELM_COMMAND", "execute order 66")
	suite.setenv("UPDATE_DEPENDENCIES", "true")
	suite.setenv("DEBUG", "true")

	cfg := Config{}
	suite.Require().NoError(cfg.Populate())

	suite.Equal("execute order 66", cfg.Command)
	suite.True(cfg.UpdateDependencies)
	suite.True(cfg.Debug)
}

func (suite *ConfigTestSuite) TestPopulateWithConfigurablePrefix() {
	suite.unsetenv("API_SERVER")
	suite.unsetenv("PLUGIN_API_SERVER")

	suite.setenv("PLUGIN_PREFIX", "prix_fixe")
	suite.setenv("PRIX_FIXE_API_SERVER", "your waiter this evening")

	cfg := Config{}
	suite.Require().NoError(cfg.Populate())

	suite.Equal("prix_fixe", cfg.Prefix)
	suite.Equal("your waiter this evening", cfg.APIServer)
}

func (suite *ConfigTestSuite) TestPrefixSettingDoesNotAffectPluginPrefix() {
	suite.setenv("PLUGIN_PREFIX", "IXFREP")
	suite.setenv("PLUGIN_HELM_COMMAND", "wake me up")
	suite.setenv("IXFREP_PLUGIN_HELM_COMMAND", "send me to sleep inside")

	cfg := Config{}
	suite.Require().NoError(cfg.Populate())

	suite.Equal("wake me up", cfg.Command)
}

func (suite *ConfigTestSuite) TestPrefixSettingMustHavePluginPrefix() {
	suite.unsetenv("PLUGIN_PREFIX")
	suite.setenv("PREFIX", "refpix")
	suite.setenv("HELM_COMMAND", "gimme more")
	suite.setenv("REFPIX_HELM_COMMAND", "gimme less")

	cfg := Config{}
	suite.Require().NoError(cfg.Populate())

	suite.Equal("gimme more", cfg.Command)
}

func (suite *ConfigTestSuite) TestPopulateWithConflictingVariables() {
	suite.setenv("PLUGIN_HELM_COMMAND", "execute order 66")
	suite.setenv("HELM_COMMAND", "defend the jedi") // values from the `environment` block override those from `settings`

	suite.setenv("PLUGIN_PREFIX", "prod")
	suite.setenv("TIMEOUT", "5m0s")
	suite.setenv("PROD_TIMEOUT", "2m30s") // values from prefixed env vars override those from non-prefixed ones

	cfg := Config{}
	suite.Require().NoError(cfg.Populate())

	suite.Equal("defend the jedi", cfg.Command)
	suite.Equal("2m30s", cfg.Timeout)
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

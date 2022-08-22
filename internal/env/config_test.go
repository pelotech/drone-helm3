package env

import (
	"fmt"
	"os"
	"strings"
	"testing"

	"github.com/stretchr/testify/suite"
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

func (suite *ConfigTestSuite) TestNewConfigWithPluginPrefix() {
	suite.unsetenv("MODE")
	suite.unsetenv("UPDATE_DEPENDENCIES")
	suite.unsetenv("DEBUG")

	suite.setenv("PLUGIN_MODE", "iambic")
	suite.setenv("PLUGIN_UPDATE_DEPENDENCIES", "true")
	suite.setenv("PLUGIN_DEBUG", "true")

	cfg, err := NewConfig(&strings.Builder{}, &strings.Builder{})
	suite.Require().NoError(err)

	suite.Equal("iambic", cfg.Command)
	suite.True(cfg.UpdateDependencies)
	suite.True(cfg.Debug)
}

func (suite *ConfigTestSuite) TestNewConfigWithNoPrefix() {
	suite.unsetenv("PLUGIN_MODE")
	suite.unsetenv("PLUGIN_UPDATE_DEPENDENCIES")
	suite.unsetenv("PLUGIN_DEBUG")

	suite.setenv("MODE", "iambic")
	suite.setenv("UPDATE_DEPENDENCIES", "true")
	suite.setenv("DEBUG", "true")

	cfg, err := NewConfig(&strings.Builder{}, &strings.Builder{})
	suite.Require().NoError(err)

	suite.Equal("iambic", cfg.Command)
	suite.True(cfg.UpdateDependencies)
	suite.True(cfg.Debug)
}

func (suite *ConfigTestSuite) TestNewConfigWithConflictingVariables() {
	suite.setenv("PLUGIN_MODE", "iambic")
	suite.setenv("MODE", "haiku") // values from the `environment` block override those from `settings`

	cfg, err := NewConfig(&strings.Builder{}, &strings.Builder{})
	suite.Require().NoError(err)

	suite.Equal("haiku", cfg.Command)
}

func (suite *ConfigTestSuite) TestNewConfigInfersNumbersAreSeconds() {
	suite.setenv("PLUGIN_TIMEOUT", "42")
	cfg, err := NewConfig(&strings.Builder{}, &strings.Builder{})
	suite.Require().NoError(err)
	suite.Equal("42s", cfg.Timeout)
}

func (suite *ConfigTestSuite) TestNewConfigWithAliases() {
	for _, varname := range []string{
		"MODE",
		"ADD_REPOS",
		"KUBE_API_SERVER",
		"KUBE_SERVICE_ACCOUNT",
		"WAIT_FOR_UPGRADE",
		"FORCE_UPGRADE",
		"KUBE_TOKEN",
		"KUBE_CERTIFICATE",
	} {
		suite.unsetenv(varname)
		suite.unsetenv("PLUGIN_" + varname)
	}
	suite.setenv("PLUGIN_HELM_COMMAND", "beware the jabberwock")
	suite.setenv("PLUGIN_HELM_REPOS", "chortle=http://calloo.callay/frabjous/day")
	suite.setenv("PLUGIN_API_SERVER", "http://tumtum.tree")
	suite.setenv("PLUGIN_SERVICE_ACCOUNT", "tulgey")
	suite.setenv("PLUGIN_WAIT", "true")
	suite.setenv("PLUGIN_FORCE", "true")
	suite.setenv("PLUGIN_KUBERNETES_TOKEN", "Y29tZSB0byBteSBhcm1z")
	suite.setenv("PLUGIN_KUBERNETES_CERTIFICATE", "d2l0aCBpdHMgaGVhZA==")

	cfg, err := NewConfig(&strings.Builder{}, &strings.Builder{})
	suite.Require().NoError(err)
	suite.Equal("beware the jabberwock", cfg.Command)
	suite.Equal([]string{"chortle=http://calloo.callay/frabjous/day"}, cfg.AddRepos)
	suite.Equal("http://tumtum.tree", cfg.APIServer)
	suite.Equal("tulgey", cfg.ServiceAccount)
	suite.True(cfg.Wait, "Wait should be aliased")
	suite.True(cfg.Force, "Force should be aliased")
	suite.Equal("Y29tZSB0byBteSBhcm1z", cfg.KubeToken, "KubeToken should be aliased")
	suite.Equal("d2l0aCBpdHMgaGVhZA==", cfg.Certificate, "Certificate should be aliased")
}

func (suite *ConfigTestSuite) TestAliasedSettingWithoutPluginPrefix() {
	suite.unsetenv("FORCE_UPGRADE")
	suite.unsetenv("PLUGIN_FORCE_UPGRADE")
	suite.unsetenv("PLUGIN_FORCE")
	suite.setenv("FORCE", "true")

	cfg, err := NewConfig(&strings.Builder{}, &strings.Builder{})
	suite.Require().NoError(err)
	suite.True(cfg.Force)
}

func (suite *ConfigTestSuite) TestNewConfigWithAliasConflicts() {
	suite.unsetenv("FORCE_UPGRADE")
	suite.setenv("PLUGIN_FORCE", "true")
	suite.setenv("PLUGIN_FORCE_UPGRADE", "false") // should override even when set to the zero value

	cfg, err := NewConfig(&strings.Builder{}, &strings.Builder{})
	suite.NoError(err)
	suite.False(cfg.Force, "official names should override alias names")
}

func (suite *ConfigTestSuite) TestNewConfigSetsWriters() {
	stdout := &strings.Builder{}
	stderr := &strings.Builder{}
	cfg, err := NewConfig(stdout, stderr)
	suite.Require().NoError(err)

	suite.Equal(stdout, cfg.Stdout)
	suite.Equal(stderr, cfg.Stderr)
}

func (suite *ConfigTestSuite) TestDeprecatedSettingWarnings() {
	for _, varname := range deprecatedVars {
		suite.setenv(varname, "deprecoat") // environment-block entries should cause warnings
	}

	suite.unsetenv("PURGE")
	suite.setenv("PLUGIN_PURGE", "true") // settings-block entries should cause warnings
	suite.setenv("UPGRADE", "")          // entries should cause warnings even when set to empty string

	stderr := &strings.Builder{}
	_, err := NewConfig(&strings.Builder{}, stderr)
	suite.NoError(err)

	for _, varname := range deprecatedVars {
		suite.Contains(stderr.String(), fmt.Sprintf("Warning: ignoring deprecated '%s' setting\n", strings.ToLower(varname)))
	}
}

func (suite *ConfigTestSuite) TestLogDebug() {
	suite.setenv("DEBUG", "true")
	suite.setenv("MODE", "upgrade")

	stderr := strings.Builder{}
	stdout := strings.Builder{}
	_, err := NewConfig(&stdout, &stderr)
	suite.Require().NoError(err)

	suite.Equal("", stdout.String())

	suite.Regexp(`^Generated config: \{Command:upgrade.*\}`, stderr.String())
}

func (suite *ConfigTestSuite) TestLogDebugCensorsKubeToken() {
	stderr := &strings.Builder{}
	kubeToken := "I'm shy! Don't put me in your build logs!"
	cfg := Config{
		Debug:     true,
		KubeToken: kubeToken,
		Stderr:    stderr,
	}

	cfg.logDebug()

	suite.Contains(stderr.String(), "KubeToken:(redacted)")
	suite.Equal(kubeToken, cfg.KubeToken) // The actual config value should be left unchanged
}

func (suite *ConfigTestSuite) TestNewConfigWithValuesSecrets() {
	suite.unsetenv("VALUES")
	suite.unsetenv("STRING_VALUES")
	suite.unsetenv("SECRET_WATER")
	suite.setenv("SECRET_FIRE", "Eru_Ilúvatar")
	suite.setenv("SECRET_RINGS", "1")
	suite.setenv("PLUGIN_VALUES", "fire=$SECRET_FIRE,water=${SECRET_WATER}")
	suite.setenv("PLUGIN_STRING_VALUES", "rings=${SECRET_RINGS}")
	suite.setenv("PLUGIN_ADD_REPOS", "testrepo=https://user:${SECRET_FIRE}@testrepo.test")

	cfg, err := NewConfig(&strings.Builder{}, &strings.Builder{})
	suite.Require().NoError(err)

	suite.Equal("fire=Eru_Ilúvatar,water=", cfg.Values)
	suite.Equal("rings=1", cfg.StringValues)
	suite.Equal(fmt.Sprintf("testrepo=https://user:%s@testrepo.test", os.Getenv("SECRET_FIRE")), cfg.AddRepos[0])
}

func (suite *ConfigTestSuite) TestValuesSecretsWithDebugLogging() {
	suite.unsetenv("VALUES")
	suite.unsetenv("SECRET_WATER")
	suite.setenv("SECRET_FIRE", "Eru_Ilúvatar")
	suite.setenv("PLUGIN_DEBUG", "true")
	suite.setenv("PLUGIN_STRING_VALUES", "fire=$SECRET_FIRE")
	suite.setenv("PLUGIN_VALUES", "fire=$SECRET_FIRE,water=$SECRET_WATER")
	stderr := strings.Builder{}
	_, err := NewConfig(&strings.Builder{}, &stderr)
	suite.Require().NoError(err)

	suite.Contains(stderr.String(), "Values:fire=Eru_Ilúvatar,water=")
	suite.Contains(stderr.String(), `$SECRET_WATER not present in environment, replaced with ""`)
}

func (suite *ConfigTestSuite) TestHistoryMax() {
	conf := NewTestConfig(suite.T())
	suite.Assert().Equal(10, conf.HistoryMax)

	suite.setenv("PLUGIN_HISTORY_MAX", "0")
	conf = NewTestConfig(suite.T())
	suite.Assert().Equal(0, conf.HistoryMax)
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

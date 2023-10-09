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
	prefix := getEnvPrefix(os.Stdout)

	suite.unsetenv("MODE")
	suite.unsetenv("UPDATE_DEPENDENCIES")
	suite.unsetenv("DEBUG")

	suite.setenv(prefix+"_MODE", "iambic")
	suite.setenv(prefix+"_UPDATE_DEPENDENCIES", "true")
	suite.setenv(prefix+"_DEBUG", "true")

	cfg, err := NewConfig(&strings.Builder{}, &strings.Builder{})
	suite.Require().NoError(err)

	suite.Equal("iambic", cfg.Command)
	suite.True(cfg.UpdateDependencies)
	suite.True(cfg.Debug)
}

func (suite *ConfigTestSuite) TestNewConfigWithNoPrefix() {
	prefix := getEnvPrefix(os.Stdout)

	suite.unsetenv(prefix + "_MODE")
	suite.unsetenv(prefix + "_UPDATE_DEPENDENCIES")
	suite.unsetenv(prefix + "_DEBUG")

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
	prefix := getEnvPrefix(os.Stdout)

	suite.setenv(prefix+"_MODE", "iambic")
	suite.setenv("MODE", "haiku") // values from the `environment` block override those from `settings`

	cfg, err := NewConfig(&strings.Builder{}, &strings.Builder{})
	suite.Require().NoError(err)

	suite.Equal("haiku", cfg.Command)
}

func (suite *ConfigTestSuite) TestNewConfigInfersNumbersAreSeconds() {
	prefix := getEnvPrefix(os.Stdout)
	suite.setenv(prefix+"_TIMEOUT", "42")
	cfg, err := NewConfig(&strings.Builder{}, &strings.Builder{})
	suite.Require().NoError(err)
	suite.Equal("42s", cfg.Timeout)
}

func (suite *ConfigTestSuite) TestNewConfigWithAliases() {
	prefix := getEnvPrefix(os.Stdout)

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
		suite.unsetenv(prefix + "_" + varname)
	}

	suite.setenv(prefix+"_HELM_COMMAND", "beware the jabberwock")
	suite.setenv(prefix+"_HELM_REPOS", "chortle=http://calloo.callay/frabjous/day")
	suite.setenv(prefix+"_API_SERVER", "http://tumtum.tree")
	suite.setenv(prefix+"_SERVICE_ACCOUNT", "tulgey")
	suite.setenv(prefix+"_WAIT", "true")
	suite.setenv(prefix+"_FORCE", "true")
	suite.setenv(prefix+"_KUBERNETES_TOKEN", "Y29tZSB0byBteSBhcm1z")
	suite.setenv(prefix+"_KUBERNETES_CERTIFICATE", "d2l0aCBpdHMgaGVhZA==")

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
	prefix := getEnvPrefix(os.Stdout)

	suite.unsetenv("FORCE_UPGRADE")
	suite.unsetenv(prefix + "_FORCE_UPGRADE")
	suite.unsetenv(prefix + "_FORCE")
	suite.setenv("FORCE", "true")

	cfg, err := NewConfig(&strings.Builder{}, &strings.Builder{})
	suite.Require().NoError(err)
	suite.True(cfg.Force)
}

func (suite *ConfigTestSuite) TestNewConfigWithAliasConflicts() {
	prefix := getEnvPrefix(os.Stdout)

	suite.unsetenv("FORCE_UPGRADE")
	suite.setenv(prefix+"_FORCE", "true")
	suite.setenv(prefix+"_FORCE_UPGRADE", "false") // should override even when set to the zero value

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
	prefix := getEnvPrefix(os.Stdout)

	for _, varname := range deprecatedVars {
		suite.setenv(varname, "deprecoat") // environment-block entries should cause warnings
	}

	suite.unsetenv("PURGE")
	suite.setenv(prefix+"_PURGE", "true") // settings-block entries should cause warnings
	suite.setenv("UPGRADE", "")           // entries should cause warnings even when set to empty string

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

	if os.Getenv("runner") == "github" {
		suite.Equal("Info: running in github runner, `runner` environment set to 'github'.\n", stdout.String())
	} else {
		suite.Equal("", stdout.String())
	}
	fmt.Println(stderr.String())
	suite.Regexp(`^Generated config: \{Command:upgrade(.|\n)*\}`, stderr.String())
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
	prefix := getEnvPrefix(os.Stdout)

	suite.unsetenv("VALUES")
	suite.unsetenv("STRING_VALUES")
	suite.unsetenv("SECRET_WATER")
	suite.setenv("SECRET_FIRE", "Eru_Ilúvatar")
	suite.setenv("SECRET_RINGS", "1")
	suite.setenv(prefix+"_VALUES", "fire=$SECRET_FIRE,water=${SECRET_WATER}")
	suite.setenv(prefix+"_STRING_VALUES", "rings=${SECRET_RINGS}")
	suite.setenv(prefix+"_ADD_REPOS", "testrepo=https://user:${SECRET_FIRE}@testrepo.test")

	cfg, err := NewConfig(&strings.Builder{}, &strings.Builder{})
	suite.Require().NoError(err)

	suite.Equal("fire=Eru_Ilúvatar,water=", cfg.Values)
	suite.Equal("rings=1", cfg.StringValues)
	suite.Equal(fmt.Sprintf("testrepo=https://user:%s@testrepo.test", os.Getenv("SECRET_FIRE")), cfg.AddRepos[0])
}

func (suite *ConfigTestSuite) TestValuesSecretsWithDebugLogging() {
	prefix := getEnvPrefix(os.Stdout)

	suite.unsetenv("VALUES")
	suite.unsetenv("SECRET_WATER")
	suite.setenv("SECRET_FIRE", "Eru_Ilúvatar")
	suite.setenv(prefix+"_DEBUG", "true")
	suite.setenv(prefix+"_STRING_VALUES", "fire=$SECRET_FIRE")
	suite.setenv(prefix+"_VALUES", "fire=$SECRET_FIRE,water=$SECRET_WATER")
	stderr := strings.Builder{}
	_, err := NewConfig(&strings.Builder{}, &stderr)
	suite.Require().NoError(err)

	suite.Contains(stderr.String(), "Values:fire=Eru_Ilúvatar,water=")
	suite.Contains(stderr.String(), `$SECRET_WATER not present in environment, replaced with ""`)
}

func (suite *ConfigTestSuite) TestHistoryMax() {
	conf := NewTestConfig(suite.T())
	prefix := getEnvPrefix(os.Stdout)

	suite.Assert().Equal(10, conf.HistoryMax)
	suite.setenv(prefix+"_HISTORY_MAX", "0")

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

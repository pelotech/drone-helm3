package helm

import (
	"fmt"
	"github.com/stretchr/testify/suite"
	"os"
	"strings"
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

func (suite *ConfigTestSuite) TestNewConfigWithPluginPrefix() {
	stdout := strings.Builder{}
	stderr := strings.Builder{}
	for _, varname := range []string{
		"MODE",
		"DRONE_BUILD_EVENT",
		"HELM_COMMAND",
		"PLUGIN_HELM_COMMAND",
		"UPDATE_DEPENDENCIES",
		"ADD_REPOS",
		"HELM_REPOS",
		"PLUGIN_HELM_REPOS",
		"DEBUG",
		"VALUES",
		"STRING_VALUES",
		"VALUES_FILES",
		"NAMESPACE",
		"KUBE_TOKEN",
		"KUBERNETES_TOKEN",
		"PLUGIN_KUBERNETES_TOKEN",
		"SKIP_TLS_VERIFY",
		"KUBE_CERTIFICATE",
		"KUBERNETES_CERTIFICATE",
		"PLUGIN_KUBERNETES_CERTIFICATE",
		"KUBE_API_SERVER",
		"API_SERVER",
		"PLUGIN_API_SERVER",
		"KUBE_SERVICE_ACCOUNT",
		"SERVICE_ACCOUNT",
		"PLUGIN_SERVICE_ACCOUNT",
		"CHART_VERSION",
		"DRY_RUN",
		"WAIT_FOR_UPGRADE",
		"WAIT",
		"PLUGIN_WAIT",
		"REUSE_VALUES",
		"KEEP_HISTORY",
		"TIMEOUT",
		"CHART",
		"RELEASE",
		"FORCE",
		"FORCE_UPGRADE",
		"PLUGIN_FORCE_UPGRADE",
		"ATOMIC_UPGRADE",
		"CLEANUP_FAILED_UPGRADE",
		"LINT_STRICTLY",
	} {
		suite.unsetenv(varname)
	}

	suite.setenv("PLUGIN_MODE", "upgrade")
	suite.setenv("PLUGIN_UPDATE_DEPENDENCIES", "true")
	suite.setenv("PLUGIN_ADD_REPOS", "foo=http://bar,goo=http://baz")
	suite.setenv("PLUGIN_DEBUG", "true")
	suite.setenv("PLUGIN_VALUES", "dog=husky")
	suite.setenv("PLUGIN_STRING_VALUES", "version=1.0")
	suite.setenv("PLUGIN_VALUES_FILES", "underrides.yml,overrides.yml")
	suite.setenv("PLUGIN_NAMESPACE", "myapp")
	suite.setenv("PLUGIN_KUBE_TOKEN", "cGxlYXNlIHNpciwgbGV0IG1lIGlu")
	suite.setenv("PLUGIN_SKIP_TLS_VERIFY", "true")
	suite.setenv("PLUGIN_KUBE_CERTIFICATE", "SSBhbSB0b3RhbGx5IHRoZSBzZXJ2ZXIgeW91IHdhbnQ=")
	suite.setenv("PLUGIN_KUBE_API_SERVER", "http://my.kube/cluster")
	suite.setenv("PLUGIN_KUBE_SERVICE_ACCOUNT", "deploybot")
	suite.setenv("PLUGIN_CHART_VERSION", "six")
	suite.setenv("PLUGIN_DRY_RUN", "true")
	suite.setenv("PLUGIN_WAIT_FOR_UPGRADE", "true")
	suite.setenv("PLUGIN_REUSE_VALUES", "true")
	suite.setenv("PLUGIN_KEEP_HISTORY", "true")
	suite.setenv("PLUGIN_TIMEOUT", "5m20s")
	suite.setenv("PLUGIN_CHART", "./helm/myapp/")
	suite.setenv("PLUGIN_RELEASE", "my_app")
	suite.setenv("PLUGIN_FORCE_UPGRADE", "true")
	suite.setenv("PLUGIN_ATOMIC_UPGRADE", "true")
	suite.setenv("PLUGIN_CLEANUP_FAILED_UPGRADE", "true")
	suite.setenv("PLUGIN_LINT_STRICTLY", "true")

	cfg, err := NewConfig(&stdout, &stderr)
	suite.Require().NoError(err)

	want := Config{
		Command:            "upgrade",
		DroneEvent:         "",
		UpdateDependencies: true,
		AddRepos:           []string{"foo=http://bar", "goo=http://baz"},
		Debug:              true,
		Values:             "dog=husky",
		StringValues:       "version=1.0",
		ValuesFiles:        []string{"underrides.yml", "overrides.yml"},
		Namespace:          "myapp",
		KubeToken:          "cGxlYXNlIHNpciwgbGV0IG1lIGlu",
		SkipTLSVerify:      true,
		Certificate:        "SSBhbSB0b3RhbGx5IHRoZSBzZXJ2ZXIgeW91IHdhbnQ=",
		APIServer:          "http://my.kube/cluster",
		ServiceAccount:     "deploybot",
		ChartVersion:       "six",
		DryRun:             true,
		Wait:               true,
		ReuseValues:        true,
		KeepHistory:        true,
		Timeout:            "5m20s",
		Chart:              "./helm/myapp/",
		Release:            "my_app",
		Force:              true,
		AtomicUpgrade:      true,
		CleanupOnFail:      true,
		LintStrictly:       true,
		Stdout:             &stdout,
		Stderr:             &stderr,
	}

	suite.Equal(&want, cfg)
}

func (suite *ConfigTestSuite) TestNewConfigWithNoPrefix() {
	stdout := strings.Builder{}
	stderr := strings.Builder{}
	for _, varname := range []string{
		"PLUGIN_MODE",
		"PLUGIN_HELM_COMMAND",
		"HELM_COMMAND",
		"PLUGIN_UPDATE_DEPENDENCIES",
		"PLUGIN_ADD_REPOS",
		"PLUGIN_HELM_REPOS",
		"HELM_REPOS",
		"PLUGIN_DEBUG",
		"PLUGIN_VALUES",
		"PLUGIN_STRING_VALUES",
		"PLUGIN_VALUES_FILES",
		"PLUGIN_NAMESPACE",
		"PLUGIN_KUBE_TOKEN",
		"PLUGIN_KUBERNETES_TOKEN",
		"KUBERNETES_TOKEN",
		"PLUGIN_SKIP_TLS_VERIFY",
		"PLUGIN_KUBE_CERTIFICATE",
		"PLUGIN_KUBERNETES_CERTIFICATE",
		"KUBERNETES_CERTIFICATE",
		"PLUGIN_KUBE_API_SERVER",
		"PLUGIN_API_SERVER",
		"API_SERVER",
		"PLUGIN_KUBE_SERVICE_ACCOUNT",
		"PLUGIN_SERVICE_ACCOUNT",
		"SERVICE_ACCOUNT",
		"PLUGIN_CHART_VERSION",
		"PLUGIN_DRY_RUN",
		"PLUGIN_WAIT_FOR_UPGRADE",
		"PLUGIN_WAIT",
		"WAIT",
		"PLUGIN_REUSE_VALUES",
		"PLUGIN_KEEP_HISTORY",
		"PLUGIN_TIMEOUT",
		"PLUGIN_CHART",
		"PLUGIN_RELEASE",
		"PLUGIN_FORCE",
		"PLUGIN_FORCE_UPGRADE",
		"FORCE_UPGRADE",
		"PLUGIN_ATOMIC_UPGRADE",
		"PLUGIN_CLEANUP_FAILED_UPGRADE",
		"PLUGIN_LINT_STRICTLY",
	} {
		suite.unsetenv(varname)
	}

	suite.setenv("MODE", "upgrade")
	suite.setenv("DRONE_BUILD_EVENT", "tag")
	suite.setenv("UPDATE_DEPENDENCIES", "true")
	suite.setenv("ADD_REPOS", "foo=http://bar,goo=http://baz")
	suite.setenv("DEBUG", "true")
	suite.setenv("VALUES", "dog=husky")
	suite.setenv("STRING_VALUES", "version=1.0")
	suite.setenv("VALUES_FILES", "underrides.yml,overrides.yml")
	suite.setenv("NAMESPACE", "myapp")
	suite.setenv("KUBE_TOKEN", "cGxlYXNlIHNpciwgbGV0IG1lIGlu")
	suite.setenv("SKIP_TLS_VERIFY", "true")
	suite.setenv("KUBE_CERTIFICATE", "SSBhbSB0b3RhbGx5IHRoZSBzZXJ2ZXIgeW91IHdhbnQ=")
	suite.setenv("KUBE_API_SERVER", "http://my.kube/cluster")
	suite.setenv("KUBE_SERVICE_ACCOUNT", "deploybot")
	suite.setenv("CHART_VERSION", "six")
	suite.setenv("DRY_RUN", "true")
	suite.setenv("WAIT_FOR_UPGRADE", "true")
	suite.setenv("REUSE_VALUES", "true")
	suite.setenv("KEEP_HISTORY", "true")
	suite.setenv("TIMEOUT", "5m20s")
	suite.setenv("CHART", "./helm/myapp/")
	suite.setenv("RELEASE", "my_app")
	suite.setenv("FORCE_UPGRADE", "true")
	suite.setenv("ATOMIC_UPGRADE", "true")
	suite.setenv("CLEANUP_FAILED_UPGRADE", "true")
	suite.setenv("LINT_STRICTLY", "true")

	cfg, err := NewConfig(&stdout, &stderr)
	suite.Require().NoError(err)

	want := Config{
		Command:            "upgrade",
		DroneEvent:         "tag",
		UpdateDependencies: true,
		AddRepos:           []string{"foo=http://bar", "goo=http://baz"},
		Debug:              true,
		Values:             "dog=husky",
		StringValues:       "version=1.0",
		ValuesFiles:        []string{"underrides.yml", "overrides.yml"},
		Namespace:          "myapp",
		KubeToken:          "cGxlYXNlIHNpciwgbGV0IG1lIGlu",
		SkipTLSVerify:      true,
		Certificate:        "SSBhbSB0b3RhbGx5IHRoZSBzZXJ2ZXIgeW91IHdhbnQ=",
		APIServer:          "http://my.kube/cluster",
		ServiceAccount:     "deploybot",
		ChartVersion:       "six",
		DryRun:             true,
		Wait:               true,
		ReuseValues:        true,
		KeepHistory:        true,
		Timeout:            "5m20s",
		Chart:              "./helm/myapp/",
		Release:            "my_app",
		Force:              true,
		AtomicUpgrade:      true,
		CleanupOnFail:      true,
		LintStrictly:       true,
		Stdout:             &stdout,
		Stderr:             &stderr,
	}

	suite.Equal(&want, cfg)
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

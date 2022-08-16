package run

import (
	"fmt"
	"strings"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/pelotech/drone-helm3/internal/env"
	"github.com/stretchr/testify/suite"
)

type UpgradeTestSuite struct {
	suite.Suite
	ctrl            *gomock.Controller
	mockCmd         *Mockcmd
	originalCommand func(string, ...string) cmd
}

func (suite *UpgradeTestSuite) BeforeTest(_, _ string) {
	suite.ctrl = gomock.NewController(suite.T())
	suite.mockCmd = NewMockcmd(suite.ctrl)

	suite.originalCommand = command
	command = func(path string, args ...string) cmd { return suite.mockCmd }
}

func (suite *UpgradeTestSuite) AfterTest(_, _ string) {
	command = suite.originalCommand
}

func TestUpgradeTestSuite(t *testing.T) {
	suite.Run(t, new(UpgradeTestSuite))
}

func (suite *UpgradeTestSuite) TestNewUpgrade() {
	cfg := env.NewTestConfig(suite.T())
	cfg.ChartVersion = "seventeen"
	cfg.DryRun = true
	cfg.Wait = true
	cfg.Values = "steadfastness,forthrightness"
	cfg.StringValues = "tensile_strength,flexibility"
	cfg.ValuesFiles = []string{"/root/price_inventory.yml"}
	cfg.ReuseValues = true
	cfg.Timeout = "go sit in the corner"
	cfg.Chart = "billboard_top_100"
	cfg.Release = "post_malone_circles"
	cfg.Force = true
	cfg.AtomicUpgrade = true
	cfg.CleanupOnFail = true

	up := NewUpgrade(*cfg)

	suite.Equal(cfg.Chart, up.chart)
	suite.Equal(cfg.Release, up.release)
	suite.Equal(cfg.ChartVersion, up.chartVersion)
	suite.Equal(true, up.dryRun)
	suite.Equal(cfg.Wait, up.wait)
	suite.Equal("steadfastness,forthrightness", up.values)
	suite.Equal("tensile_strength,flexibility", up.stringValues)
	suite.Equal([]string{"/root/price_inventory.yml"}, up.valuesFiles)
	suite.Equal(cfg.ReuseValues, up.reuseValues)
	suite.Equal(cfg.Timeout, up.timeout)
	suite.Equal(cfg.Force, up.force)
	suite.Equal(true, up.atomic)
	suite.Equal(true, up.cleanupOnFail)
	suite.NotNil(up.config)
	suite.NotNil(up.certs)
}

func (suite *UpgradeTestSuite) TestPrepareAndExecute() {
	defer suite.ctrl.Finish()

	cfg := env.NewTestConfig(suite.T())
	cfg.Chart = "at40"
	cfg.Release = "jonas_brothers_only_human"

	u := NewUpgrade(*cfg)

	command = func(path string, args ...string) cmd {
		suite.Equal(helmBin, path)
		suite.Equal([]string{"upgrade", "--install", "--history-max=10", "jonas_brothers_only_human", "at40"}, args)

		return suite.mockCmd
	}

	suite.mockCmd.EXPECT().
		Stdout(gomock.Any())
	suite.mockCmd.EXPECT().
		Stderr(gomock.Any())
	suite.mockCmd.EXPECT().
		Run().
		Times(1)

	err := u.Prepare()
	suite.Require().Nil(err)
	u.Execute()
}

func (suite *UpgradeTestSuite) TestPrepareNamespaceFlag() {
	defer suite.ctrl.Finish()

	cfg := env.NewTestConfig(suite.T())
	cfg.Namespace = "melt"
	cfg.Chart = "at40"
	cfg.Release = "shaed_trampoline"

	u := NewUpgrade(*cfg)

	command = func(path string, args ...string) cmd {
		suite.Equal(helmBin, path)
		suite.Equal([]string{"--namespace", "melt", "upgrade", "--install", "--history-max=10", "shaed_trampoline", "at40"}, args)

		return suite.mockCmd
	}

	suite.mockCmd.EXPECT().Stdout(gomock.Any())
	suite.mockCmd.EXPECT().Stderr(gomock.Any())

	err := u.Prepare()
	suite.Require().Nil(err)
}

func (suite *UpgradeTestSuite) TestPrepareWithUpgradeFlags() {
	defer suite.ctrl.Finish()

	cfg := env.NewTestConfig(suite.T())
	cfg.Chart = "hot_ac"
	cfg.Release = "maroon_5_memories"
	cfg.ChartVersion = "radio_edit"
	cfg.DryRun = true
	cfg.Wait = true
	cfg.Values = "age=35"
	cfg.StringValues = "height=5ft10in"
	cfg.ValuesFiles = []string{"/usr/local/stats", "/usr/local/grades"}
	cfg.ReuseValues = true
	cfg.Timeout = "sit_in_the_corner"
	cfg.Force = true
	cfg.AtomicUpgrade = true
	cfg.CleanupOnFail = true

	u := NewUpgrade(*cfg)
	// inject a ca cert filename so repoCerts won't create any files that we'd have to clean up
	u.certs.caCertFilename = "local_ca.cert"

	command = func(path string, args ...string) cmd {
		suite.Equal(helmBin, path)
		suite.Equal([]string{"upgrade", "--install",
			"--version", "radio_edit",
			"--dry-run",
			"--wait",
			"--reuse-values",
			"--timeout", "sit_in_the_corner",
			"--force",
			"--atomic",
			"--cleanup-on-fail",
			"--set", "age=35",
			"--set-string", "height=5ft10in",
			"--values", "/usr/local/stats",
			"--values", "/usr/local/grades",
			"--ca-file", "local_ca.cert",
			"--history-max=10",
			"maroon_5_memories", "hot_ac"}, args)

		return suite.mockCmd
	}

	suite.mockCmd.EXPECT().Stdout(gomock.Any())
	suite.mockCmd.EXPECT().Stderr(gomock.Any())

	err := u.Prepare()
	suite.Require().Nil(err)
}

func (suite *UpgradeTestSuite) TestRequiresChartAndRelease() {
	// These aren't really expected, but allowing them gives clearer test-failure messages
	suite.mockCmd.EXPECT().Stdout(gomock.Any()).AnyTimes()
	suite.mockCmd.EXPECT().Stderr(gomock.Any()).AnyTimes()

	u := NewUpgrade(env.Config{})
	u.release = "seth_everman_unskippable_cutscene"

	err := u.Prepare()
	suite.EqualError(err, "chart is required", "Chart should be mandatory")

	u.release = ""
	u.chart = "billboard_top_zero"

	err = u.Prepare()
	suite.EqualError(err, "release is required", "Release should be mandatory")
}

func (suite *UpgradeTestSuite) TestPrepareDebugFlag() {
	stdout := strings.Builder{}
	stderr := strings.Builder{}

	cfg := env.NewTestConfig(suite.T())
	cfg.Chart = "at40"
	cfg.Release = "lewis_capaldi_someone_you_loved"
	cfg.Debug = true
	cfg.Stdout = &stdout
	cfg.Stderr = &stderr

	u := NewUpgrade(*cfg)

	command = func(path string, args ...string) cmd {
		suite.mockCmd.EXPECT().
			String().
			Return(fmt.Sprintf("%s %s", path, strings.Join(args, " ")))

		return suite.mockCmd
	}

	suite.mockCmd.EXPECT().
		Stdout(&stdout)
	suite.mockCmd.EXPECT().
		Stderr(&stderr)

	u.Prepare()

	want := fmt.Sprintf(
		"Generated command: '%s --debug upgrade --install --history-max=10 lewis_capaldi_someone_you_loved at40'\n",
		helmBin,
	)
	suite.Equal(want, stderr.String())
	suite.Equal("", stdout.String())
}

func (suite *UpgradeTestSuite) TestPrepareSkipCrdsFlag() {
	defer suite.ctrl.Finish()

	cfg := env.NewTestConfig(suite.T())
	cfg.Chart = "at40"
	cfg.Release = "cabbages_smell_great"
	cfg.SkipCrds = true

	u := NewUpgrade(*cfg)

	command = func(path string, args ...string) cmd {
		suite.Equal(helmBin, path)
		suite.Equal([]string{"upgrade", "--install", "--skip-crds", "--history-max=10", "cabbages_smell_great", "at40"}, args)

		return suite.mockCmd
	}

	suite.mockCmd.EXPECT().Stdout(gomock.Any())
	suite.mockCmd.EXPECT().Stderr(gomock.Any())

	err := u.Prepare()
	suite.Require().Nil(err)
}

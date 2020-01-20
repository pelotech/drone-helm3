package run

import (
	"fmt"
	"github.com/golang/mock/gomock"
	"github.com/pelotech/drone-helm3/internal/env"
	"github.com/stretchr/testify/suite"
	"strings"
	"testing"
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
	cfg := env.Config{
		ChartVersion:  "seventeen",
		DryRun:        true,
		Wait:          true,
		Values:        "steadfastness,forthrightness",
		StringValues:  "tensile_strength,flexibility",
		ValuesFiles:   []string{"/root/price_inventory.yml"},
		ReuseValues:   true,
		Timeout:       "go sit in the corner",
		Chart:         "billboard_top_100",
		Release:       "post_malone_circles",
		Force:         true,
		AtomicUpgrade: true,
		CleanupOnFail: true,
	}

	up := NewUpgrade(cfg)
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
}

func (suite *UpgradeTestSuite) TestPrepareAndExecute() {
	defer suite.ctrl.Finish()

	cfg := env.Config{
		Chart:   "at40",
		Release: "jonas_brothers_only_human",
	}
	u := NewUpgrade(cfg)

	command = func(path string, args ...string) cmd {
		suite.Equal(helmBin, path)
		suite.Equal([]string{"upgrade", "--install", "jonas_brothers_only_human", "at40"}, args)

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

	cfg := env.Config{
		Namespace: "melt",
		Chart:     "at40",
		Release:   "shaed_trampoline",
	}
	u := NewUpgrade(cfg)

	command = func(path string, args ...string) cmd {
		suite.Equal(helmBin, path)
		suite.Equal([]string{"--namespace", "melt", "upgrade", "--install", "shaed_trampoline", "at40"}, args)

		return suite.mockCmd
	}

	suite.mockCmd.EXPECT().Stdout(gomock.Any())
	suite.mockCmd.EXPECT().Stderr(gomock.Any())

	err := u.Prepare()
	suite.Require().Nil(err)
}

func (suite *UpgradeTestSuite) TestPrepareWithUpgradeFlags() {
	defer suite.ctrl.Finish()

	cfg := env.Config{
		Chart:         "hot_ac",
		Release:       "maroon_5_memories",
		ChartVersion:  "radio_edit",
		DryRun:        true,
		Wait:          true,
		Values:        "age=35",
		StringValues:  "height=5ft10in",
		ValuesFiles:   []string{"/usr/local/stats", "/usr/local/grades"},
		ReuseValues:   true,
		Timeout:       "sit_in_the_corner",
		Force:         true,
		AtomicUpgrade: true,
		CleanupOnFail: true,
		RepoCAFile:    "local_ca.cert",
	}
	u := NewUpgrade(cfg)

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
	cfg := env.Config{
		Chart:   "at40",
		Release: "lewis_capaldi_someone_you_loved",
		Debug:   true,
		Stdout:  &stdout,
		Stderr:  &stderr,
	}
	u := NewUpgrade(cfg)

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

	want := fmt.Sprintf("Generated command: '%s --debug upgrade "+
		"--install lewis_capaldi_someone_you_loved at40'\n", helmBin)
	suite.Equal(want, stderr.String())
	suite.Equal("", stdout.String())
}

package run

import (
	"fmt"
	"github.com/golang/mock/gomock"
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

func (suite *UpgradeTestSuite) TestPrepareAndExecute() {
	defer suite.ctrl.Finish()

	u := Upgrade{
		Chart:   "at40",
		Release: "jonas_brothers_only_human",
	}

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

	cfg := Config{}
	err := u.Prepare(cfg)
	suite.Require().Nil(err)
	u.Execute(cfg)
}

func (suite *UpgradeTestSuite) TestPrepareNamespaceFlag() {
	defer suite.ctrl.Finish()

	u := Upgrade{
		Chart:   "at40",
		Release: "shaed_trampoline",
	}

	command = func(path string, args ...string) cmd {
		suite.Equal(helmBin, path)
		suite.Equal([]string{"--namespace", "melt", "upgrade", "--install", "shaed_trampoline", "at40"}, args)

		return suite.mockCmd
	}

	suite.mockCmd.EXPECT().Stdout(gomock.Any())
	suite.mockCmd.EXPECT().Stderr(gomock.Any())

	cfg := Config{
		Namespace: "melt",
	}
	err := u.Prepare(cfg)
	suite.Require().Nil(err)
}

func (suite *UpgradeTestSuite) TestPrepareWithUpgradeFlags() {
	defer suite.ctrl.Finish()

	u := Upgrade{
		Chart:        "hot_ac",
		Release:      "maroon_5_memories",
		ChartVersion: "radio_edit",
		DryRun:       true,
		Wait:         true,
		Values:       "age=35",
		StringValues: "height=5ft10in",
		ValuesFiles:  []string{"/usr/local/stats", "/usr/local/grades"},
		ReuseValues:  true,
		Timeout:      "sit_in_the_corner",
		Force:        true,
		Atomic:       true,
	}

	cfg := Config{}

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
			"--set", "age=35",
			"--set-string", "height=5ft10in",
			"--values", "/usr/local/stats",
			"--values", "/usr/local/grades",
			"maroon_5_memories", "hot_ac"}, args)

		return suite.mockCmd
	}

	suite.mockCmd.EXPECT().Stdout(gomock.Any())
	suite.mockCmd.EXPECT().Stderr(gomock.Any())

	err := u.Prepare(cfg)
	suite.Require().Nil(err)
}

func (suite *UpgradeTestSuite) TestRequiresChartAndRelease() {
	// These aren't really expected, but allowing them gives clearer test-failure messages
	suite.mockCmd.EXPECT().Stdout(gomock.Any()).AnyTimes()
	suite.mockCmd.EXPECT().Stderr(gomock.Any()).AnyTimes()

	u := Upgrade{
		Release: "seth_everman_unskippable_cutscene",
	}

	err := u.Prepare(Config{})
	suite.EqualError(err, "chart is required", "Chart should be mandatory")

	u = Upgrade{
		Chart: "billboard_top_zero",
	}

	err = u.Prepare(Config{})
	suite.EqualError(err, "release is required", "Release should be mandatory")
}

func (suite *UpgradeTestSuite) TestPrepareDebugFlag() {
	u := Upgrade{
		Chart:   "at40",
		Release: "lewis_capaldi_someone_you_loved",
	}

	stdout := strings.Builder{}
	stderr := strings.Builder{}
	cfg := Config{
		Debug:  true,
		Stdout: &stdout,
		Stderr: &stderr,
	}

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

	u.Prepare(cfg)

	want := fmt.Sprintf("Generated command: '%s --debug upgrade "+
		"--install lewis_capaldi_someone_you_loved at40'\n", helmBin)
	suite.Equal(want, stderr.String())
	suite.Equal("", stdout.String())
}

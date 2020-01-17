package run

import (
	"fmt"
	"github.com/golang/mock/gomock"
	"github.com/pelotech/drone-helm3/internal/env"
	"github.com/stretchr/testify/suite"
	"strings"
	"testing"
)

type UninstallTestSuite struct {
	suite.Suite
	ctrl            *gomock.Controller
	mockCmd         *Mockcmd
	actualArgs      []string
	originalCommand func(string, ...string) cmd
}

func (suite *UninstallTestSuite) BeforeTest(_, _ string) {
	suite.ctrl = gomock.NewController(suite.T())
	suite.mockCmd = NewMockcmd(suite.ctrl)

	suite.originalCommand = command
	command = func(path string, args ...string) cmd {
		suite.actualArgs = args
		return suite.mockCmd
	}
}

func (suite *UninstallTestSuite) AfterTest(_, _ string) {
	command = suite.originalCommand
}

func TestUninstallTestSuite(t *testing.T) {
	suite.Run(t, new(UninstallTestSuite))
}

func (suite *UninstallTestSuite) TestNewUninstall() {
	cfg := env.Config{
		DryRun:      true,
		Release:     "jetta_id_love_to_change_the_world",
		KeepHistory: true,
	}
	u := NewUninstall(cfg)
	suite.Equal("jetta_id_love_to_change_the_world", u.release)
	suite.Equal(true, u.dryRun)
	suite.Equal(true, u.keepHistory)
	suite.NotNil(u.config)
}

func (suite *UninstallTestSuite) TestPrepareAndExecute() {
	defer suite.ctrl.Finish()

	cfg := env.Config{
		Release: "zayde_wølf_king",
	}
	u := NewUninstall(cfg)

	actual := []string{}
	command = func(path string, args ...string) cmd {
		suite.Equal(helmBin, path)
		actual = args

		return suite.mockCmd
	}

	suite.mockCmd.EXPECT().
		Stdout(gomock.Any())
	suite.mockCmd.EXPECT().
		Stderr(gomock.Any())
	suite.mockCmd.EXPECT().
		Run().
		Times(1)

	suite.NoError(u.Prepare())
	expected := []string{"uninstall", "zayde_wølf_king"}
	suite.Equal(expected, actual)

	u.Execute()
}

func (suite *UninstallTestSuite) TestPrepareDryRunFlag() {
	cfg := env.Config{
		Release: "firefox_ak_wildfire",
		DryRun:  true,
	}
	u := NewUninstall(cfg)

	suite.mockCmd.EXPECT().Stdout(gomock.Any()).AnyTimes()
	suite.mockCmd.EXPECT().Stderr(gomock.Any()).AnyTimes()

	suite.NoError(u.Prepare())
	expected := []string{"uninstall", "--dry-run", "firefox_ak_wildfire"}
	suite.Equal(expected, suite.actualArgs)
}

func (suite *UninstallTestSuite) TestPrepareKeepHistoryFlag() {
	cfg := env.Config{
		Release:     "perturbator_sentient",
		KeepHistory: true,
	}
	u := NewUninstall(cfg)

	suite.mockCmd.EXPECT().Stdout(gomock.Any()).AnyTimes()
	suite.mockCmd.EXPECT().Stderr(gomock.Any()).AnyTimes()

	suite.NoError(u.Prepare())
	expected := []string{"uninstall", "--keep-history", "perturbator_sentient"}
	suite.Equal(expected, suite.actualArgs)
}

func (suite *UninstallTestSuite) TestPrepareNamespaceFlag() {
	cfg := env.Config{
		Release:   "carly_simon_run_away_with_me",
		Namespace: "emotion",
	}
	u := NewUninstall(cfg)

	suite.mockCmd.EXPECT().Stdout(gomock.Any()).AnyTimes()
	suite.mockCmd.EXPECT().Stderr(gomock.Any()).AnyTimes()

	suite.NoError(u.Prepare())
	expected := []string{"--namespace", "emotion", "uninstall", "carly_simon_run_away_with_me"}
	suite.Equal(expected, suite.actualArgs)
}

func (suite *UninstallTestSuite) TestPrepareDebugFlag() {
	stderr := strings.Builder{}
	cfg := env.Config{
		Release: "just_a_band_huff_and_puff",
		Debug:   true,
		Stderr:  &stderr,
	}
	u := NewUninstall(cfg)

	command = func(path string, args ...string) cmd {
		suite.mockCmd.EXPECT().
			String().
			Return(fmt.Sprintf("%s %s", path, strings.Join(args, " ")))

		return suite.mockCmd
	}

	suite.mockCmd.EXPECT().Stdout(gomock.Any()).AnyTimes()
	suite.mockCmd.EXPECT().Stderr(&stderr).AnyTimes()

	suite.NoError(u.Prepare())
	suite.Equal(fmt.Sprintf("Generated command: '%s --debug "+
		"uninstall just_a_band_huff_and_puff'\n", helmBin), stderr.String())
}

func (suite *UninstallTestSuite) TestPrepareRequiresRelease() {
	// These aren't really expected, but allowing them gives clearer test-failure messages
	suite.mockCmd.EXPECT().Stdout(gomock.Any()).AnyTimes()
	suite.mockCmd.EXPECT().Stderr(gomock.Any()).AnyTimes()

	u := NewUninstall(env.Config{})
	err := u.Prepare()
	suite.EqualError(err, "release is required", "Uninstall.Release should be mandatory")
}

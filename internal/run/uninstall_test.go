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
	suite.Equal(&Uninstall{
		Release:     "jetta_id_love_to_change_the_world",
		DryRun:      true,
		KeepHistory: true,
	}, u)
}

func (suite *UninstallTestSuite) TestPrepareAndExecute() {
	defer suite.ctrl.Finish()

	u := Uninstall{
		Release: "zayde_wølf_king",
	}

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

	cfg := Config{}
	suite.NoError(u.Prepare(cfg))
	expected := []string{"uninstall", "zayde_wølf_king"}
	suite.Equal(expected, actual)

	u.Execute()
}

func (suite *UninstallTestSuite) TestPrepareDryRunFlag() {
	u := Uninstall{
		Release: "firefox_ak_wildfire",
		DryRun:  true,
	}
	cfg := Config{}

	suite.mockCmd.EXPECT().Stdout(gomock.Any()).AnyTimes()
	suite.mockCmd.EXPECT().Stderr(gomock.Any()).AnyTimes()

	suite.NoError(u.Prepare(cfg))
	expected := []string{"uninstall", "--dry-run", "firefox_ak_wildfire"}
	suite.Equal(expected, suite.actualArgs)
}

func (suite *UninstallTestSuite) TestPrepareKeepHistoryFlag() {
	u := Uninstall{
		Release:     "perturbator_sentient",
		KeepHistory: true,
	}
	cfg := Config{}

	suite.mockCmd.EXPECT().Stdout(gomock.Any()).AnyTimes()
	suite.mockCmd.EXPECT().Stderr(gomock.Any()).AnyTimes()

	suite.NoError(u.Prepare(cfg))
	expected := []string{"uninstall", "--keep-history", "perturbator_sentient"}
	suite.Equal(expected, suite.actualArgs)
}

func (suite *UninstallTestSuite) TestPrepareNamespaceFlag() {
	u := Uninstall{
		Release: "carly_simon_run_away_with_me",
	}
	cfg := Config{
		Namespace: "emotion",
	}

	suite.mockCmd.EXPECT().Stdout(gomock.Any()).AnyTimes()
	suite.mockCmd.EXPECT().Stderr(gomock.Any()).AnyTimes()

	suite.NoError(u.Prepare(cfg))
	expected := []string{"--namespace", "emotion", "uninstall", "carly_simon_run_away_with_me"}
	suite.Equal(expected, suite.actualArgs)
}

func (suite *UninstallTestSuite) TestPrepareDebugFlag() {
	u := Uninstall{
		Release: "just_a_band_huff_and_puff",
	}
	stderr := strings.Builder{}
	cfg := Config{
		Debug:  true,
		Stderr: &stderr,
	}

	command = func(path string, args ...string) cmd {
		suite.mockCmd.EXPECT().
			String().
			Return(fmt.Sprintf("%s %s", path, strings.Join(args, " ")))

		return suite.mockCmd
	}

	suite.mockCmd.EXPECT().Stdout(gomock.Any()).AnyTimes()
	suite.mockCmd.EXPECT().Stderr(&stderr).AnyTimes()

	suite.NoError(u.Prepare(cfg))
	suite.Equal(fmt.Sprintf("Generated command: '%s --debug "+
		"uninstall just_a_band_huff_and_puff'\n", helmBin), stderr.String())
}

func (suite *UninstallTestSuite) TestPrepareRequiresRelease() {
	// These aren't really expected, but allowing them gives clearer test-failure messages
	suite.mockCmd.EXPECT().Stdout(gomock.Any()).AnyTimes()
	suite.mockCmd.EXPECT().Stderr(gomock.Any()).AnyTimes()

	u := Uninstall{}
	err := u.Prepare(Config{})
	suite.EqualError(err, "release is required", "Uninstall.Release should be mandatory")
}

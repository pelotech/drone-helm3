package run

import (
	"fmt"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/suite"
	"strings"
	"testing"
)

type DeleteTestSuite struct {
	suite.Suite
	ctrl            *gomock.Controller
	mockCmd         *Mockcmd
	actualArgs      []string
	originalCommand func(string, ...string) cmd
}

func (suite *DeleteTestSuite) BeforeTest(_, _ string) {
	suite.ctrl = gomock.NewController(suite.T())
	suite.mockCmd = NewMockcmd(suite.ctrl)

	suite.originalCommand = command
	command = func(path string, args ...string) cmd {
		suite.actualArgs = args
		return suite.mockCmd
	}
}

func (suite *DeleteTestSuite) AfterTest(_, _ string) {
	command = suite.originalCommand
}

func TestDeleteTestSuite(t *testing.T) {
	suite.Run(t, new(DeleteTestSuite))
}

func (suite *DeleteTestSuite) TestPrepareAndExecute() {
	defer suite.ctrl.Finish()

	d := Delete{
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

	cfg := Config{
		KubeConfig: "/root/.kube/config",
	}
	suite.NoError(d.Prepare(cfg))
	expected := []string{"--kubeconfig", "/root/.kube/config", "delete", "zayde_wølf_king"}
	suite.Equal(expected, actual)

	d.Execute(cfg)
}

func (suite *DeleteTestSuite) TestPrepareDryRunFlag() {
	d := Delete{
		Release: "firefox_ak_wildfire",
		DryRun:  true,
	}
	cfg := Config{
		KubeConfig: "/root/.kube/config",
	}

	suite.mockCmd.EXPECT().Stdout(gomock.Any()).AnyTimes()
	suite.mockCmd.EXPECT().Stderr(gomock.Any()).AnyTimes()

	suite.NoError(d.Prepare(cfg))
	expected := []string{"--kubeconfig", "/root/.kube/config", "delete", "--dry-run", "firefox_ak_wildfire"}
	suite.Equal(expected, suite.actualArgs)
}

func (suite *DeleteTestSuite) TestPrepareNamespaceFlag() {
	d := Delete{
		Release: "carly_simon_run_away_with_me",
	}
	cfg := Config{
		KubeConfig: "/root/.kube/config",
		Namespace:  "emotion",
	}

	suite.mockCmd.EXPECT().Stdout(gomock.Any()).AnyTimes()
	suite.mockCmd.EXPECT().Stderr(gomock.Any()).AnyTimes()

	suite.NoError(d.Prepare(cfg))
	expected := []string{"--kubeconfig", "/root/.kube/config",
		"--namespace", "emotion", "delete", "carly_simon_run_away_with_me"}
	suite.Equal(expected, suite.actualArgs)
}

func (suite *DeleteTestSuite) TestPrepareDebugFlag() {
	d := Delete{
		Release: "just_a_band_huff_and_puff",
	}
	stderr := strings.Builder{}
	cfg := Config{
		KubeConfig: "/root/.kube/config",
		Debug:      true,
		Stderr:     &stderr,
	}

	command = func(path string, args ...string) cmd {
		suite.mockCmd.EXPECT().
			String().
			Return(fmt.Sprintf("%s %s", path, strings.Join(args, " ")))

		return suite.mockCmd
	}

	suite.mockCmd.EXPECT().Stdout(gomock.Any()).AnyTimes()
	suite.mockCmd.EXPECT().Stderr(&stderr).AnyTimes()

	suite.NoError(d.Prepare(cfg))
	suite.Equal(fmt.Sprintf("Generated command: '%s --kubeconfig /root/.kube/config "+
		"--debug delete just_a_band_huff_and_puff'\n", helmBin), stderr.String())
}

func (suite *DeleteTestSuite) TestPrepareRequiresRelease() {
	// These aren't really expected, but allowing them gives clearer test-failure messages
	suite.mockCmd.EXPECT().Stdout(gomock.Any()).AnyTimes()
	suite.mockCmd.EXPECT().Stderr(gomock.Any()).AnyTimes()

	d := Delete{}
	err := d.Prepare(Config{})
	suite.EqualError(err, "release is required", "Delete.Release should be mandatory")
}

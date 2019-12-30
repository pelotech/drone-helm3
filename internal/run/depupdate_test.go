package run

import (
	"fmt"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/suite"
	"strings"
	"testing"
)

type DepUpdateTestSuite struct {
	suite.Suite
	ctrl            *gomock.Controller
	mockCmd         *Mockcmd
	originalCommand func(string, ...string) cmd
}

func (suite *DepUpdateTestSuite) BeforeTest(_, _ string) {
	suite.ctrl = gomock.NewController(suite.T())
	suite.mockCmd = NewMockcmd(suite.ctrl)

	suite.originalCommand = command
	command = func(path string, args ...string) cmd { return suite.mockCmd }
}

func (suite *DepUpdateTestSuite) AfterTest(_, _ string) {
	command = suite.originalCommand
}

func TestDepUpdateTestSuite(t *testing.T) {
	suite.Run(t, new(DepUpdateTestSuite))
}

func (suite *DepUpdateTestSuite) TestPrepareAndExecute() {
	defer suite.ctrl.Finish()

	stdout := strings.Builder{}
	stderr := strings.Builder{}
	cfg := Config{
		Stdout: &stdout,
		Stderr: &stderr,
	}

	command = func(path string, args ...string) cmd {
		suite.Equal(helmBin, path)
		suite.Equal([]string{"dependency", "update", "your_top_songs_2019"}, args)

		return suite.mockCmd
	}
	suite.mockCmd.EXPECT().
		Stdout(&stdout)
	suite.mockCmd.EXPECT().
		Stderr(&stderr)
	suite.mockCmd.EXPECT().
		Run().
		Times(1)

	d := DepUpdate{
		Chart: "your_top_songs_2019",
	}

	suite.Require().NoError(d.Prepare(cfg))
	suite.NoError(d.Execute(cfg))
}

func (suite *DepUpdateTestSuite) TestPrepareNamespaceFlag() {
	defer suite.ctrl.Finish()

	cfg := Config{
		Namespace: "spotify",
	}

	command = func(path string, args ...string) cmd {
		suite.Equal([]string{"--namespace", "spotify", "dependency", "update", "your_top_songs_2019"}, args)

		return suite.mockCmd
	}
	suite.mockCmd.EXPECT().Stdout(gomock.Any()).AnyTimes()
	suite.mockCmd.EXPECT().Stderr(gomock.Any()).AnyTimes()

	d := DepUpdate{
		Chart: "your_top_songs_2019",
	}

	suite.Require().NoError(d.Prepare(cfg))
}

func (suite *DepUpdateTestSuite) TestPrepareDebugFlag() {
	defer suite.ctrl.Finish()

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
	suite.mockCmd.EXPECT().Stdout(gomock.Any()).AnyTimes()
	suite.mockCmd.EXPECT().Stderr(gomock.Any()).AnyTimes()

	d := DepUpdate{
		Chart: "your_top_songs_2019",
	}

	suite.Require().NoError(d.Prepare(cfg))

	want := fmt.Sprintf("Generated command: '%s --debug dependency update your_top_songs_2019'\n", helmBin)
	suite.Equal(want, stderr.String())
	suite.Equal("", stdout.String())
}

func (suite *DepUpdateTestSuite) TestPrepareChartRequired() {
	d := DepUpdate{}

	suite.mockCmd.EXPECT().Stdout(gomock.Any()).AnyTimes()
	suite.mockCmd.EXPECT().Stderr(gomock.Any()).AnyTimes()

	err := d.Prepare(Config{})
	suite.EqualError(err, "chart is required")
}

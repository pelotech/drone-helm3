package run

import (
  "errors"
  "github.com/golang/mock/gomock"
  "github.com/pelotech/drone-helm3/internal/env"
  "github.com/stretchr/testify/suite"
  "strings"
  "testing"
)

type DepActionTestSuite struct {
  suite.Suite
  ctrl            *gomock.Controller
  mockCmd         *Mockcmd
  originalCommand func(string, ...string) cmd
}

func (suite *DepActionTestSuite) BeforeTest(_, _ string) {
  suite.ctrl = gomock.NewController(suite.T())
  suite.mockCmd = NewMockcmd(suite.ctrl)

  suite.originalCommand = command
  command = func(path string, args ...string) cmd { return suite.mockCmd }
}

func (suite *DepActionTestSuite) AfterTest(_, _ string) {
  command = suite.originalCommand
}

func TestDepActionTestSuite(t *testing.T) {
  suite.Run(t, new(DepActionTestSuite))
}

func (suite *DepActionTestSuite) TestNewDepAction() {
  cfg := env.Config{
    Chart: "scatterplot",
  }
  d := NewDepAction(cfg)
  suite.Equal("scatterplot", d.chart)
}

func (suite *DepActionTestSuite) TestPrepareAndExecuteBuild() {
  defer suite.ctrl.Finish()

  stdout := strings.Builder{}
  stderr := strings.Builder{}
  cfg := env.Config{
    Chart:              "your_top_songs_2019",
    Stdout:             &stdout,
    Stderr:             &stderr,
    DependenciesAction: "build",
  }

  command = func(path string, args ...string) cmd {
    suite.Equal(helmBin, path)
    suite.Equal([]string{"dependency", "build", "your_top_songs_2019"}, args)

    return suite.mockCmd
  }
  suite.mockCmd.EXPECT().
    Stdout(&stdout)
  suite.mockCmd.EXPECT().
    Stderr(&stderr)
  suite.mockCmd.EXPECT().
    Run().
    Times(1)

  d := NewDepAction(cfg)

  suite.Require().NoError(d.Prepare())
  suite.NoError(d.Execute())
}

func (suite *DepActionTestSuite) TestPrepareAndExecuteUpdate() {
  defer suite.ctrl.Finish()

  stdout := strings.Builder{}
  stderr := strings.Builder{}
  cfg := env.Config{
    Chart:              "your_top_songs_2019",
    Stdout:             &stdout,
    Stderr:             &stderr,
    DependenciesAction: "update",
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

  d := NewDepAction(cfg)

  suite.Require().NoError(d.Prepare())
  suite.NoError(d.Execute())
}

func (suite *DepActionTestSuite) TestPrepareAndExecuteUnknown() {
  defer suite.ctrl.Finish()

  stdout := strings.Builder{}
  stderr := strings.Builder{}
  cfg := env.Config{
    Chart:              "your_top_songs_2019",
    Stdout:             &stdout,
    Stderr:             &stderr,
    DependenciesAction: "downgrade",
  }

  d := NewDepAction(cfg)
  suite.Require().Equal(errors.New("unknown dependency_action: downgrade"), d.Prepare())
}

func (suite *DepActionTestSuite) TestPrepareChartRequired() {
  d := NewDepAction(env.Config{})

  suite.mockCmd.EXPECT().Stdout(gomock.Any()).AnyTimes()
  suite.mockCmd.EXPECT().Stderr(gomock.Any()).AnyTimes()

  err := d.Prepare()
  suite.EqualError(err, "chart is required")
}

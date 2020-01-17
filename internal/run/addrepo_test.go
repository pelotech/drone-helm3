package run

import (
	"fmt"
	"github.com/golang/mock/gomock"
	"github.com/pelotech/drone-helm3/internal/env"
	"github.com/stretchr/testify/suite"
	"strings"
	"testing"
)

type AddRepoTestSuite struct {
	suite.Suite
	ctrl            *gomock.Controller
	mockCmd         *Mockcmd
	originalCommand func(string, ...string) cmd
	commandPath     string
	commandArgs     []string
}

func (suite *AddRepoTestSuite) BeforeTest(_, _ string) {
	suite.ctrl = gomock.NewController(suite.T())
	suite.mockCmd = NewMockcmd(suite.ctrl)

	suite.originalCommand = command
	command = func(path string, args ...string) cmd {
		suite.commandPath = path
		suite.commandArgs = args
		return suite.mockCmd
	}
}

func (suite *AddRepoTestSuite) AfterTest(_, _ string) {
	suite.ctrl.Finish()
	command = suite.originalCommand
}

func TestAddRepoTestSuite(t *testing.T) {
	suite.Run(t, new(AddRepoTestSuite))
}

func (suite *AddRepoTestSuite) TestNewAddRepo() {
	repo := NewAddRepo(env.Config{}, "picompress=https://github.com/caleb_phipps/picompress")
	suite.Require().NotNil(repo)
	suite.Equal("picompress=https://github.com/caleb_phipps/picompress", repo.repo)
	suite.NotNil(repo.config)
}

func (suite *AddRepoTestSuite) TestPrepareAndExecute() {
	stdout := strings.Builder{}
	stderr := strings.Builder{}
	cfg := env.Config{
		Stdout: &stdout,
		Stderr: &stderr,
	}
	a := NewAddRepo(cfg, "edeath=https://github.com/n_marks/e-death")

	suite.mockCmd.EXPECT().
		Stdout(&stdout).
		Times(1)
	suite.mockCmd.EXPECT().
		Stderr(&stderr).
		Times(1)

	suite.Require().NoError(a.Prepare())
	suite.Equal(helmBin, suite.commandPath)
	suite.Equal([]string{"repo", "add", "edeath", "https://github.com/n_marks/e-death"}, suite.commandArgs)

	suite.mockCmd.EXPECT().
		Run().
		Times(1)

	suite.Require().NoError(a.Execute())

}

func (suite *AddRepoTestSuite) TestPrepareRepoIsRequired() {
	// These aren't really expected, but allowing them gives clearer test-failure messages
	suite.mockCmd.EXPECT().Stdout(gomock.Any()).AnyTimes()
	suite.mockCmd.EXPECT().Stderr(gomock.Any()).AnyTimes()
	a := NewAddRepo(env.Config{}, "")

	err := a.Prepare()
	suite.EqualError(err, "repo is required")
}

func (suite *AddRepoTestSuite) TestPrepareMalformedRepo() {
	a := NewAddRepo(env.Config{}, "dwim")
	err := a.Prepare()
	suite.EqualError(err, "bad repo spec 'dwim'")
}

func (suite *AddRepoTestSuite) TestPrepareWithEqualSignInURL() {
	suite.mockCmd.EXPECT().Stdout(gomock.Any()).AnyTimes()
	suite.mockCmd.EXPECT().Stderr(gomock.Any()).AnyTimes()
	a := NewAddRepo(env.Config{}, "samaritan=https://github.com/arthur_claypool/samaritan?version=2.1")
	suite.NoError(a.Prepare())
	suite.Contains(suite.commandArgs, "https://github.com/arthur_claypool/samaritan?version=2.1")
}

func (suite *AddRepoTestSuite) TestNamespaceFlag() {
	suite.mockCmd.EXPECT().Stdout(gomock.Any()).AnyTimes()
	suite.mockCmd.EXPECT().Stderr(gomock.Any()).AnyTimes()
	cfg := env.Config{
		Namespace: "alliteration",
	}
	a := NewAddRepo(cfg, "edeath=https://github.com/theater_guy/e-death")

	suite.NoError(a.Prepare())
	suite.Equal(suite.commandPath, helmBin)
	suite.Equal(suite.commandArgs, []string{"--namespace", "alliteration",
		"repo", "add", "edeath", "https://github.com/theater_guy/e-death"})
}

func (suite *AddRepoTestSuite) TestDebugFlag() {
	suite.mockCmd.EXPECT().Stdout(gomock.Any()).AnyTimes()
	suite.mockCmd.EXPECT().Stderr(gomock.Any()).AnyTimes()

	stderr := strings.Builder{}

	command = func(path string, args ...string) cmd {
		suite.mockCmd.EXPECT().
			String().
			Return(fmt.Sprintf("%s %s", path, strings.Join(args, " ")))

		return suite.mockCmd
	}

	cfg := env.Config{
		Debug:  true,
		Stderr: &stderr,
	}
	a := NewAddRepo(cfg, "edeath=https://github.com/the_bug/e-death")

	suite.Require().NoError(a.Prepare())
	suite.Equal(fmt.Sprintf("Generated command: '%s --debug "+
		"repo add edeath https://github.com/the_bug/e-death'\n", helmBin), stderr.String())
}

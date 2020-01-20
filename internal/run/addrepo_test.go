package run

import (
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
	suite.NotNil(repo.certs)
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

func (suite *AddRepoTestSuite) TestRepoAddFlags() {
	suite.mockCmd.EXPECT().Stdout(gomock.Any()).AnyTimes()
	suite.mockCmd.EXPECT().Stderr(gomock.Any()).AnyTimes()
	cfg := env.Config{}
	a := NewAddRepo(cfg, "machine=https://github.com/harold_finch/themachine")

	// inject a ca cert filename so repoCerts won't create any files that we'd have to clean up
	a.certs.caCertFilename = "./helm/reporepo.cert"
	suite.NoError(a.Prepare())
	suite.Equal([]string{"repo", "add", "--ca-file", "./helm/reporepo.cert",
		"machine", "https://github.com/harold_finch/themachine"}, suite.commandArgs)
}

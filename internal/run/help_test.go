package run

import (
	"fmt"
	"github.com/golang/mock/gomock"
	"github.com/pelotech/drone-helm3/internal/env"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
	"strings"
	"testing"
)

type HelpTestSuite struct {
	suite.Suite
}

func TestHelpTestSuite(t *testing.T) {
	suite.Run(t, new(HelpTestSuite))
}

func (suite *HelpTestSuite) TestNewHelp() {
	cfg := env.Config{
		Command: "everybody dance NOW!!",
	}
	help := NewHelp(cfg)
	suite.Require().NotNil(help)
	suite.Equal("everybody dance NOW!!", help.helmCommand)
}

func (suite *HelpTestSuite) TestPrepare() {
	ctrl := gomock.NewController(suite.T())
	defer ctrl.Finish()

	mCmd := NewMockcmd(ctrl)
	originalCommand := command

	command = func(path string, args ...string) cmd {
		assert.Equal(suite.T(), helmBin, path)
		assert.Equal(suite.T(), []string{"help"}, args)
		return mCmd
	}
	defer func() { command = originalCommand }()

	stdout := strings.Builder{}
	stderr := strings.Builder{}

	mCmd.EXPECT().
		Stdout(&stdout)
	mCmd.EXPECT().
		Stderr(&stderr)

	cfg := env.Config{
		Stdout: &stdout,
		Stderr: &stderr,
	}

	h := NewHelp(cfg)
	err := h.Prepare()
	suite.NoError(err)
}

func (suite *HelpTestSuite) TestExecute() {
	ctrl := gomock.NewController(suite.T())
	defer ctrl.Finish()
	mCmd := NewMockcmd(ctrl)

	mCmd.EXPECT().
		Run().
		Times(2)

	help := NewHelp(env.Config{Command: "help"})
	help.cmd = mCmd
	suite.NoError(help.Execute())

	help.helmCommand = "get down on friday"
	suite.EqualError(help.Execute(), "unknown command 'get down on friday'")
}

func (suite *HelpTestSuite) TestPrepareDebugFlag() {
	stdout := strings.Builder{}
	stderr := strings.Builder{}
	cfg := env.Config{
		Debug:  true,
		Stdout: &stdout,
		Stderr: &stderr,
	}

	help := NewHelp(cfg)
	help.Prepare()

	want := fmt.Sprintf("Generated command: '%s --debug help'\n", helmBin)
	suite.Equal(want, stderr.String())
	suite.Equal("", stdout.String())
}

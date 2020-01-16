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
	suite.Equal("everybody dance NOW!!", help.HelmCommand)
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

	cfg := Config{
		Stdout: &stdout,
		Stderr: &stderr,
	}

	h := Help{}
	err := h.Prepare(cfg)
	suite.NoError(err)
}

func (suite *HelpTestSuite) TestExecute() {
	ctrl := gomock.NewController(suite.T())
	defer ctrl.Finish()
	mCmd := NewMockcmd(ctrl)
	originalCommand := command
	command = func(_ string, _ ...string) cmd {
		return mCmd
	}
	defer func() { command = originalCommand }()

	mCmd.EXPECT().
		Run().
		Times(2)

	cfg := Config{}
	help := Help{
		HelmCommand: "help",
		cmd:         mCmd,
	}
	suite.NoError(help.Execute(cfg))

	help.HelmCommand = "get down on friday"
	suite.EqualError(help.Execute(cfg), "unknown command 'get down on friday'")
}

func (suite *HelpTestSuite) TestPrepareDebugFlag() {
	help := Help{}

	stdout := strings.Builder{}
	stderr := strings.Builder{}
	cfg := Config{
		Debug:  true,
		Stdout: &stdout,
		Stderr: &stderr,
	}

	help.Prepare(cfg)

	want := fmt.Sprintf("Generated command: '%s --debug help'\n", helmBin)
	suite.Equal(want, stderr.String())
	suite.Equal("", stdout.String())
}

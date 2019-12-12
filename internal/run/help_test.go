package run

import (
	"fmt"
	"github.com/golang/mock/gomock"
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
	mCmd.EXPECT().
		Run().
		Times(1)

	cfg := Config{
		Stdout: &stdout,
		Stderr: &stderr,
	}

	h := Help{}
	err := h.Prepare(cfg)
	suite.Require().Nil(err)
	h.Execute(cfg)
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

package run

import (
	"fmt"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/suite"
	"strings"
	"testing"
)

type UpgradeTestSuite struct {
	suite.Suite
	ctrl            *gomock.Controller
	mockCmd         *Mockcmd
	originalCommand func(string, ...string) cmd
}

func (suite *UpgradeTestSuite) BeforeTest(_, _ string) {
	suite.ctrl = gomock.NewController(suite.T())
	suite.mockCmd = NewMockcmd(suite.ctrl)

	suite.originalCommand = command
	command = func(path string, args ...string) cmd { return suite.mockCmd }
}

func (suite *UpgradeTestSuite) AfterTest(_, _ string) {
	command = suite.originalCommand
}

func TestUpgradeTestSuite(t *testing.T) {
	suite.Run(t, new(UpgradeTestSuite))
}

func (suite *UpgradeTestSuite) TestPrepare() {
	defer suite.ctrl.Finish()

	u := Upgrade{
		Chart:   "at40",
		Release: "jonas_brothers_only_human",
	}

	command = func(path string, args ...string) cmd {
		suite.Equal(helmBin, path)
		suite.Equal([]string{"upgrade", "--install", "jonas_brothers_only_human", "at40"}, args)

		return suite.mockCmd
	}

	suite.mockCmd.EXPECT().
		Stdout(gomock.Any())
	suite.mockCmd.EXPECT().
		Stderr(gomock.Any())
	suite.mockCmd.EXPECT().
		Run().
		Times(1)

	err := u.Prepare(Config{})
	suite.Require().Nil(err)
	u.Execute()
}

func (suite *UpgradeTestSuite) TestPrepareDebugFlag() {
	u := Upgrade{
		Chart:   "at40",
		Release: "lewis_capaldi_someone_you_loved",
	}

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

	suite.mockCmd.EXPECT().
		Stdout(&stdout)
	suite.mockCmd.EXPECT().
		Stderr(&stderr)

	u.Prepare(cfg)

	want := fmt.Sprintf("Generated command: '%s --debug upgrade --install lewis_capaldi_someone_you_loved at40'\n", helmBin)
	suite.Equal(want, stderr.String())
	suite.Equal("", stdout.String())
}

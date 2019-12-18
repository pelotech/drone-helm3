package run

import (
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/suite"
	"testing"
)

type LintTestSuite struct {
	suite.Suite
	ctrl            *gomock.Controller
	mockCmd         *Mockcmd
	originalCommand func(string, ...string) cmd
}

func (suite *LintTestSuite) BeforeTest(_, _ string) {
	suite.ctrl = gomock.NewController(suite.T())
	suite.mockCmd = NewMockcmd(suite.ctrl)

	suite.originalCommand = command
	command = func(path string, args ...string) cmd { return suite.mockCmd }
}

func (suite *LintTestSuite) AfterTest(_, _ string) {
	command = suite.originalCommand
}

func TestLintTestSuite(t *testing.T) {
	suite.Run(t, new(LintTestSuite))
}

func (suite *LintTestSuite) TestPrepareAndExecute() {
	defer suite.ctrl.Finish()

	l := Lint{
		Chart: "./epic/mychart",
	}

	command = func(path string, args ...string) cmd {
		suite.Equal(helmBin, path)
		suite.Equal([]string{"lint", "./epic/mychart"}, args)

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
	err := l.Prepare(cfg)
	suite.Require().Nil(err)
	l.Execute(cfg)
}

func (suite *LintTestSuite) TestPrepareWithLintFlags() {
	defer suite.ctrl.Finish()

	cfg := Config{
		Values:       "width=5",
		StringValues: "version=2.0",
		ValuesFiles:  []string{"/usr/local/underrides", "/usr/local/overrides"},
	}

	l := Lint{
		Chart: "./uk/top_40",
	}

	command = func(path string, args ...string) cmd {
		suite.Equal(helmBin, path)
		suite.Equal([]string{"lint",
			"--set", "width=5",
			"--set-string", "version=2.0",
			"--values", "/usr/local/underrides",
			"--values", "/usr/local/overrides",
			"./uk/top_40"}, args)

		return suite.mockCmd
	}

	suite.mockCmd.EXPECT().Stdout(gomock.Any())
	suite.mockCmd.EXPECT().Stderr(gomock.Any())

	err := l.Prepare(cfg)
	suite.Require().Nil(err)
}

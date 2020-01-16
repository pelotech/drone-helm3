package run

import (
	"fmt"
	"github.com/golang/mock/gomock"
	"github.com/pelotech/drone-helm3/internal/env"
	"github.com/stretchr/testify/suite"
	"strings"
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

func (suite *LintTestSuite) TestNewLint() {
	cfg := env.Config{
		Chart:        "./flow",
		Values:       "steadfastness,forthrightness",
		StringValues: "tensile_strength,flexibility",
		ValuesFiles:  []string{"/root/price_inventory.yml"},
		LintStrictly: true,
	}
	lint := NewLint(cfg)
	suite.Require().NotNil(lint)
	suite.Equal(&Lint{
		Chart:        "./flow",
		Values:       "steadfastness,forthrightness",
		StringValues: "tensile_strength,flexibility",
		ValuesFiles:  []string{"/root/price_inventory.yml"},
		Strict:       true,
	}, lint)
}

func (suite *LintTestSuite) TestPrepareAndExecute() {
	defer suite.ctrl.Finish()

	stdout := strings.Builder{}
	stderr := strings.Builder{}

	l := Lint{
		Chart: "./epic/mychart",
	}
	cfg := Config{
		Stdout: &stdout,
		Stderr: &stderr,
	}

	command = func(path string, args ...string) cmd {
		suite.Equal(helmBin, path)
		suite.Equal([]string{"lint", "./epic/mychart"}, args)

		return suite.mockCmd
	}

	suite.mockCmd.EXPECT().
		Stdout(&stdout)
	suite.mockCmd.EXPECT().
		Stderr(&stderr)
	suite.mockCmd.EXPECT().
		Run().
		Times(1)

	err := l.Prepare(cfg)
	suite.Require().Nil(err)
	l.Execute()
}

func (suite *LintTestSuite) TestPrepareRequiresChart() {
	// These aren't really expected, but allowing them gives clearer test-failure messages
	suite.mockCmd.EXPECT().Stdout(gomock.Any()).AnyTimes()
	suite.mockCmd.EXPECT().Stderr(gomock.Any()).AnyTimes()

	cfg := Config{}
	l := Lint{}

	err := l.Prepare(cfg)
	suite.EqualError(err, "chart is required", "Chart should be mandatory")
}

func (suite *LintTestSuite) TestPrepareWithLintFlags() {
	defer suite.ctrl.Finish()

	cfg := Config{}

	l := Lint{
		Chart:        "./uk/top_40",
		Values:       "width=5",
		StringValues: "version=2.0",
		ValuesFiles:  []string{"/usr/local/underrides", "/usr/local/overrides"},
		Strict:       true,
	}

	command = func(path string, args ...string) cmd {
		suite.Equal(helmBin, path)
		suite.Equal([]string{"lint",
			"--set", "width=5",
			"--set-string", "version=2.0",
			"--values", "/usr/local/underrides",
			"--values", "/usr/local/overrides",
			"--strict",
			"./uk/top_40"}, args)

		return suite.mockCmd
	}

	suite.mockCmd.EXPECT().Stdout(gomock.Any()).AnyTimes()
	suite.mockCmd.EXPECT().Stderr(gomock.Any()).AnyTimes()

	err := l.Prepare(cfg)
	suite.Require().Nil(err)
}

func (suite *LintTestSuite) TestPrepareWithDebugFlag() {
	defer suite.ctrl.Finish()

	stderr := strings.Builder{}

	cfg := Config{
		Debug:  true,
		Stderr: &stderr,
	}

	l := Lint{
		Chart: "./scotland/top_40",
	}

	command = func(path string, args ...string) cmd {
		suite.mockCmd.EXPECT().
			String().
			Return(fmt.Sprintf("%s %s", path, strings.Join(args, " ")))

		return suite.mockCmd
	}

	suite.mockCmd.EXPECT().Stdout(gomock.Any())
	suite.mockCmd.EXPECT().Stderr(&stderr)

	err := l.Prepare(cfg)
	suite.Require().Nil(err)

	want := fmt.Sprintf("Generated command: '%s --debug lint ./scotland/top_40'\n", helmBin)
	suite.Equal(want, stderr.String())
}

func (suite *LintTestSuite) TestPrepareWithNamespaceFlag() {
	defer suite.ctrl.Finish()

	cfg := Config{
		Namespace: "table-service",
	}

	l := Lint{
		Chart: "./wales/top_40",
	}

	actual := []string{}
	command = func(path string, args ...string) cmd {
		actual = args
		return suite.mockCmd
	}

	suite.mockCmd.EXPECT().Stdout(gomock.Any()).AnyTimes()
	suite.mockCmd.EXPECT().Stderr(gomock.Any()).AnyTimes()

	err := l.Prepare(cfg)
	suite.Require().Nil(err)

	expected := []string{"--namespace", "table-service", "lint", "./wales/top_40"}
	suite.Equal(expected, actual)
}

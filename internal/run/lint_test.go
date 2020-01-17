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
	suite.Equal("./flow", lint.chart)
	suite.Equal("steadfastness,forthrightness", lint.values)
	suite.Equal("tensile_strength,flexibility", lint.stringValues)
	suite.Equal([]string{"/root/price_inventory.yml"}, lint.valuesFiles)
	suite.Equal(true, lint.strict)
	suite.NotNil(lint.config)
}

func (suite *LintTestSuite) TestPrepareAndExecute() {
	defer suite.ctrl.Finish()

	stdout := strings.Builder{}
	stderr := strings.Builder{}

	cfg := env.Config{
		Chart:  "./epic/mychart",
		Stdout: &stdout,
		Stderr: &stderr,
	}
	l := NewLint(cfg)

	command = func(path string, args ...string) cmd {
		suite.Equal(helmBin, path)
		suite.Equal([]string{"lint", "./epic/mychart"}, args)

		return suite.mockCmd
	}

	suite.mockCmd.EXPECT().String().AnyTimes()
	suite.mockCmd.EXPECT().
		Stdout(&stdout)
	suite.mockCmd.EXPECT().
		Stderr(&stderr)
	suite.mockCmd.EXPECT().
		Run().
		Times(1)

	err := l.Prepare()
	suite.Require().Nil(err)
	l.Execute()
}

func (suite *LintTestSuite) TestPrepareRequiresChart() {
	// These aren't really expected, but allowing them gives clearer test-failure messages
	suite.mockCmd.EXPECT().Stdout(gomock.Any()).AnyTimes()
	suite.mockCmd.EXPECT().Stderr(gomock.Any()).AnyTimes()

	l := NewLint(env.Config{})
	err := l.Prepare()
	suite.EqualError(err, "chart is required", "Chart should be mandatory")
}

func (suite *LintTestSuite) TestPrepareWithLintFlags() {
	defer suite.ctrl.Finish()

	cfg := env.Config{
		Chart:        "./uk/top_40",
		Values:       "width=5",
		StringValues: "version=2.0",
		ValuesFiles:  []string{"/usr/local/underrides", "/usr/local/overrides"},
		LintStrictly: true,
	}
	l := NewLint(cfg)

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
	suite.mockCmd.EXPECT().String().AnyTimes()

	err := l.Prepare()
	suite.Require().Nil(err)
}

func (suite *LintTestSuite) TestPrepareWithDebugFlag() {
	defer suite.ctrl.Finish()

	stderr := strings.Builder{}

	cfg := env.Config{
		Debug:  true,
		Stderr: &stderr,
		Chart:  "./scotland/top_40",
	}
	l := NewLint(cfg)

	command = func(path string, args ...string) cmd {
		suite.mockCmd.EXPECT().
			String().
			Return(fmt.Sprintf("%s %s", path, strings.Join(args, " ")))

		return suite.mockCmd
	}

	suite.mockCmd.EXPECT().Stdout(gomock.Any())
	suite.mockCmd.EXPECT().Stderr(&stderr)

	err := l.Prepare()
	suite.Require().Nil(err)

	want := fmt.Sprintf("Generated command: '%s --debug lint ./scotland/top_40'\n", helmBin)
	suite.Equal(want, stderr.String())
}

func (suite *LintTestSuite) TestPrepareWithNamespaceFlag() {
	defer suite.ctrl.Finish()

	cfg := env.Config{
		Namespace: "table-service",
		Chart:     "./wales/top_40",
	}
	l := NewLint(cfg)

	actual := []string{}
	command = func(path string, args ...string) cmd {
		actual = args
		return suite.mockCmd
	}

	suite.mockCmd.EXPECT().Stdout(gomock.Any()).AnyTimes()
	suite.mockCmd.EXPECT().Stderr(gomock.Any()).AnyTimes()

	err := l.Prepare()
	suite.Require().Nil(err)

	expected := []string{"--namespace", "table-service", "lint", "./wales/top_40"}
	suite.Equal(expected, actual)
}

package run

import (
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/pelotech/drone-helm3/internal/env"
	"github.com/stretchr/testify/suite"
)

type ChartTestSuite struct {
	suite.Suite
	ctrl            *gomock.Controller
	mockCmd         *Mockcmd
	actualArgs      []string
	originalCommand func(string, ...string) cmd
}

func (suite *ChartTestSuite) BeforeTest(_, _ string) {
	suite.ctrl = gomock.NewController(suite.T())
	suite.mockCmd = NewMockcmd(suite.ctrl)

	suite.originalCommand = command
	command = func(path string, args ...string) cmd {
		suite.actualArgs = args
		return suite.mockCmd
	}
}

func (suite *ChartTestSuite) AfterTest(_, _ string) {
	command = suite.originalCommand
}

func TestChartTestSuite(t *testing.T) {
	suite.Run(t, new(ChartTestSuite))
}

func (suite *ChartTestSuite) TestNewChart() {
	cfg := env.Config{
		Chart:            "./",
		RegistryRepoName: "repo_name",
		RegistryURL:      "registry_url",
		ChartVersion:     "0.0.1",
	}
	cs := NewChart("save", cfg)
	cp := NewChart("push", cfg)

	suite.Equal("./", cs.chartPath)
	suite.Equal("repo_name", cs.registryRepoName)
	suite.Equal("registry_url", cs.registryURL)
	suite.Equal("repo_name", cs.registryRepoName)
	suite.Equal("save", cs.subCommand)
	suite.Equal("push", cp.subCommand)
	suite.NotNil(cs.config)
}

func (suite *ChartTestSuite) TestPrepareAndExecuteSave() {
	defer suite.ctrl.Finish()

	cfg := env.Config{
		Chart:            "./",
		RegistryRepoName: "repo_name",
		RegistryURL:      "registry_url",
		ChartVersion:     "0.0.1",
	}

	c := NewChart("save", cfg)

	actual := []string{}
	command = func(path string, args ...string) cmd {
		suite.Equal(helmBin, path)
		actual = args

		return suite.mockCmd
	}

	suite.mockCmd.EXPECT().
		Stdout(gomock.Any())
	suite.mockCmd.EXPECT().
		Stderr(gomock.Any())
	suite.mockCmd.EXPECT().
		Run().
		Times(1)

	suite.NoError(c.Prepare())
	expected := []string{"chart", "save", "./", "registry_url/repo_name:0.0.1"}
	suite.Equal(expected, actual)

	c.Execute()
}

func (suite *ChartTestSuite) TestPrepareAndExecutePush() {
	defer suite.ctrl.Finish()

	cfg := env.Config{
		RegistryRepoName: "repo_name",
		RegistryURL:      "registry_url",
		ChartVersion:     "0.0.1",
	}

	c := NewChart("push", cfg)

	actual := []string{}
	command = func(path string, args ...string) cmd {
		suite.Equal(helmBin, path)
		actual = args

		return suite.mockCmd
	}

	suite.mockCmd.EXPECT().
		Stdout(gomock.Any())
	suite.mockCmd.EXPECT().
		Stderr(gomock.Any())
	suite.mockCmd.EXPECT().
		Run().
		Times(1)

	suite.NoError(c.Prepare())
	expected := []string{"chart", "push", "registry_url/repo_name:0.0.1"}
	suite.Equal(expected, actual)

	c.Execute()
}

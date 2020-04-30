package run

import (
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/pelotech/drone-helm3/internal/env"
	"github.com/stretchr/testify/suite"
)

type RegistryTestSuite struct {
	suite.Suite
	ctrl            *gomock.Controller
	mockCmd         *Mockcmd
	actualArgs      []string
	originalCommand func(string, ...string) cmd
}

func (suite *RegistryTestSuite) BeforeTest(_, _ string) {
	suite.ctrl = gomock.NewController(suite.T())
	suite.mockCmd = NewMockcmd(suite.ctrl)

	suite.originalCommand = command
	command = func(path string, args ...string) cmd {
		suite.actualArgs = args
		return suite.mockCmd
	}
}

func (suite *RegistryTestSuite) AfterTest(_, _ string) {
	command = suite.originalCommand
}

func TestRegistryTestSuite(t *testing.T) {
	suite.Run(t, new(RegistryTestSuite))
}

func (suite *RegistryTestSuite) TestNewRegistry() {
	cfg := env.Config{
		RegistryLoginUserID:   "johndoe",
		RegistryLoginPassword: "super_secret_password",
		RegistryURL:           "registry_url",
	}
	r := NewRegistry("login", cfg)
	ro := NewRegistry("logout", cfg)

	suite.Equal("johndoe", r.userID)
	suite.Equal("super_secret_password", r.password)
	suite.Equal("registry_url", r.registryURL)
	suite.Equal("login", r.subCommand)
	suite.Equal("logout", ro.subCommand)
	suite.NotNil(r.config)
}

func (suite *RegistryTestSuite) TestPrepareAndExecuteLogin() {
	defer suite.ctrl.Finish()

	cfg := env.Config{
		RegistryLoginUserID:   "johndoe",
		RegistryLoginPassword: "super_secret_password",
		RegistryURL:           "registry_url",
	}

	u := NewRegistry("login", cfg)

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

	suite.NoError(u.Prepare())
	expected := []string{"registry", "login", "-u", "johndoe", "-p", "super_secret_password", "registry_url"}
	suite.Equal(expected, actual)

	u.Execute()
}

func (suite *RegistryTestSuite) TestPrepareAndExecuteLogout() {
	defer suite.ctrl.Finish()

	cfg := env.Config{
		RegistryURL: "registry_url",
	}

	u := NewRegistry("logout", cfg)

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

	suite.NoError(u.Prepare())
	expected := []string{"registry", "logout", "registry_url"}
	suite.Equal(expected, actual)

	u.Execute()
}

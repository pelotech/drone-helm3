package helm

import (
	"fmt"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/suite"
	"os"
	"testing"

	"github.com/pelotech/drone-helm3/internal/run"
)

type PlanTestSuite struct {
	suite.Suite
}

func TestPlanTestSuite(t *testing.T) {
	suite.Run(t, new(PlanTestSuite))
}

func (suite *PlanTestSuite) TestNewPlan() {
	ctrl := gomock.NewController(suite.T())
	stepOne := NewMockStep(ctrl)
	stepTwo := NewMockStep(ctrl)

	origHelp := help
	help = func(cfg Config) []Step {
		return []Step{stepOne, stepTwo}
	}
	defer func() { help = origHelp }()

	cfg := Config{
		Command:      "help",
		Debug:        false,
		KubeConfig:   "/branch/.sfere/profig",
		Values:       "steadfastness,forthrightness",
		StringValues: "tensile_strength,flexibility",
		ValuesFiles:  []string{"/root/price_inventory.yml"},
		Namespace:    "outer",
	}

	runCfg := run.Config{
		Debug:        false,
		KubeConfig:   "/branch/.sfere/profig",
		Values:       "steadfastness,forthrightness",
		StringValues: "tensile_strength,flexibility",
		ValuesFiles:  []string{"/root/price_inventory.yml"},
		Namespace:    "outer",
		Stdout:       os.Stdout,
		Stderr:       os.Stderr,
	}

	stepOne.EXPECT().
		Prepare(runCfg)
	stepTwo.EXPECT().
		Prepare(runCfg)

	plan, err := NewPlan(cfg)
	suite.Require().Nil(err)
	suite.Equal(cfg, plan.cfg)
	suite.Equal(runCfg, plan.runCfg)
}

func (suite *PlanTestSuite) TestNewPlanAbortsOnError() {
	ctrl := gomock.NewController(suite.T())
	stepOne := NewMockStep(ctrl)
	stepTwo := NewMockStep(ctrl)

	origHelp := help
	help = func(cfg Config) []Step {
		return []Step{stepOne, stepTwo}
	}
	defer func() { help = origHelp }()

	cfg := Config{
		Command: "help",
	}

	stepOne.EXPECT().
		Prepare(gomock.Any()).
		Return(fmt.Errorf("I'm starry Dave, aye, cat blew that"))

	_, err := NewPlan(cfg)
	suite.Require().NotNil(err)
	suite.EqualError(err, "while preparing *helm.MockStep step: I'm starry Dave, aye, cat blew that")
}

func (suite *PlanTestSuite) TestUpgrade() {
	cfg := Config{
		ChartVersion: "seventeen",
		DryRun:       true,
		Wait:         true,
		ReuseValues:  true,
		Timeout:      "go sit in the corner",
		Chart:        "billboard_top_100",
		Release:      "post_malone_circles",
		Force:        true,
	}

	steps := upgrade(cfg)
	suite.Require().Equal(2, len(steps), "upgrade should return 2 steps")
	suite.Require().IsType(&run.InitKube{}, steps[0])

	suite.Require().IsType(&run.Upgrade{}, steps[1])
	upgrade, _ := steps[1].(*run.Upgrade)

	expected := &run.Upgrade{
		Chart:        cfg.Chart,
		Release:      cfg.Release,
		ChartVersion: cfg.ChartVersion,
		DryRun:       true,
		Wait:         cfg.Wait,
		ReuseValues:  cfg.ReuseValues,
		Timeout:      cfg.Timeout,
		Force:        cfg.Force,
	}

	suite.Equal(expected, upgrade)
}

func (suite *PlanTestSuite) TestDel() {
	cfg := Config{
		KubeToken:      "b2YgbXkgYWZmZWN0aW9u",
		SkipTLSVerify:  true,
		Certificate:    "cHJvY2xhaW1zIHdvbmRlcmZ1bCBmcmllbmRzaGlw",
		APIServer:      "98.765.43.21",
		ServiceAccount: "greathelm",
		DryRun:         true,
		Timeout:        "think about what you did",
		Release:        "jetta_id_love_to_change_the_world",
	}

	steps := del(cfg)
	suite.Require().Equal(2, len(steps), "del should return 2 steps")

	suite.Require().IsType(&run.InitKube{}, steps[0])
	init, _ := steps[0].(*run.InitKube)
	var expected Step = &run.InitKube{
		SkipTLSVerify:  true,
		Certificate:    "cHJvY2xhaW1zIHdvbmRlcmZ1bCBmcmllbmRzaGlw",
		APIServer:      "98.765.43.21",
		ServiceAccount: "greathelm",
		Token:          "b2YgbXkgYWZmZWN0aW9u",
		TemplateFile:   kubeConfigTemplate,
	}

	suite.Equal(expected, init)

	suite.Require().IsType(&run.Delete{}, steps[1])
	actual, _ := steps[1].(*run.Delete)
	expected = &run.Delete{
		Release: "jetta_id_love_to_change_the_world",
		DryRun:  true,
	}
	suite.Equal(expected, actual)
}

func (suite *PlanTestSuite) TestInitKube() {
	cfg := Config{
		KubeToken:      "cXVlZXIgY2hhcmFjdGVyCg==",
		SkipTLSVerify:  true,
		Certificate:    "b2Ygd29rZW5lc3MK",
		APIServer:      "123.456.78.9",
		ServiceAccount: "helmet",
	}

	steps := initKube(cfg)
	suite.Require().Equal(1, len(steps), "initKube should return one step")
	suite.Require().IsType(&run.InitKube{}, steps[0])
	init, _ := steps[0].(*run.InitKube)

	expected := &run.InitKube{
		SkipTLSVerify:  true,
		Certificate:    "b2Ygd29rZW5lc3MK",
		APIServer:      "123.456.78.9",
		ServiceAccount: "helmet",
		Token:          "cXVlZXIgY2hhcmFjdGVyCg==",
		TemplateFile:   kubeConfigTemplate,
	}
	suite.Equal(expected, init)
}

func (suite *PlanTestSuite) TestDeterminePlanUpgradeCommand() {
	cfg := Config{
		Command: "upgrade",
	}
	stepsMaker := determineSteps(cfg)
	suite.Same(&upgrade, stepsMaker)
}

func (suite *PlanTestSuite) TestDeterminePlanUpgradeFromDroneEvent() {
	cfg := Config{}

	upgradeEvents := []string{"push", "tag", "deployment", "pull_request", "promote", "rollback"}
	for _, event := range upgradeEvents {
		cfg.DroneEvent = event
		stepsMaker := determineSteps(cfg)
		suite.Same(&upgrade, stepsMaker, fmt.Sprintf("for event type '%s'", event))
	}
}

func (suite *PlanTestSuite) TestDeterminePlanDeleteCommand() {
	cfg := Config{
		Command: "delete",
	}
	stepsMaker := determineSteps(cfg)
	suite.Same(&del, stepsMaker)
}

func (suite *PlanTestSuite) TestDeterminePlanDeleteFromDroneEvent() {
	cfg := Config{
		DroneEvent: "delete",
	}
	stepsMaker := determineSteps(cfg)
	suite.Same(&del, stepsMaker)
}

func (suite *PlanTestSuite) TestDeterminePlanHelpCommand() {
	cfg := Config{
		Command: "help",
	}

	stepsMaker := determineSteps(cfg)
	suite.Same(&help, stepsMaker)
}

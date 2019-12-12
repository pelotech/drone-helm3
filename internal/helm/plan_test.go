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
		KubeToken:      "cXVlZXIgY2hhcmFjdGVyCg==",
		SkipTLSVerify:  true,
		Certificate:    "b2Ygd29rZW5lc3MK",
		APIServer:      "123.456.78.9",
		ServiceAccount: "helmet",
		ChartVersion:   "seventeen",
		Wait:           true,
		ReuseValues:    true,
		Timeout:        "go sit in the corner",
		Chart:          "billboard_top_100",
		Release:        "post_malone_circles",
		Force:          true,
	}

	steps := upgrade(cfg)

	suite.Equal(2, len(steps))

	suite.Require().IsType(&run.InitKube{}, steps[0])
	init, _ := steps[0].(*run.InitKube)

	var expected Step = &run.InitKube{
		SkipTLSVerify:  cfg.SkipTLSVerify,
		Certificate:    cfg.Certificate,
		APIServer:      cfg.APIServer,
		ServiceAccount: cfg.ServiceAccount,
		Token:          cfg.KubeToken,
		TemplateFile:   kubeConfigTemplate,
	}

	suite.Equal(expected, init)

	suite.Require().IsType(&run.Upgrade{}, steps[1])
	upgrade, _ := steps[1].(*run.Upgrade)

	expected = &run.Upgrade{
		Chart:        cfg.Chart,
		Release:      cfg.Release,
		ChartVersion: cfg.ChartVersion,
		Wait:         cfg.Wait,
		ReuseValues:  cfg.ReuseValues,
		Timeout:      cfg.Timeout,
		Force:        cfg.Force,
	}

	suite.Equal(expected, upgrade)
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

func (suite *PlanTestSuite) TestDeterminePlanHelpCommand() {
	cfg := Config{
		Command: "help",
	}

	stepsMaker := determineSteps(cfg)
	suite.Same(&help, stepsMaker)
}

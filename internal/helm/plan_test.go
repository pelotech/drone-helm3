package helm

import (
	"fmt"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/suite"
	"strings"
	"testing"

	"github.com/pelotech/drone-helm3/internal/env"
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
	defer ctrl.Finish()
	stepOne := NewMockStep(ctrl)
	stepTwo := NewMockStep(ctrl)

	origHelp := help
	help = func(cfg env.Config) []Step {
		return []Step{stepOne, stepTwo}
	}
	defer func() { help = origHelp }()

	stdout := strings.Builder{}
	stderr := strings.Builder{}
	cfg := env.Config{
		Command:   "help",
		Debug:     false,
		Namespace: "outer",
		Stdout:    &stdout,
		Stderr:    &stderr,
	}

	runCfg := run.Config{
		Debug:     false,
		Namespace: "outer",
		Stdout:    &stdout,
		Stderr:    &stderr,
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
	defer ctrl.Finish()
	stepOne := NewMockStep(ctrl)
	stepTwo := NewMockStep(ctrl)

	origHelp := help
	help = func(cfg env.Config) []Step {
		return []Step{stepOne, stepTwo}
	}
	defer func() { help = origHelp }()

	cfg := env.Config{
		Command: "help",
	}

	stepOne.EXPECT().
		Prepare(gomock.Any()).
		Return(fmt.Errorf("I'm starry Dave, aye, cat blew that"))

	_, err := NewPlan(cfg)
	suite.Require().NotNil(err)
	suite.EqualError(err, "while preparing *helm.MockStep step: I'm starry Dave, aye, cat blew that")
}

func (suite *PlanTestSuite) TestExecute() {
	ctrl := gomock.NewController(suite.T())
	defer ctrl.Finish()
	stepOne := NewMockStep(ctrl)
	stepTwo := NewMockStep(ctrl)

	runCfg := run.Config{}

	plan := Plan{
		steps:  []Step{stepOne, stepTwo},
		runCfg: runCfg,
	}

	stepOne.EXPECT().
		Execute().
		Times(1)
	stepTwo.EXPECT().
		Execute().
		Times(1)

	suite.NoError(plan.Execute())
}

func (suite *PlanTestSuite) TestExecuteAbortsOnError() {
	ctrl := gomock.NewController(suite.T())
	defer ctrl.Finish()
	stepOne := NewMockStep(ctrl)
	stepTwo := NewMockStep(ctrl)

	runCfg := run.Config{}

	plan := Plan{
		steps:  []Step{stepOne, stepTwo},
		runCfg: runCfg,
	}

	stepOne.EXPECT().
		Execute().
		Times(1).
		Return(fmt.Errorf("oh, he'll gnaw"))

	err := plan.Execute()
	suite.EqualError(err, "while executing *helm.MockStep step: oh, he'll gnaw")
}

func (suite *PlanTestSuite) TestUpgrade() {
	steps := upgrade(env.Config{})
	suite.Require().Equal(2, len(steps), "upgrade should return 2 steps")
	suite.IsType(&run.InitKube{}, steps[0])
	suite.IsType(&run.Upgrade{}, steps[1])
}

func (suite *PlanTestSuite) TestUpgradeWithUpdateDependencies() {
	cfg := env.Config{
		UpdateDependencies: true,
	}
	steps := upgrade(cfg)
	suite.Require().Equal(3, len(steps), "upgrade should have a third step when DepUpdate is true")
	suite.IsType(&run.InitKube{}, steps[0])
	suite.IsType(&run.DepUpdate{}, steps[1])
}

func (suite *PlanTestSuite) TestUpgradeWithAddRepos() {
	cfg := env.Config{
		AddRepos: []string{
			"machine=https://github.com/harold_finch/themachine",
		},
	}
	steps := upgrade(cfg)
	suite.Require().True(len(steps) > 1, "upgrade should generate at least two steps")
	suite.IsType(&run.AddRepo{}, steps[1])
}

func (suite *PlanTestSuite) TestUninstall() {
	steps := uninstall(env.Config{})
	suite.Require().Equal(2, len(steps), "uninstall should return 2 steps")

	suite.IsType(&run.InitKube{}, steps[0])
	suite.IsType(&run.Uninstall{}, steps[1])
}

func (suite *PlanTestSuite) TestUninstallWithUpdateDependencies() {
	cfg := env.Config{
		UpdateDependencies: true,
	}
	steps := uninstall(cfg)
	suite.Require().Equal(3, len(steps), "uninstall should have a third step when DepUpdate is true")
	suite.IsType(&run.InitKube{}, steps[0])
	suite.IsType(&run.DepUpdate{}, steps[1])
}

func (suite *PlanTestSuite) TestLint() {
	steps := lint(env.Config{})
	suite.Require().Equal(1, len(steps))
	suite.IsType(&run.Lint{}, steps[0])
}

func (suite *PlanTestSuite) TestLintWithUpdateDependencies() {
	cfg := env.Config{
		UpdateDependencies: true,
	}
	steps := lint(cfg)
	suite.Require().Equal(2, len(steps), "lint should have a second step when DepUpdate is true")
	suite.IsType(&run.DepUpdate{}, steps[0])
}

func (suite *PlanTestSuite) TestLintWithAddRepos() {
	cfg := env.Config{
		AddRepos: []string{"friendczar=https://github.com/logan_pierce/friendczar"},
	}
	steps := lint(cfg)
	suite.Require().True(len(steps) > 0, "lint should return at least one step")
	suite.IsType(&run.AddRepo{}, steps[0])
}

func (suite *PlanTestSuite) TestDeterminePlanUpgradeCommand() {
	cfg := env.Config{
		Command: "upgrade",
	}
	stepsMaker := determineSteps(cfg)
	suite.Same(&upgrade, stepsMaker)
}

func (suite *PlanTestSuite) TestDeterminePlanUpgradeFromDroneEvent() {
	cfg := env.Config{}

	upgradeEvents := []string{"push", "tag", "deployment", "pull_request", "promote", "rollback"}
	for _, event := range upgradeEvents {
		cfg.DroneEvent = event
		stepsMaker := determineSteps(cfg)
		suite.Same(&upgrade, stepsMaker, fmt.Sprintf("for event type '%s'", event))
	}
}

func (suite *PlanTestSuite) TestDeterminePlanUninstallCommand() {
	cfg := env.Config{
		Command: "uninstall",
	}
	stepsMaker := determineSteps(cfg)
	suite.Same(&uninstall, stepsMaker)
}

// helm_command = delete is provided as an alias for backward-compatibility with drone-helm
func (suite *PlanTestSuite) TestDeterminePlanDeleteCommand() {
	cfg := env.Config{
		Command: "delete",
	}
	stepsMaker := determineSteps(cfg)
	suite.Same(&uninstall, stepsMaker)
}

func (suite *PlanTestSuite) TestDeterminePlanDeleteFromDroneEvent() {
	cfg := env.Config{
		DroneEvent: "delete",
	}
	stepsMaker := determineSteps(cfg)
	suite.Same(&uninstall, stepsMaker)
}

func (suite *PlanTestSuite) TestDeterminePlanLintCommand() {
	cfg := env.Config{
		Command: "lint",
	}

	stepsMaker := determineSteps(cfg)
	suite.Same(&lint, stepsMaker)
}

func (suite *PlanTestSuite) TestDeterminePlanHelpCommand() {
	cfg := env.Config{
		Command: "help",
	}

	stepsMaker := determineSteps(cfg)
	suite.Same(&help, stepsMaker)
}

package helm

import (
	"fmt"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/suite"
	"strings"
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
	defer ctrl.Finish()
	stepOne := NewMockStep(ctrl)
	stepTwo := NewMockStep(ctrl)

	origHelp := help
	help = func(cfg Config) []Step {
		return []Step{stepOne, stepTwo}
	}
	defer func() { help = origHelp }()

	stdout := strings.Builder{}
	stderr := strings.Builder{}
	cfg := Config{
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
		Execute(runCfg).
		Times(1)
	stepTwo.EXPECT().
		Execute(runCfg).
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
		Execute(runCfg).
		Times(1).
		Return(fmt.Errorf("oh, he'll gnaw"))

	err := plan.Execute()
	suite.EqualError(err, "while executing *helm.MockStep step: oh, he'll gnaw")
}

func (suite *PlanTestSuite) TestUpgrade() {
	cfg := Config{
		ChartVersion:  "seventeen",
		DryRun:        true,
		Wait:          true,
		Values:        "steadfastness,forthrightness",
		StringValues:  "tensile_strength,flexibility",
		ValuesFiles:   []string{"/root/price_inventory.yml"},
		ReuseValues:   true,
		Timeout:       "go sit in the corner",
		Chart:         "billboard_top_100",
		Release:       "post_malone_circles",
		Force:         true,
		AtomicUpgrade: true,
		CleanupOnFail: true,
		RepoCAFile:    "state_licensure.repo.cert",
	}

	steps := upgrade(cfg)
	suite.Require().Equal(2, len(steps), "upgrade should return 2 steps")
	suite.Require().IsType(&run.InitKube{}, steps[0])

	suite.Require().IsType(&run.Upgrade{}, steps[1])
	upgrade, _ := steps[1].(*run.Upgrade)

	expected := &run.Upgrade{
		Chart:         cfg.Chart,
		Release:       cfg.Release,
		ChartVersion:  cfg.ChartVersion,
		DryRun:        true,
		Wait:          cfg.Wait,
		Values:        "steadfastness,forthrightness",
		StringValues:  "tensile_strength,flexibility",
		ValuesFiles:   []string{"/root/price_inventory.yml"},
		ReuseValues:   cfg.ReuseValues,
		Timeout:       cfg.Timeout,
		Force:         cfg.Force,
		Atomic:        true,
		CleanupOnFail: true,
		CAFile:        "state_licensure.repo.cert",
	}

	suite.Equal(expected, upgrade)
}

func (suite *PlanTestSuite) TestUpgradeWithUpdateDependencies() {
	cfg := Config{
		UpdateDependencies: true,
	}
	steps := upgrade(cfg)
	suite.Require().Equal(3, len(steps), "upgrade should have a third step when DepUpdate is true")
	suite.IsType(&run.InitKube{}, steps[0])
	suite.IsType(&run.DepUpdate{}, steps[1])
}

func (suite *PlanTestSuite) TestUpgradeWithAddRepos() {
	cfg := Config{
		AddRepos: []string{
			"machine=https://github.com/harold_finch/themachine",
		},
	}
	steps := upgrade(cfg)
	suite.Require().True(len(steps) > 1, "upgrade should generate at least two steps")
	suite.IsType(&run.AddRepo{}, steps[1])
}

func (suite *PlanTestSuite) TestUninstall() {
	cfg := Config{
		KubeToken:      "b2YgbXkgYWZmZWN0aW9u",
		SkipTLSVerify:  true,
		Certificate:    "cHJvY2xhaW1zIHdvbmRlcmZ1bCBmcmllbmRzaGlw",
		APIServer:      "98.765.43.21",
		ServiceAccount: "greathelm",
		DryRun:         true,
		Timeout:        "think about what you did",
		Release:        "jetta_id_love_to_change_the_world",
		KeepHistory:    true,
	}

	steps := uninstall(cfg)
	suite.Require().Equal(2, len(steps), "uninstall should return 2 steps")

	suite.Require().IsType(&run.InitKube{}, steps[0])
	init, _ := steps[0].(*run.InitKube)
	var expected Step = &run.InitKube{
		SkipTLSVerify:  true,
		Certificate:    "cHJvY2xhaW1zIHdvbmRlcmZ1bCBmcmllbmRzaGlw",
		APIServer:      "98.765.43.21",
		ServiceAccount: "greathelm",
		Token:          "b2YgbXkgYWZmZWN0aW9u",
		TemplateFile:   kubeConfigTemplate,
		ConfigFile:     kubeConfigFile,
	}

	suite.Equal(expected, init)

	suite.Require().IsType(&run.Uninstall{}, steps[1])
	actual, _ := steps[1].(*run.Uninstall)
	expected = &run.Uninstall{
		Release:     "jetta_id_love_to_change_the_world",
		DryRun:      true,
		KeepHistory: true,
	}
	suite.Equal(expected, actual)
}

func (suite *PlanTestSuite) TestUninstallWithUpdateDependencies() {
	cfg := Config{
		UpdateDependencies: true,
	}
	steps := uninstall(cfg)
	suite.Require().Equal(3, len(steps), "uninstall should have a third step when DepUpdate is true")
	suite.IsType(&run.InitKube{}, steps[0])
	suite.IsType(&run.DepUpdate{}, steps[1])
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
		ConfigFile:     kubeConfigFile,
	}
	suite.Equal(expected, init)
}

func (suite *PlanTestSuite) TestDepUpdate() {
	cfg := Config{
		UpdateDependencies: true,
		Chart:              "scatterplot",
	}

	steps := depUpdate(cfg)
	suite.Require().Equal(1, len(steps), "depUpdate should return one step")
	suite.Require().IsType(&run.DepUpdate{}, steps[0])
	update, _ := steps[0].(*run.DepUpdate)

	expected := &run.DepUpdate{
		Chart: "scatterplot",
	}
	suite.Equal(expected, update)
}

func (suite *PlanTestSuite) TestAddRepos() {
	cfg := Config{
		AddRepos: []string{
			"first=https://add.repos/one",
			"second=https://add.repos/two",
		},
		RepoCAFile: "state_licensure.repo.cert",
	}
	steps := addRepos(cfg)
	suite.Require().Equal(2, len(steps), "addRepos should add one step per repo")
	suite.Require().IsType(&run.AddRepo{}, steps[0])
	suite.Require().IsType(&run.AddRepo{}, steps[1])
	first := steps[0].(*run.AddRepo)
	second := steps[1].(*run.AddRepo)

	suite.Equal(first.Repo, "first=https://add.repos/one")
	suite.Equal(second.Repo, "second=https://add.repos/two")
	suite.Equal(first.CAFile, "state_licensure.repo.cert")
	suite.Equal(second.CAFile, "state_licensure.repo.cert")
}

func (suite *PlanTestSuite) TestLint() {
	cfg := Config{
		Chart:        "./flow",
		Values:       "steadfastness,forthrightness",
		StringValues: "tensile_strength,flexibility",
		ValuesFiles:  []string{"/root/price_inventory.yml"},
		LintStrictly: true,
	}

	steps := lint(cfg)
	suite.Equal(1, len(steps))

	want := &run.Lint{
		Chart:        "./flow",
		Values:       "steadfastness,forthrightness",
		StringValues: "tensile_strength,flexibility",
		ValuesFiles:  []string{"/root/price_inventory.yml"},
		Strict:       true,
	}
	suite.Equal(want, steps[0])
}

func (suite *PlanTestSuite) TestLintWithUpdateDependencies() {
	cfg := Config{
		UpdateDependencies: true,
	}
	steps := lint(cfg)
	suite.Require().Equal(2, len(steps), "lint should have a second step when DepUpdate is true")
	suite.IsType(&run.DepUpdate{}, steps[0])
}

func (suite *PlanTestSuite) TestLintWithAddRepos() {
	cfg := Config{
		AddRepos: []string{"friendczar=https://github.com/logan_pierce/friendczar"},
	}
	steps := lint(cfg)
	suite.Require().True(len(steps) > 0, "lint should return at least one step")
	suite.IsType(&run.AddRepo{}, steps[0])
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

func (suite *PlanTestSuite) TestDeterminePlanUninstallCommand() {
	cfg := Config{
		Command: "uninstall",
	}
	stepsMaker := determineSteps(cfg)
	suite.Same(&uninstall, stepsMaker)
}

// helm_command = delete is provided as an alias for backward-compatibility with drone-helm
func (suite *PlanTestSuite) TestDeterminePlanDeleteCommand() {
	cfg := Config{
		Command: "delete",
	}
	stepsMaker := determineSteps(cfg)
	suite.Same(&uninstall, stepsMaker)
}

func (suite *PlanTestSuite) TestDeterminePlanDeleteFromDroneEvent() {
	cfg := Config{
		DroneEvent: "delete",
	}
	stepsMaker := determineSteps(cfg)
	suite.Same(&uninstall, stepsMaker)
}

func (suite *PlanTestSuite) TestDeterminePlanLintCommand() {
	cfg := Config{
		Command: "lint",
	}

	stepsMaker := determineSteps(cfg)
	suite.Same(&lint, stepsMaker)
}

func (suite *PlanTestSuite) TestDeterminePlanHelpCommand() {
	cfg := Config{
		Command: "help",
	}

	stepsMaker := determineSteps(cfg)
	suite.Same(&help, stepsMaker)
}

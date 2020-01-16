package helm

import (
	"fmt"
	"github.com/pelotech/drone-helm3/internal/env"
	"github.com/pelotech/drone-helm3/internal/run"
	"os"
)

const (
	kubeConfigTemplate = "/root/.kube/config.tpl"
	kubeConfigFile     = "/root/.kube/config"
)

// A Step is one step in the plan.
type Step interface {
	Prepare(run.Config) error
	Execute(run.Config) error
}

// A Plan is a series of steps to perform.
type Plan struct {
	steps  []Step
	cfg    env.Config
	runCfg run.Config
}

// NewPlan makes a plan for running a helm operation.
func NewPlan(cfg env.Config) (*Plan, error) {
	p := Plan{
		cfg: cfg,
		runCfg: run.Config{
			Debug:     cfg.Debug,
			Namespace: cfg.Namespace,
			Stdout:    cfg.Stdout,
			Stderr:    cfg.Stderr,
		},
	}

	p.steps = (*determineSteps(cfg))(cfg)

	for i, step := range p.steps {
		if cfg.Debug {
			fmt.Fprintf(os.Stderr, "calling %T.Prepare (step %d)\n", step, i)
		}

		if err := step.Prepare(p.runCfg); err != nil {
			err = fmt.Errorf("while preparing %T step: %w", step, err)
			return nil, err
		}
	}

	return &p, nil
}

// determineSteps is primarily for the tests' convenience: it allows testing the "which stuff should
// we do" logic without building a config that meets all the steps' requirements.
func determineSteps(cfg env.Config) *func(env.Config) []Step {
	switch cfg.Command {
	case "upgrade":
		return &upgrade
	case "uninstall", "delete":
		return &uninstall
	case "lint":
		return &lint
	case "help":
		return &help
	default:
		switch cfg.DroneEvent {
		case "push", "tag", "deployment", "pull_request", "promote", "rollback":
			return &upgrade
		case "delete":
			return &uninstall
		default:
			return &help
		}
	}
}

// Execute runs each step in the plan, aborting and reporting on error
func (p *Plan) Execute() error {
	for i, step := range p.steps {
		if p.cfg.Debug {
			fmt.Fprintf(p.cfg.Stderr, "calling %T.Execute (step %d)\n", step, i)
		}

		if err := step.Execute(p.runCfg); err != nil {
			return fmt.Errorf("while executing %T step: %w", step, err)
		}
	}

	return nil
}

var upgrade = func(cfg env.Config) []Step {
	steps := initKube(cfg)
	steps = append(steps, addRepos(cfg)...)
	if cfg.UpdateDependencies {
		steps = append(steps, depUpdate(cfg)...)
	}
	steps = append(steps, run.NewUpgrade(cfg))

	return steps
}

var uninstall = func(cfg env.Config) []Step {
	steps := initKube(cfg)
	if cfg.UpdateDependencies {
		steps = append(steps, depUpdate(cfg)...)
	}
	steps = append(steps, run.NewUninstall(cfg))

	return steps
}

var lint = func(cfg env.Config) []Step {
	steps := addRepos(cfg)
	if cfg.UpdateDependencies {
		steps = append(steps, depUpdate(cfg)...)
	}
	steps = append(steps, run.NewLint(cfg))
	return steps
}

var help = func(cfg env.Config) []Step {
	return []Step{run.NewHelp(cfg)}
}

func initKube(cfg env.Config) []Step {
	return []Step{run.NewInitKube(cfg, kubeConfigTemplate, kubeConfigFile)}
}

func addRepos(cfg env.Config) []Step {
	steps := make([]Step, 0)
	for _, repo := range cfg.AddRepos {
		steps = append(steps, run.NewAddRepo(repo))
	}

	return steps
}

func depUpdate(cfg env.Config) []Step {
	return []Step{run.NewDepUpdate(cfg)}
}

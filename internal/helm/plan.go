package helm

import (
	"errors"
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
	Prepare() error
	Execute() error
}

// A Plan is a series of steps to perform.
type Plan struct {
	steps []Step
	cfg   env.Config
}

// NewPlan makes a plan for running a helm operation.
func NewPlan(cfg env.Config) (*Plan, error) {
	p := Plan{
		cfg: cfg,
	}

	if cfg.UpdateDependencies && cfg.DependenciesAction != "" {
		return nil, errors.New("update_dependencies is deprecated and cannot be provided together with dependencies_action")
	}

	p.steps = (*determineSteps(cfg))(cfg)

	for i, step := range p.steps {
		if cfg.Debug {
			fmt.Fprintf(os.Stderr, "calling %T.Prepare (step %d)\n", step, i)
		}

		if err := step.Prepare(); err != nil {
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

		if err := step.Execute(); err != nil {
			return fmt.Errorf("while executing %T step: %w", step, err)
		}
	}

	return nil
}

var upgrade = func(cfg env.Config) []Step {
	var steps []Step
	if !cfg.KubeInitSkip {
		steps = append(steps, run.NewInitKube(cfg, kubeConfigTemplate, kubeConfigFile))
	}
	for _, repo := range cfg.AddRepos {
		steps = append(steps, run.NewAddRepo(cfg, repo))
	}

	if cfg.DependenciesAction != "" {
		steps = append(steps, run.NewDepAction(cfg))
	}

	if cfg.UpdateDependencies {
		steps = append(steps, run.NewDepUpdate(cfg))
	}

	steps = append(steps, run.NewUpgrade(cfg))

	return steps
}

var uninstall = func(cfg env.Config) []Step {
	var steps []Step
	if !cfg.KubeInitSkip {
		steps = append(steps, run.NewInitKube(cfg, kubeConfigTemplate, kubeConfigFile))
	}
	if cfg.UpdateDependencies {
		steps = append(steps, run.NewDepUpdate(cfg))
	}
	steps = append(steps, run.NewUninstall(cfg))

	return steps
}

var lint = func(cfg env.Config) []Step {
	var steps []Step
	for _, repo := range cfg.AddRepos {
		steps = append(steps, run.NewAddRepo(cfg, repo))
	}
	if cfg.UpdateDependencies {
		steps = append(steps, run.NewDepUpdate(cfg))
	}
	steps = append(steps, run.NewLint(cfg))
	return steps
}

var help = func(cfg env.Config) []Step {
	return []Step{run.NewHelp(cfg)}
}

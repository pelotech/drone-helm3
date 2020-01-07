package helm

import (
	"fmt"
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
	cfg    Config
	runCfg run.Config
}

// NewPlan makes a plan for running a helm operation.
func NewPlan(cfg Config) (*Plan, error) {
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
func determineSteps(cfg Config) *func(Config) []Step {
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

var upgrade = func(cfg Config) []Step {
	steps := initKube(cfg)
	steps = append(steps, addRepos(cfg)...)
	if cfg.UpdateDependencies {
		steps = append(steps, depUpdate(cfg)...)
	}
	steps = append(steps, &run.Upgrade{
		Chart:         cfg.Chart,
		Release:       cfg.Release,
		ChartVersion:  cfg.ChartVersion,
		DryRun:        cfg.DryRun,
		Wait:          cfg.Wait,
		Values:        cfg.Values,
		StringValues:  cfg.StringValues,
		ValuesFiles:   cfg.ValuesFiles,
		ReuseValues:   cfg.ReuseValues,
		Timeout:       cfg.Timeout,
		Force:         cfg.Force,
		Atomic:        cfg.AtomicUpgrade,
		CleanupOnFail: cfg.CleanupOnFail,
	})

	return steps
}

var uninstall = func(cfg Config) []Step {
	steps := initKube(cfg)
	if cfg.UpdateDependencies {
		steps = append(steps, depUpdate(cfg)...)
	}
	steps = append(steps, &run.Uninstall{
		Release:     cfg.Release,
		DryRun:      cfg.DryRun,
		KeepHistory: cfg.KeepHistory,
	})

	return steps
}

var lint = func(cfg Config) []Step {
	steps := addRepos(cfg)
	if cfg.UpdateDependencies {
		steps = append(steps, depUpdate(cfg)...)
	}
	steps = append(steps, &run.Lint{
		Chart:        cfg.Chart,
		Values:       cfg.Values,
		StringValues: cfg.StringValues,
		ValuesFiles:  cfg.ValuesFiles,
		Strict:       cfg.LintStrictly,
	})

	return steps
}

var help = func(cfg Config) []Step {
	help := &run.Help{
		HelmCommand: cfg.Command,
	}
	return []Step{help}
}

func initKube(cfg Config) []Step {
	return []Step{
		&run.InitKube{
			SkipTLSVerify:  cfg.SkipTLSVerify,
			Certificate:    cfg.Certificate,
			APIServer:      cfg.APIServer,
			ServiceAccount: cfg.ServiceAccount,
			Token:          cfg.KubeToken,
			TemplateFile:   kubeConfigTemplate,
			ConfigFile:     kubeConfigFile,
		},
	}
}

func addRepos(cfg Config) []Step {
	steps := make([]Step, 0)
	for _, repo := range cfg.AddRepos {
		steps = append(steps, &run.AddRepo{
			Repo: repo,
		})
	}

	return steps
}

func depUpdate(cfg Config) []Step {
	return []Step{
		&run.DepUpdate{
			Chart: cfg.Chart,
		},
	}
}

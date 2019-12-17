package helm

import (
	"fmt"
	"github.com/pelotech/drone-helm3/internal/run"
	"os"
)

const kubeConfigTemplate = "/root/.kube/config.tpl"

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
			Debug:        cfg.Debug,
			KubeConfig:   cfg.KubeConfig,
			Values:       cfg.Values,
			StringValues: cfg.StringValues,
			ValuesFiles:  cfg.ValuesFiles,
			Namespace:    cfg.Namespace,
			Stdout:       os.Stdout,
			Stderr:       os.Stderr,
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
	case "delete":
		panic("not implemented")
	case "lint":
		panic("not implemented")
	case "help":
		return &help
	default:
		switch cfg.DroneEvent {
		case "push", "tag", "deployment", "pull_request", "promote", "rollback":
			return &upgrade
		default:
			panic("not implemented")
		}
	}
}

// Execute runs each step in the plan, aborting and reporting on error
func (p *Plan) Execute() error {
	for i, step := range p.steps {
		if p.cfg.Debug {
			fmt.Fprintf(os.Stderr, "calling %T.Execute (step %d)\n", step, i)
		}

		if err := step.Execute(p.runCfg); err != nil {
			return fmt.Errorf("in execution step %d: %w", i, err)
		}
	}

	return nil
}

var upgrade = func(cfg Config) []Step {
	steps := make([]Step, 0)

	steps = append(steps, &run.InitKube{
		SkipTLSVerify:  cfg.SkipTLSVerify,
		Certificate:    cfg.Certificate,
		APIServer:      cfg.APIServer,
		ServiceAccount: cfg.ServiceAccount,
		Token:          cfg.KubeToken,
		TemplateFile:   kubeConfigTemplate,
	})

	steps = append(steps, &run.Upgrade{
		Chart:        cfg.Chart,
		Release:      cfg.Release,
		ChartVersion: cfg.ChartVersion,
		DryRun:       cfg.DryRun,
		Wait:         cfg.Wait,
		ReuseValues:  cfg.ReuseValues,
		Timeout:      cfg.Timeout,
		Force:        cfg.Force,
	})

	return steps
}

var help = func(cfg Config) []Step {
	help := &run.Help{}
	return []Step{help}
}

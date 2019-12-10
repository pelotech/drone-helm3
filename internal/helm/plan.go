package helm

import (
	"errors"
	"fmt"
	"github.com/pelotech/drone-helm3/internal/run"
	"os"
)

// A Step is one step in the plan.
type Step interface {
	Prepare(run.Config) error
	Execute() error
}

// A Plan is a series of steps to perform.
type Plan struct {
	steps []Step
}

// NewPlan makes a plan for running a helm operation.
func NewPlan(cfg Config) (*Plan, error) {
	runCfg := run.Config{
		Debug:          cfg.Debug,
		KubeConfig:     cfg.KubeConfig,
		Values:         cfg.Values,
		StringValues:   cfg.StringValues,
		ValuesFiles:    cfg.ValuesFiles,
		Namespace:      cfg.Namespace,
		Token:          cfg.Token,
		SkipTLSVerify:  cfg.SkipTLSVerify,
		Certificate:    cfg.Certificate,
		APIServer:      cfg.APIServer,
		ServiceAccount: cfg.ServiceAccount,
		Stdout:         os.Stdout,
		Stderr:         os.Stderr,
	}

	p := Plan{}
	switch cfg.Command {
	case "upgrade":
		steps, err := upgrade(cfg, runCfg)
		if err != nil {
			return nil, err
		}
		p.steps = steps
	case "delete":
		return nil, errors.New("not implemented")
	case "lint":
		return nil, errors.New("not implemented")
	case "help":
		steps, err := help(cfg, runCfg)
		if err != nil {
			return nil, err
		}
		p.steps = steps
	default:
		switch cfg.DroneEvent {
		case "push", "tag", "deployment", "pull_request", "promote", "rollback":
			steps, err := upgrade(cfg, runCfg)
			if err != nil {
				return nil, err
			}
			p.steps = steps
		default:
			return nil, errors.New("not implemented")
		}
	}

	return &p, nil
}

// Execute runs each step in the plan, aborting and reporting on error
func (p *Plan) Execute() error {
	for _, step := range p.steps {
		if err := step.Execute(); err != nil {
			return err
		}
	}

	return nil
}

func upgrade(cfg Config, runCfg run.Config) ([]Step, error) {
	steps := make([]Step, 0)
	upgrade := &run.Upgrade{
		Chart:        cfg.Chart,
		Release:      cfg.Release,
		ChartVersion: cfg.ChartVersion,
		Wait:         cfg.Wait,
		ReuseValues:  cfg.ReuseValues,
		Timeout:      cfg.Timeout,
		Force:        cfg.Force,
	}
	if err := upgrade.Prepare(runCfg); err != nil {
		err = fmt.Errorf("while preparing upgrade step: %w", err)
		return steps, err
	}
	steps = append(steps, upgrade)

	return steps, nil
}

func help(cfg Config, runCfg run.Config) ([]Step, error) {
	help := &run.Help{}

	if err := help.Prepare(runCfg); err != nil {
		err = fmt.Errorf("while preparing help step: %w", err)
		return []Step{}, err
	}

	return []Step{help}, nil
}

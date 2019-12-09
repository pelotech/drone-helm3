package helm

import (
	"errors"
	"github.com/pelotech/drone-helm3/internal/run"
)

type Step interface {
	Run() error
}

type Plan struct {
	steps []Step
}

func NewPlan(cfg Config) (*Plan, error) {
	p := Plan{}
	switch cfg.Command {
	case "upgrade":
		steps, err := upgrade(cfg)
		if err != nil {
			return nil, err
		}
		p.steps = steps
	case "delete":
		return nil, errors.New("not implemented")
	case "lint":
		return nil, errors.New("not implemented")
	case "help":
		return nil, errors.New("not implemented")
	default:
		switch cfg.DroneEvent {
		case "push", "tag", "deployment", "pull_request", "promote", "rollback":
			steps, err := upgrade(cfg)
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

func (p *Plan) Execute() error {
	for _, step := range p.steps {
		if err := step.Run(); err != nil {
			return err
		}
	}

	return nil
}

func upgrade(cfg Config) ([]Step, error) {
	steps := make([]Step, 0)
	steps = append(steps, run.NewUpgrade(cfg.Release, cfg.Chart))

	return steps, nil
}

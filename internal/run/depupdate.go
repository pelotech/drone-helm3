package run

import (
	"fmt"
	"github.com/pelotech/drone-helm3/internal/env"
)

// DepUpdate is an execution step that calls `helm dependency update` when executed.
type DepUpdate struct {
	Chart string
	cmd   cmd
}

// NewDepUpdate creates a DepUpdate using fields from the given Config. No validation is performed at this time.
func NewDepUpdate(cfg env.Config) *DepUpdate {
	return &DepUpdate{
		Chart: cfg.Chart,
	}
}

// Execute executes the `helm upgrade` command.
func (d *DepUpdate) Execute() error {
	return d.cmd.Run()
}

// Prepare gets the DepUpdate ready to execute.
func (d *DepUpdate) Prepare(cfg Config) error {
	if d.Chart == "" {
		return fmt.Errorf("chart is required")
	}

	args := make([]string, 0)

	if cfg.Namespace != "" {
		args = append(args, "--namespace", cfg.Namespace)
	}
	if cfg.Debug {
		args = append(args, "--debug")
	}

	args = append(args, "dependency", "update", d.Chart)

	d.cmd = command(helmBin, args...)
	d.cmd.Stdout(cfg.Stdout)
	d.cmd.Stderr(cfg.Stderr)

	if cfg.Debug {
		fmt.Fprintf(cfg.Stderr, "Generated command: '%s'\n", d.cmd.String())
	}

	return nil
}

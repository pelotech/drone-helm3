package run

import (
	"fmt"
	"github.com/pelotech/drone-helm3/internal/env"
)

// DepUpdate is an execution step that calls `helm dependency update` when executed.
type DepUpdate struct {
	*config
	chart string
	cmd   cmd
}

// NewDepUpdate creates a DepUpdate using fields from the given Config. No validation is performed at this time.
func NewDepUpdate(cfg env.Config) *DepUpdate {
	return &DepUpdate{
		config: newConfig(cfg),
		chart:  cfg.Chart,
	}
}

// Execute executes the `helm upgrade` command.
func (d *DepUpdate) Execute() error {
	return d.cmd.Run()
}

// Prepare gets the DepUpdate ready to execute.
func (d *DepUpdate) Prepare() error {
	if d.chart == "" {
		return fmt.Errorf("chart is required")
	}

	args := make([]string, 0)

	if d.namespace != "" {
		args = append(args, "--namespace", d.namespace)
	}
	if d.debug {
		args = append(args, "--debug")
	}

	args = append(args, "dependency", "update", d.chart)

	d.cmd = command(helmBin, args...)
	d.cmd.Stdout(d.stdout)
	d.cmd.Stderr(d.stderr)

	if d.debug {
		fmt.Fprintf(d.stderr, "Generated command: '%s'\n", d.cmd.String())
	}

	return nil
}

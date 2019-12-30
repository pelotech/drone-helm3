package run

import (
	"fmt"
)

// DepUpdate is an execution step that calls `helm dependency update` when executed.
type DepUpdate struct {
	Chart string
	cmd   cmd
}

// Execute executes the `helm upgrade` command.
func (d *DepUpdate) Execute(_ Config) error {
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

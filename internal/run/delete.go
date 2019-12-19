package run

import (
	"fmt"
)

// Delete is an execution step that calls `helm upgrade` when executed.
type Delete struct {
	Release string
	DryRun  bool
	cmd     cmd
}

// Execute executes the `helm upgrade` command.
func (d *Delete) Execute(_ Config) error {
	return d.cmd.Run()
}

// Prepare gets the Delete ready to execute.
func (d *Delete) Prepare(cfg Config) error {
	if d.Release == "" {
		return fmt.Errorf("release is required")
	}

	args := []string{"--kubeconfig", cfg.KubeConfig}

	if cfg.Namespace != "" {
		args = append(args, "--namespace", cfg.Namespace)
	}
	if cfg.Debug {
		args = append(args, "--debug")
	}

	args = append(args, "delete")

	if d.DryRun {
		args = append(args, "--dry-run")
	}

	args = append(args, d.Release)

	d.cmd = command(helmBin, args...)
	d.cmd.Stdout(cfg.Stdout)
	d.cmd.Stderr(cfg.Stderr)

	if cfg.Debug {
		fmt.Fprintf(cfg.Stderr, "Generated command: '%s'\n", d.cmd.String())
	}

	return nil
}

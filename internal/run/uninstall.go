package run

import (
	"fmt"
)

// Uninstall is an execution step that calls `helm uninstall` when executed.
type Uninstall struct {
	Release     string
	DryRun      bool
	KeepHistory bool
	cmd         cmd
}

// Execute executes the `helm uninstall` command.
func (u *Uninstall) Execute(_ Config) error {
	return u.cmd.Run()
}

// Prepare gets the Uninstall ready to execute.
func (u *Uninstall) Prepare(cfg Config) error {
	if u.Release == "" {
		return fmt.Errorf("release is required")
	}

	args := make([]string, 0)

	if cfg.Namespace != "" {
		args = append(args, "--namespace", cfg.Namespace)
	}
	if cfg.Debug {
		args = append(args, "--debug")
	}

	args = append(args, "uninstall")

	if u.DryRun {
		args = append(args, "--dry-run")
	}
	if u.KeepHistory {
		args = append(args, "--keep-history")
	}

	args = append(args, u.Release)

	u.cmd = command(helmBin, args...)
	u.cmd.Stdout(cfg.Stdout)
	u.cmd.Stderr(cfg.Stderr)

	if cfg.Debug {
		fmt.Fprintf(cfg.Stderr, "Generated command: '%s'\n", u.cmd.String())
	}

	return nil
}

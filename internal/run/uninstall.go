package run

import (
	"fmt"
	"github.com/pelotech/drone-helm3/internal/env"
)

// Uninstall is an execution step that calls `helm uninstall` when executed.
type Uninstall struct {
	*config
	release     string
	dryRun      bool
	keepHistory bool
	cmd         cmd
}

// NewUninstall creates an Uninstall using fields from the given Config. No validation is performed at this time.
func NewUninstall(cfg env.Config) *Uninstall {
	return &Uninstall{
		config:      newConfig(cfg),
		release:     cfg.Release,
		dryRun:      cfg.DryRun,
		keepHistory: cfg.KeepHistory,
	}
}

// Execute executes the `helm uninstall` command.
func (u *Uninstall) Execute() error {
	return u.cmd.Run()
}

// Prepare gets the Uninstall ready to execute.
func (u *Uninstall) Prepare() error {
	if u.release == "" {
		return fmt.Errorf("release is required")
	}

	args := u.globalFlags()
	args = append(args, "uninstall")

	if u.dryRun {
		args = append(args, "--dry-run")
	}
	if u.keepHistory {
		args = append(args, "--keep-history")
	}

	args = append(args, u.release)

	u.cmd = command(helmBin, args...)
	u.cmd.Stdout(u.stdout)
	u.cmd.Stderr(u.stderr)

	if u.debug {
		fmt.Fprintf(u.stderr, "Generated command: '%s'\n", u.cmd.String())
	}

	return nil
}

package run

import (
	"fmt"
)

// Upgrade is an execution step that calls `helm upgrade` when executed.
type Upgrade struct {
	Chart   string
	Release string

	ChartVersion string
	Wait         bool
	ReuseValues  bool
	Timeout      string
	Force        bool

	cmd cmd
}

// Execute executes the `helm upgrade` command.
func (u *Upgrade) Execute(_ Config) error {
	return u.cmd.Run()
}

// Prepare gets the Upgrade ready to execute.
func (u *Upgrade) Prepare(cfg Config) error {
	args := []string{"--kubeconfig", cfg.KubeConfig}

	if cfg.Namespace != "" {
		args = append(args, "--namespace", cfg.Namespace)
	}

	args = append(args, "upgrade", "--install", u.Release, u.Chart)

	if cfg.Debug {
		args = append([]string{"--debug"}, args...)
	}

	u.cmd = command(helmBin, args...)
	u.cmd.Stdout(cfg.Stdout)
	u.cmd.Stderr(cfg.Stderr)

	if cfg.Debug {
		fmt.Fprintf(cfg.Stderr, "Generated command: '%s'\n", u.cmd.String())
	}

	return nil
}

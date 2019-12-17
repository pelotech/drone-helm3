package run

import (
	"fmt"
)

// Upgrade is an execution step that calls `helm upgrade` when executed.
type Upgrade struct {
	Chart   string
	Release string

	ChartVersion string
	DryRun       bool
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
	if u.Chart == "" {
		return fmt.Errorf("chart is required")
	}
	if u.Release == "" {
		return fmt.Errorf("release is required")
	}

	args := []string{"--kubeconfig", cfg.KubeConfig}

	if cfg.Namespace != "" {
		args = append(args, "--namespace", cfg.Namespace)
	}
	if cfg.Debug {
		args = append(args, "--debug")
	}

	args = append(args, "upgrade", "--install")

	if u.ChartVersion != "" {
		args = append(args, "--version", u.ChartVersion)
	}
	if u.DryRun {
		args = append(args, "--dry-run")
	}
	if u.Wait {
		args = append(args, "--wait")
	}
	if u.ReuseValues {
		args = append(args, "--reuse-values")
	}
	if u.Timeout != "" {
		args = append(args, "--timeout", u.Timeout)
	}
	if u.Force {
		args = append(args, "--force")
	}
	if cfg.Values != "" {
		args = append(args, "--set", cfg.Values)
	}
	if cfg.StringValues != "" {
		args = append(args, "--set-string", cfg.StringValues)
	}
	for _, vFile := range cfg.ValuesFiles {
		args = append(args, "--values", vFile)
	}

	args = append(args, u.Release, u.Chart)
	u.cmd = command(helmBin, args...)
	u.cmd.Stdout(cfg.Stdout)
	u.cmd.Stderr(cfg.Stderr)

	if cfg.Debug {
		fmt.Fprintf(cfg.Stderr, "Generated command: '%s'\n", u.cmd.String())
	}

	return nil
}

package run

import (
	"fmt"
	"github.com/pelotech/drone-helm3/internal/env"
)

// Upgrade is an execution step that calls `helm upgrade` when executed.
type Upgrade struct {
	Chart   string
	Release string

	ChartVersion  string
	DryRun        bool
	Wait          bool
	Values        string
	StringValues  string
	ValuesFiles   []string
	ReuseValues   bool
	Timeout       string
	Force         bool
	Atomic        bool
	CleanupOnFail bool

	cmd cmd
}

// NewUpgrade creates an Upgrade using fields from the given Config. No validation is performed at this time.
func NewUpgrade(cfg env.Config) *Upgrade {
	return &Upgrade{
		Chart:         cfg.Chart,
		Release:       cfg.Release,
		ChartVersion:  cfg.ChartVersion,
		DryRun:        cfg.DryRun,
		Wait:          cfg.Wait,
		Values:        cfg.Values,
		StringValues:  cfg.StringValues,
		ValuesFiles:   cfg.ValuesFiles,
		ReuseValues:   cfg.ReuseValues,
		Timeout:       cfg.Timeout,
		Force:         cfg.Force,
		Atomic:        cfg.AtomicUpgrade,
		CleanupOnFail: cfg.CleanupOnFail,
	}
}

// Execute executes the `helm upgrade` command.
func (u *Upgrade) Execute() error {
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

	args := make([]string, 0)

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
	if u.Atomic {
		args = append(args, "--atomic")
	}
	if u.CleanupOnFail {
		args = append(args, "--cleanup-on-fail")
	}
	if u.Values != "" {
		args = append(args, "--set", u.Values)
	}
	if u.StringValues != "" {
		args = append(args, "--set-string", u.StringValues)
	}
	for _, vFile := range u.ValuesFiles {
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

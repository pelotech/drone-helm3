package run

import (
	"fmt"
	"github.com/pelotech/drone-helm3/internal/env"
)

// Upgrade is an execution step that calls `helm upgrade` when executed.
type Upgrade struct {
	*config
	chart   string
	release string

	chartVersion  string
	dryRun        bool
	wait          bool
	values        string
	stringValues  string
	valuesFiles   []string
	reuseValues   bool
	timeout       string
	force         bool
	atomic        bool
	cleanupOnFail bool
	certs         *repoCerts

	cmd cmd
}

// NewUpgrade creates an Upgrade using fields from the given Config. No validation is performed at this time.
func NewUpgrade(cfg env.Config) *Upgrade {
	return &Upgrade{
		config:        newConfig(cfg),
		chart:         cfg.Chart,
		release:       cfg.Release,
		chartVersion:  cfg.ChartVersion,
		dryRun:        cfg.DryRun,
		wait:          cfg.Wait,
		values:        cfg.Values,
		stringValues:  cfg.StringValues,
		valuesFiles:   cfg.ValuesFiles,
		reuseValues:   cfg.ReuseValues,
		timeout:       cfg.Timeout,
		force:         cfg.Force,
		atomic:        cfg.AtomicUpgrade,
		cleanupOnFail: cfg.CleanupOnFail,
		certs:         newRepoCerts(cfg),
	}
}

// Execute executes the `helm upgrade` command.
func (u *Upgrade) Execute() error {
	return u.cmd.Run()
}

// Prepare gets the Upgrade ready to execute.
func (u *Upgrade) Prepare() error {
	if u.chart == "" {
		return fmt.Errorf("chart is required")
	}
	if u.release == "" {
		return fmt.Errorf("release is required")
	}

	args := u.globalFlags()
	args = append(args, "upgrade", "--install")

	if u.chartVersion != "" {
		args = append(args, "--version", u.chartVersion)
	}
	if u.dryRun {
		args = append(args, "--dry-run")
	}
	if u.wait {
		args = append(args, "--wait")
	}
	if u.reuseValues {
		args = append(args, "--reuse-values")
	}
	if u.timeout != "" {
		args = append(args, "--timeout", u.timeout)
	}
	if u.force {
		args = append(args, "--force")
	}
	if u.atomic {
		args = append(args, "--atomic")
	}
	if u.cleanupOnFail {
		args = append(args, "--cleanup-on-fail")
	}
	if u.values != "" {
		args = append(args, "--set", u.values)
	}
	if u.stringValues != "" {
		args = append(args, "--set-string", u.stringValues)
	}
	for _, vFile := range u.valuesFiles {
		args = append(args, "--values", vFile)
	}
	args = append(args, u.certs.flags()...)

	args = append(args, u.release, u.chart)
	u.cmd = command(helmBin, args...)
	u.cmd.Stdout(u.stdout)
	u.cmd.Stderr(u.stderr)

	if u.debug {
		fmt.Fprintf(u.stderr, "Generated command: '%s'\n", u.cmd.String())
	}

	return nil
}

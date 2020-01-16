package run

import (
	"fmt"
	"github.com/pelotech/drone-helm3/internal/env"
)

// Lint is an execution step that calls `helm lint` when executed.
type Lint struct {
	Chart        string
	Values       string
	StringValues string
	ValuesFiles  []string
	Strict       bool
	cmd          cmd
}

// NewLint creates a Lint using fields from the given Config. No validation is performed at this time.
func NewLint(cfg env.Config) *Lint {
	return &Lint{
		Chart:        cfg.Chart,
		Values:       cfg.Values,
		StringValues: cfg.StringValues,
		ValuesFiles:  cfg.ValuesFiles,
		Strict:       cfg.LintStrictly,
	}
}

// Execute executes the `helm lint` command.
func (l *Lint) Execute(_ Config) error {
	return l.cmd.Run()
}

// Prepare gets the Lint ready to execute.
func (l *Lint) Prepare(cfg Config) error {
	if l.Chart == "" {
		return fmt.Errorf("chart is required")
	}

	args := make([]string, 0)

	if cfg.Namespace != "" {
		args = append(args, "--namespace", cfg.Namespace)
	}
	if cfg.Debug {
		args = append(args, "--debug")
	}

	args = append(args, "lint")

	if l.Values != "" {
		args = append(args, "--set", l.Values)
	}
	if l.StringValues != "" {
		args = append(args, "--set-string", l.StringValues)
	}
	for _, vFile := range l.ValuesFiles {
		args = append(args, "--values", vFile)
	}
	if l.Strict {
		args = append(args, "--strict")
	}

	args = append(args, l.Chart)

	l.cmd = command(helmBin, args...)
	l.cmd.Stdout(cfg.Stdout)
	l.cmd.Stderr(cfg.Stderr)

	if cfg.Debug {
		fmt.Fprintf(cfg.Stderr, "Generated command: '%s'\n", l.cmd.String())
	}

	return nil
}

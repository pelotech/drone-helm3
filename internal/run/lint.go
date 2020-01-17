package run

import (
	"fmt"
	"github.com/pelotech/drone-helm3/internal/env"
)

// Lint is an execution step that calls `helm lint` when executed.
type Lint struct {
	*config
	chart        string
	values       string
	stringValues string
	valuesFiles  []string
	strict       bool
	cmd          cmd
}

// NewLint creates a Lint using fields from the given Config. No validation is performed at this time.
func NewLint(cfg env.Config) *Lint {
	return &Lint{
		config:       newConfig(cfg),
		chart:        cfg.Chart,
		values:       cfg.Values,
		stringValues: cfg.StringValues,
		valuesFiles:  cfg.ValuesFiles,
		strict:       cfg.LintStrictly,
	}
}

// Execute executes the `helm lint` command.
func (l *Lint) Execute() error {
	return l.cmd.Run()
}

// Prepare gets the Lint ready to execute.
func (l *Lint) Prepare() error {
	if l.chart == "" {
		return fmt.Errorf("chart is required")
	}

	args := make([]string, 0)

	if l.namespace != "" {
		args = append(args, "--namespace", l.namespace)
	}
	if l.debug {
		args = append(args, "--debug")
	}

	args = append(args, "lint")

	if l.values != "" {
		args = append(args, "--set", l.values)
	}
	if l.stringValues != "" {
		args = append(args, "--set-string", l.stringValues)
	}
	for _, vFile := range l.valuesFiles {
		args = append(args, "--values", vFile)
	}
	if l.strict {
		args = append(args, "--strict")
	}

	args = append(args, l.chart)

	l.cmd = command(helmBin, args...)
	l.cmd.Stdout(l.stdout)
	l.cmd.Stderr(l.stderr)

	if l.debug {
		fmt.Fprintf(l.stderr, "Generated command: '%s'\n", l.cmd.String())
	}

	return nil
}

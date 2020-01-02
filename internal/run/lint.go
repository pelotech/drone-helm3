package run

import (
	"fmt"
)

// Lint is an execution step that calls `helm lint` when executed.
type Lint struct {
	Chart        string
	Values       string
	StringValues string
	ValuesFiles  []string
	cmd          cmd
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

	args = append(args, l.Chart)

	l.cmd = command(helmBin, args...)
	l.cmd.Stdout(cfg.Stdout)
	l.cmd.Stderr(cfg.Stderr)

	if cfg.Debug {
		fmt.Fprintf(cfg.Stderr, "Generated command: '%s'\n", l.cmd.String())
	}

	return nil
}

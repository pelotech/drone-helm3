package run

import (
// "fmt"
)

// Lint is an execution step that calls `helm lint` when executed.
type Lint struct {
	Chart string

	cmd cmd
}

// Execute executes the `helm lint` command.
func (l *Lint) Execute(_ Config) error {
	return l.cmd.Run()
}

// Prepare gets the Lint ready to execute.
func (l *Lint) Prepare(cfg Config) error {
	args := []string{"lint"}

	if cfg.Values != "" {
		args = append(args, "--set", cfg.Values)
	}
	if cfg.StringValues != "" {
		args = append(args, "--set-string", cfg.StringValues)
	}
	for _, vFile := range cfg.ValuesFiles {
		args = append(args, "--values", vFile)
	}

	args = append(args, l.Chart)

	l.cmd = command(helmBin, args...)
	l.cmd.Stdout(cfg.Stdout)
	l.cmd.Stderr(cfg.Stderr)

	return nil
}

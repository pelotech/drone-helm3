package run

import (
	"fmt"
)

// AddRepo is an execution step that calls `helm repo add` when executed.
type AddRepo struct {
	Name string
	URL  string
	cmd  cmd
}

// Execute executes the `helm repo add` command.
func (a *AddRepo) Execute(_ Config) error {
	return a.cmd.Run()
}

// Prepare gets the AddRepo ready to execute.
func (a *AddRepo) Prepare(cfg Config) error {
	if a.Name == "" {
		return fmt.Errorf("repo name is required")
	}
	if a.URL == "" {
		return fmt.Errorf("repo URL is required")
	}

	args := make([]string, 0)

	if cfg.Namespace != "" {
		args = append(args, "--namespace", cfg.Namespace)
	}
	if cfg.Debug {
		args = append(args, "--debug")
	}

	args = append(args, "repo", "add", a.Name, a.URL)

	a.cmd = command(helmBin, args...)
	a.cmd.Stdout(cfg.Stdout)
	a.cmd.Stderr(cfg.Stderr)

	if cfg.Debug {
		fmt.Fprintf(cfg.Stderr, "Generated command: '%s'\n", a.cmd.String())
	}

	return nil
}

package run

import (
	"fmt"
	"strings"
)

// AddRepo is an execution step that calls `helm repo add` when executed.
type AddRepo struct {
	Repo   string
	CAFile string
	cmd    cmd
}

// Execute executes the `helm repo add` command.
func (a *AddRepo) Execute(_ Config) error {
	return a.cmd.Run()
}

// Prepare gets the AddRepo ready to execute.
func (a *AddRepo) Prepare(cfg Config) error {
	if a.Repo == "" {
		return fmt.Errorf("repo is required")
	}
	split := strings.SplitN(a.Repo, "=", 2)
	if len(split) != 2 {
		return fmt.Errorf("bad repo spec '%s'", a.Repo)
	}

	name := split[0]
	url := split[1]

	args := make([]string, 0)

	if cfg.Namespace != "" {
		args = append(args, "--namespace", cfg.Namespace)
	}
	if cfg.Debug {
		args = append(args, "--debug")
	}

	args = append(args, "repo", "add")
	if a.CAFile != "" {
		args = append(args, "--ca-file", a.CAFile)
	}
	args = append(args, name, url)

	a.cmd = command(helmBin, args...)
	a.cmd.Stdout(cfg.Stdout)
	a.cmd.Stderr(cfg.Stderr)

	if cfg.Debug {
		fmt.Fprintf(cfg.Stderr, "Generated command: '%s'\n", a.cmd.String())
	}

	return nil
}

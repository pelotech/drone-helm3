package run

import (
	"fmt"
	"github.com/pelotech/drone-helm3/internal/env"
	"strings"
)

// AddRepo is an execution step that calls `helm repo add` when executed.
type AddRepo struct {
	*config
	repo string
	cmd  cmd
}

// NewAddRepo creates an AddRepo for the given repo-spec. No validation is performed at this time.
func NewAddRepo(cfg env.Config, repo string) *AddRepo {
	return &AddRepo{
		config: newConfig(cfg),
		repo:   repo,
	}
}

// Execute executes the `helm repo add` command.
func (a *AddRepo) Execute() error {
	return a.cmd.Run()
}

// Prepare gets the AddRepo ready to execute.
func (a *AddRepo) Prepare() error {
	if a.repo == "" {
		return fmt.Errorf("repo is required")
	}
	split := strings.SplitN(a.repo, "=", 2)
	if len(split) != 2 {
		return fmt.Errorf("bad repo spec '%s'", a.repo)
	}

	name := split[0]
	url := split[1]

	args := a.globalFlags()
	args = append(args, "repo", "add", name, url)

	a.cmd = command(helmBin, args...)
	a.cmd.Stdout(a.stdout)
	a.cmd.Stderr(a.stderr)

	if a.debug {
		fmt.Fprintf(a.stderr, "Generated command: '%s'\n", a.cmd.String())
	}

	return nil
}

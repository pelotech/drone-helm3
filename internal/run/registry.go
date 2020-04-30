package run

import (
	"fmt"

	"github.com/pelotech/drone-helm3/internal/env"
)

// Registry is an execution step that calls `helm registry` when executed.
type Registry struct {
	*config
	userID      string
	password    string
	registryURL string
	subCommand  string
	cmd         cmd
}

// NewRegistry creates a Registry using fields from the given Config. No validation is performed at this time.
func NewRegistry(subCommand string, cfg env.Config) *Registry {
	return &Registry{
		config:      newConfig(cfg),
		userID:      cfg.RegistryLoginUserID,
		password:    cfg.RegistryLoginPassword,
		registryURL: cfg.RegistryURL,
		subCommand:  subCommand,
	}
}

// Execute executes the `helm registry` command.
func (reg *Registry) Execute() error {
	return reg.cmd.Run()
}

// Prepare gets the Registry ready to execute.
func (reg *Registry) Prepare() error {
	args := []string{}

	args = append(args, "registry")

	if reg.subCommand == "login" {
		args = append(args, "login")
		args = append(args, "-u", reg.userID)
		args = append(args, "-p", reg.password)
		args = append(args, reg.registryURL)
	} else if reg.subCommand == "logout" {
		args = append(args, "logout")
		args = append(args, reg.registryURL)
	}

	reg.cmd = command(helmBin, args...)
	reg.cmd.Stdout(reg.stdout)
	reg.cmd.Stderr(reg.stderr)

	if reg.debug {
		fmt.Fprintf(reg.stderr, "Generated command: '%s'\n", reg.cmd.String())
	}

	return nil
}

package run

import (
	"fmt"
	"github.com/pelotech/drone-helm3/internal/env"
)

// Help is a step in a helm Plan that calls `helm help`.
type Help struct {
	*config
	helmCommand string
	cmd         cmd
}

// NewHelp creates a Help using fields from the given Config. No validation is performed at this time.
func NewHelp(cfg env.Config) *Help {
	return &Help{
		config:      newConfig(cfg),
		helmCommand: cfg.Command,
	}
}

// Execute executes the `helm help` command.
func (h *Help) Execute() error {
	if err := h.cmd.Run(); err != nil {
		return fmt.Errorf("while running '%s': %w", h.cmd.String(), err)
	}

	if h.helmCommand == "help" {
		return nil
	}
	return fmt.Errorf("unknown command '%s'", h.helmCommand)
}

// Prepare gets the Help ready to execute.
func (h *Help) Prepare() error {
	args := h.globalFlags()
	args = append(args, "help")

	h.cmd = command(helmBin, args...)
	h.cmd.Stdout(h.stdout)
	h.cmd.Stderr(h.stderr)

	if h.debug {
		fmt.Fprintf(h.stderr, "Generated command: '%s'\n", h.cmd.String())
	}

	return nil
}

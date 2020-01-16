package run

import (
	"fmt"
	"github.com/pelotech/drone-helm3/internal/env"
)

// Help is a step in a helm Plan that calls `helm help`.
type Help struct {
	HelmCommand string
	cmd         cmd
}

// NewHelp creates a Help using fields from the given Config. No validation is performed at this time.
func NewHelp(cfg env.Config) *Help {
	return &Help{
		HelmCommand: cfg.Command,
	}
}

// Execute executes the `helm help` command.
func (h *Help) Execute() error {
	if err := h.cmd.Run(); err != nil {
		return fmt.Errorf("while running '%s': %w", h.cmd.String(), err)
	}

	if h.HelmCommand == "help" {
		return nil
	}
	return fmt.Errorf("unknown command '%s'", h.HelmCommand)
}

// Prepare gets the Help ready to execute.
func (h *Help) Prepare(cfg Config) error {
	args := []string{"help"}
	if cfg.Debug {
		args = append([]string{"--debug"}, args...)
	}

	h.cmd = command(helmBin, args...)
	h.cmd.Stdout(cfg.Stdout)
	h.cmd.Stderr(cfg.Stderr)

	if cfg.Debug {
		fmt.Fprintf(cfg.Stderr, "Generated command: '%s'\n", h.cmd.String())
	}

	return nil
}

package run

import (
	"fmt"
)

// Help is a step in a helm Plan that calls `helm help`.
type Help struct {
	cmd cmd
}

// Execute executes the `helm help` command.
func (h *Help) Execute(cfg Config) error {
	if err := h.cmd.Run(); err != nil {
		return fmt.Errorf("while running '%s': %w", h.cmd.String(), err)
	}

	if cfg.HelmCommand == "help" {
		return nil
	}
	return fmt.Errorf("unknown command '%s'", cfg.HelmCommand)
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

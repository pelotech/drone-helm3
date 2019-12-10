package run

import (
	"fmt"
)

// Help is a step in a helm Plan that calls `helm help`.
type Help struct {
	cmd cmd
}

// Execute executes the `helm help` command.
func (h *Help) Execute() error {
	return h.cmd.Run()
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

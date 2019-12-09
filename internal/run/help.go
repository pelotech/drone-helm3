package run

import (
	"os"
)

type Help struct {
	cmd cmd
}

func (h *Help) Run() error {
	return h.cmd.Run()
}

func NewHelp() *Help {
	h := Help{}

	h.cmd = Command(HELM_BIN, "help")
	h.cmd.Stdout(os.Stdout)
	h.cmd.Stderr(os.Stderr)

	return &h
}

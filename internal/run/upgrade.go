package run

import (
	"os"
)

// Upgrade is a step in a helm Plan that calls `helm upgrade`.
type Upgrade struct {
	Chart   string
	Release string
	cmd     cmd
}

// Run launches the command.
func (u *Upgrade) Run() error {
	return u.cmd.Run()
}

// NewUpgrade creates a new Upgrade.
func NewUpgrade(release, chart string) *Upgrade {
	u := Upgrade{
		Chart:   chart,
		Release: release,
		cmd:     command(helmBin, "upgrade", "--install", release, chart),
	}

	u.cmd.Stdout(os.Stdout)
	u.cmd.Stderr(os.Stderr)

	return &u
}

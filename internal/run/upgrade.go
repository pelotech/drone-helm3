package run

import (
	"os"
)

type Upgrade struct {
	Chart   string
	Release string
	cmd     cmd
}

func (u *Upgrade) Run() error {
	return u.cmd.Run()
}

func NewUpgrade(release, chart string) *Upgrade {
	u := Upgrade{
		Chart:   chart,
		Release: release,
		cmd:     Command(HELM_BIN, "upgrade", "--install", release, chart),
	}

	u.cmd.Stdout(os.Stdout)
	u.cmd.Stderr(os.Stderr)

	return &u
}

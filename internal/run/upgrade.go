package run

import (
	"os"
)

func Upgrade(args ...string) error {
	args = append([]string{"upgrade"}, args...)
	cmd := Command(HELM_BIN, args...)

	cmd.Stdout(os.Stdout)
	cmd.Stderr(os.Stderr)

	return cmd.Run()
}

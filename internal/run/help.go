package run

import (
	"os"
)

func Help(args ...string) error {
	args = append([]string{"help"}, args...)

	cmd := Command(HELM_BIN, args...)
	cmd.Stdout(os.Stdout)
	cmd.Stderr(os.Stderr)

	return cmd.Run()
}

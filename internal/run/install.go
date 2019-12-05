package run

import ()

const HELM_BIN = "/usr/bin/helm"

func Install(args ...string) error {
	cmd := Command()
	cmd.Path(HELM_BIN)

	args = append([]string{"install"}, args...)
	cmd.Args(args)

	return cmd.Run()
}

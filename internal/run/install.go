package run

import ()

const HELM_BIN = "/usr/bin/helm"

func Install(args ...string) error {
	cmd := &execCmd{}
	cmd.Path(HELM_BIN)

	return install(cmd, args)
}

func install(cmd cmd, args []string) error {
	args = append([]string{"install"}, args...)
	cmd.Args(args)

	return cmd.Run()
}

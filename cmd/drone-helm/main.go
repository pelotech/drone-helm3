package main

import (
	"fmt"
	"github.com/urfave/cli"
	"os"

	"github.com/pelotech/drone-helm3/internal/run"
)

func main() {
	app := cli.NewApp()
	app.Name = "helm plugin"
	app.Usage = "helm plugin"
	app.Action = execute
	app.Version = "0.0.1α"
	app.Flags = []cli.Flag{
		cli.StringFlag{
			Name:   "helm_command",
			Usage:  "Helm command to execute",
			EnvVar: "PLUGIN_HELM_COMMAND,HELM_COMMAND",
		},
	}

	if err := app.Run(os.Args); err != nil {
		fmt.Fprintf(os.Stderr, "%s\n", err.Error())
	}
}


func execute(c *cli.Context) error {
	switch c.String("helm_command") {
	case "upgrade":
		run.Upgrade()
	case "help":
		run.Help()
	default:
		switch os.Getenv("DRONE_BUILD_EVENT") {
		case "push", "tag", "deployment", "pull_request", "promote", "rollback":
			run.Upgrade()
		default:
			run.Help()
	}
	return nil
}

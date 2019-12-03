package main

import (
	"fmt"
	"os"

	"github.com/urfave/cli"
)

func main() {
	_ = fmt.Println
	_ = os.Exit

	app := cli.NewApp()
	app.Name = "helm plugin"
	app.Usage = "helm plugin"
	app.Action = run
	app.Version = "0.0.1α"
	app.Flags = []cli.Flag{
		cli.StringFlag{
			Name:   "echo, e",
			Usage:  "this text'll be ech'll",
			EnvVar: "ECHO",
		},
	}

	if err := app.Run(os.Args); err != nil {
		fmt.Fprintf(os.Stderr, "%s\n", err.Error())
	}
}

func run(c *cli.Context) error {
	fmt.Println(c.String("echo"))
	return nil
}

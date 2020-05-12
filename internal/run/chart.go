package run

import (
	"fmt"

	"github.com/pelotech/drone-helm3/internal/env"
)

// Chart is an execution step that calls `helm chart` when executed.
type Chart struct {
	*config
	chartPath        string
	registryURL      string
	registryRepoName string
	chartVersion     string
	subCommand       string
	cmd              cmd
}

// NewChart creates a Chart using fields from the given Config. No validation is performed at this time.
func NewChart(subCommand string, cfg env.Config) *Chart {
	return &Chart{
		config:           newConfig(cfg),
		chartPath:        cfg.Chart,
		registryURL:      cfg.RegistryURL,
		registryRepoName: cfg.RegistryRepoName,
		chartVersion:     cfg.ChartVersion,
		subCommand:       subCommand,
	}
}

// Execute executes the `helm chart` command.
func (reg *Chart) Execute() error {
	return reg.cmd.Run()
}

// Prepare gets the Chart ready to execute.
func (reg *Chart) Prepare() error {
	args := []string{}

	args = append(args, "chart")

	if reg.subCommand == "save" {
		args = append(args, "save")
		args = append(args, reg.chartPath)
		cmd := fmt.Sprintf("%s/%s:%s", reg.registryURL, reg.registryRepoName, reg.chartVersion)
		args = append(args, cmd)
	} else if reg.subCommand == "push" {
		args = append(args, "push")
		cmd := fmt.Sprintf("%s/%s:%s", reg.registryURL, reg.registryRepoName, reg.chartVersion)
		args = append(args, cmd)
	}

	reg.cmd = command(helmBin, args...)
	reg.cmd.Stdout(reg.stdout)
	reg.cmd.Stderr(reg.stderr)

	if reg.debug {
		fmt.Fprintf(reg.stderr, "Generated command: '%s'\n", reg.cmd.String())
	}

	return nil
}

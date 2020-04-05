package run

import (
  "errors"
  "fmt"
  "github.com/pelotech/drone-helm3/internal/env"
)

const (
  actionBuild  = "build"
  actionUpdate = "update"
)

// DepAction is an execution step that calls `helm dependency update` or `helm dependency build` when executed.
type DepAction struct {
  *config
  chart  string
  cmd    cmd
  action string
}

// NewDepAction creates a DepAction using fields from the given Config. No validation is performed at this time.
func NewDepAction(cfg env.Config) *DepAction {
  return &DepAction{
    config: newConfig(cfg),
    chart:  cfg.Chart,
    action: cfg.DependenciesAction,
  }
}

// Execute executes the `helm upgrade` command.
func (d *DepAction) Execute() error {
  return d.cmd.Run()
}

// Prepare gets the DepAction ready to execute.
func (d *DepAction) Prepare() error {
  if d.chart == "" {
    return fmt.Errorf("chart is required")
  }

  args := d.globalFlags()

  if d.action != actionBuild && d.action != actionUpdate {
    return errors.New("unknown dependency_action: " + d.action)
  }

  args = append(args, "dependency", d.action, d.chart)

  d.cmd = command(helmBin, args...)
  d.cmd.Stdout(d.stdout)
  d.cmd.Stderr(d.stderr)

  if d.debug {
    fmt.Fprintf(d.stderr, "Generated command: '%s'\n", d.cmd.String())
  }

  return nil
}

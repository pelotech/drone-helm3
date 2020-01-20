package run

import (
	"github.com/pelotech/drone-helm3/internal/env"
	"io"
)

type config struct {
	debug     bool
	namespace string
	stdout    io.Writer
	stderr    io.Writer
}

func newConfig(cfg env.Config) *config {
	return &config{
		debug:     cfg.Debug,
		namespace: cfg.Namespace,
		stdout:    cfg.Stdout,
		stderr:    cfg.Stderr,
	}
}

func (cfg *config) globalFlags() []string {
	flags := []string{}
	if cfg.debug {
		flags = append(flags, "--debug")
	}
	if cfg.namespace != "" {
		flags = append(flags, "--namespace", cfg.namespace)
	}
	return flags
}

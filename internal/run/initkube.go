package run

import (
	"errors"
	"fmt"
	"github.com/pelotech/drone-helm3/internal/env"
	"io"
	"os"
	"text/template"
)

// InitKube is a step in a helm Plan that initializes the kubernetes config file.
type InitKube struct {
	SkipTLSVerify  bool
	Certificate    string
	APIServer      string
	ServiceAccount string
	Token          string
	TemplateFile   string
	ConfigFile     string

	template   *template.Template
	configFile io.WriteCloser
	values     kubeValues
}

type kubeValues struct {
	SkipTLSVerify  bool
	Certificate    string
	APIServer      string
	Namespace      string
	ServiceAccount string
	Token          string
}

// NewInitKube creates a InitKube using the given Config and filepaths. No validation is performed at this time.
func NewInitKube(cfg env.Config, templateFile, configFile string) *InitKube {
	return &InitKube{
		SkipTLSVerify:  cfg.SkipTLSVerify,
		Certificate:    cfg.Certificate,
		APIServer:      cfg.APIServer,
		ServiceAccount: cfg.ServiceAccount,
		Token:          cfg.KubeToken,
		TemplateFile:   templateFile,
		ConfigFile:     configFile,
	}
}

// Execute generates a kubernetes config file from drone-helm3's template.
func (i *InitKube) Execute(cfg Config) error {
	if cfg.Debug {
		fmt.Fprintf(cfg.Stderr, "writing kubeconfig file to %s\n", i.ConfigFile)
	}
	defer i.configFile.Close()
	return i.template.Execute(i.configFile, i.values)
}

// Prepare ensures all required configuration is present and that the config file is writable.
func (i *InitKube) Prepare(cfg Config) error {
	var err error

	if i.APIServer == "" {
		return errors.New("an API Server is needed to deploy")
	}
	if i.Token == "" {
		return errors.New("token is needed to deploy")
	}

	if i.ServiceAccount == "" {
		i.ServiceAccount = "helm"
	}

	if cfg.Debug {
		fmt.Fprintf(cfg.Stderr, "loading kubeconfig template from %s\n", i.TemplateFile)
	}
	i.template, err = template.ParseFiles(i.TemplateFile)
	if err != nil {
		return fmt.Errorf("could not load kubeconfig template: %w", err)
	}

	i.values = kubeValues{
		SkipTLSVerify:  i.SkipTLSVerify,
		Certificate:    i.Certificate,
		APIServer:      i.APIServer,
		ServiceAccount: i.ServiceAccount,
		Token:          i.Token,
		Namespace:      cfg.Namespace,
	}

	if cfg.Debug {
		if _, err := os.Stat(i.ConfigFile); err != nil {
			// non-nil err here isn't an actual error state; the kubeconfig just doesn't exist
			fmt.Fprint(cfg.Stderr, "creating ")
		} else {
			fmt.Fprint(cfg.Stderr, "truncating ")
		}
		fmt.Fprintf(cfg.Stderr, "kubeconfig file at %s\n", i.ConfigFile)
	}

	i.configFile, err = os.Create(i.ConfigFile)
	if err != nil {
		return fmt.Errorf("could not open kubeconfig file for writing: %w", err)
	}
	return nil
}

package helm

import (
	"fmt"
	"strings"
)

type Config struct {
	// Configuration for drone-helm itself
	Command            HelmCommand `envconfig:"HELM_COMMAND"`      // Helm command to run
	DroneEvent         string      `envconfig:"DRONE_BUILD_EVENT"` // Drone event that invoked this plugin.
	UpdateDependencies bool        `split_words:"true"`            // call `helm dependency update` before the main command
	Repos              []string    `envconfig:"HELM_REPOS"`        // call `helm repo add` before the main command
	Prefix             string      ``                              // Prefix to use when looking up secret env vars

	// Global helm config
	Debug          bool     ``                                                // global helm flag (also applies to drone-helm itself)
	KubeConfig     string   `split_words:"true" default:"/root/.kube/config"` // path to the kube config file
	Values         string   ``
	StringValues   string   `split_words:"true"`
	ValuesFiles    []string `split_words:"true"`
	Namespace      string   ``
	Token          string   `envconfig:"KUBERNETES_TOKEN"`
	SkipTLSVerify  bool     `envconfig:"SKIP_TLS_VERIFY"`
	Certificate    string   `envconfig:"KUBERNETES_CERTIFICATE"`
	APIServer      string   `envconfig:"API_SERVER"`
	ServiceAccount string   `envconfig:"SERVICE_ACCOUNT"` // Can't just use split_words; need envconfig to find the non-prefixed form

	// Config specifically for `helm upgrade`
	ChartVersion string `split_words:"true"` //
	DryRun       bool   `split_words:"true"` // also available for `delete`
	Wait         bool   ``                   //
	ReuseValues  bool   `split_words:"true"` //
	Timeout      string ``                   //
	Chart        string ``                   // Also available for `lint`, in which case it must be a path to a chart directory
	Release      string ``
	Force        bool   `` //
}

type HelmCommand string

// HelmCommand.Decode checks the given value against the list of known commands and generates a helpful error if the command is unknown.
func (cmd *HelmCommand) Decode(value string) error {
	known := []string{"upgrade", "delete", "lint", "help"}
	for _, c := range known {
		if value == c {
			*cmd = HelmCommand(value)
			return nil
		}
	}

	if value == "" {
		return nil
	}
	known[len(known)-1] = fmt.Sprintf("or %s", known[len(known)-1])
	return fmt.Errorf("Unknown command '%s'. If specified, command must be %s.",
		value, strings.Join(known, ", "))
}

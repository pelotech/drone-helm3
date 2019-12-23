package helm

import (
	"fmt"
	"strings"
)

// The Config struct captures the `settings` and `environment` blocks in the application's drone
// config. Configuration in drone's `settings` block arrives as uppercase env vars matching the
// config key, prefixed with `PLUGIN_`. Config from the `environment` block is *not* prefixed; any
// keys that are likely to be in that block (i.e. things that use `from_secret` need an explicit
// `envconfig:` tag so that envconfig will look for a non-prefixed env var.
type Config struct {
	// Configuration for drone-helm itself
	Command            helmCommand `envconfig:"HELM_COMMAND"`           // Helm command to run
	DroneEvent         string      `envconfig:"DRONE_BUILD_EVENT"`      // Drone event that invoked this plugin.
	UpdateDependencies bool        `split_words:"true"`                 // Call `helm dependency update` before the main command
	Repos              []string    `envconfig:"HELM_REPOS"`             // Call `helm repo add` before the main command
	Prefix             string      ``                                   // Prefix to use when looking up secret env vars
	Debug              bool        ``                                   // Generate debug output and pass --debug to all helm commands
	Values             string      ``                                   // Argument to pass to --set in applicable helm commands
	StringValues       string      `split_words:"true"`                 // Argument to pass to --set-string in applicable helm commands
	ValuesFiles        []string    `split_words:"true"`                 // Arguments to pass to --values in applicable helm commands
	Namespace          string      ``                                   // Kubernetes namespace for all helm commands
	KubeToken          string      `envconfig:"KUBERNETES_TOKEN"`       // Kubernetes authentication token to put in .kube/config
	SkipTLSVerify      bool        `envconfig:"SKIP_TLS_VERIFY"`        // Put insecure-skip-tls-verify in .kube/config
	Certificate        string      `envconfig:"KUBERNETES_CERTIFICATE"` // The Kubernetes cluster CA's self-signed certificate (must be base64-encoded)
	APIServer          string      `envconfig:"API_SERVER"`             // The Kubernetes cluster's API endpoint
	ServiceAccount     string      `envconfig:"SERVICE_ACCOUNT"`        // Account to use for connecting to the Kubernetes cluster
	ChartVersion       string      `split_words:"true"`                 // Specific chart version to use in `helm upgrade`
	DryRun             bool        `split_words:"true"`                 // Pass --dry-run to applicable helm commands
	Wait               bool        ``                                   // Pass --wait to applicable helm commands
	ReuseValues        bool        `split_words:"true"`                 // Pass --reuse-values to `helm upgrade`
	Timeout            string      ``                                   // Argument to pass to --timeout in applicable helm commands
	Chart              string      ``                                   // Chart argument to use in applicable helm commands
	Release            string      ``                                   // Release argument to use in applicable helm commands
	Force              bool        ``                                   // Pass --force to applicable helm commands
}

type helmCommand string

// helmCommand.Decode checks the given value against the list of known commands and generates a helpful error if the command is unknown.
func (cmd *helmCommand) Decode(value string) error {
	known := []string{"upgrade", "delete", "lint", "help"}
	for _, c := range known {
		if value == c {
			*cmd = helmCommand(value)
			return nil
		}
	}

	if value == "" {
		return nil
	}
	known[len(known)-1] = fmt.Sprintf("or %s", known[len(known)-1])
	return fmt.Errorf("unknown command '%s'. If specified, command must be %s",
		value, strings.Join(known, ", "))
}

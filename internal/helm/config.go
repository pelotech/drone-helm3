package helm

import (
	"fmt"
	"github.com/kelseyhightower/envconfig"
	"io"
	"os"
	"regexp"
	"strings"
)

var (
	justNumbers    = regexp.MustCompile(`^\d+$`)
	deprecatedVars = []string{"PURGE", "RECREATE_PODS", "TILLER_NS", "UPGRADE", "CANARY_IMAGE", "CLIENT_ONLY", "STABLE_REPO_URL"}
)

// The Config struct captures the `settings` and `environment` blocks in the application's drone
// config. Configuration in drone's `settings` block arrives as uppercase env vars matching the
// config key, prefixed with `PLUGIN_`. Config from the `environment` block is uppercased, but does
// not have the `PLUGIN_` prefix.
type Config struct {
	// Configuration for drone-helm itself
	Command            string   `envconfig:"mode"`                   // Helm command to run
	DroneEvent         string   `envconfig:"DRONE_BUILD_EVENT"`      // Drone event that invoked this plugin.
	UpdateDependencies bool     `split_words:"true"`                 // Call `helm dependency update` before the main command
	AddRepos           []string `split_words:"true"`                 // Call `helm repo add` before the main command
	RepoCAFile         string   `envconfig:"repo_ca_file"`           // CA certificate for `helm repo add`
	Debug              bool     ``                                   // Generate debug output and pass --debug to all helm commands
	Values             string   ``                                   // Argument to pass to --set in applicable helm commands
	StringValues       string   `split_words:"true"`                 // Argument to pass to --set-string in applicable helm commands
	ValuesFiles        []string `split_words:"true"`                 // Arguments to pass to --values in applicable helm commands
	Namespace          string   ``                                   // Kubernetes namespace for all helm commands
	KubeToken          string   `split_words:"true"`                 // Kubernetes authentication token to put in .kube/config
	SkipTLSVerify      bool     `envconfig:"SKIP_TLS_VERIFY"`        // Put insecure-skip-tls-verify in .kube/config
	Certificate        string   `envconfig:"kube_certificate"`       // The Kubernetes cluster CA's self-signed certificate (must be base64-encoded)
	APIServer          string   `envconfig:"kube_api_server"`        // The Kubernetes cluster's API endpoint
	ServiceAccount     string   `envconfig:"kube_service_account"`   // Account to use for connecting to the Kubernetes cluster
	ChartVersion       string   `split_words:"true"`                 // Specific chart version to use in `helm upgrade`
	DryRun             bool     `split_words:"true"`                 // Pass --dry-run to applicable helm commands
	Wait               bool     `envconfig:"wait_for_upgrade"`       // Pass --wait to applicable helm commands
	ReuseValues        bool     `split_words:"true"`                 // Pass --reuse-values to `helm upgrade`
	KeepHistory        bool     `split_words:"true"`                 // Pass --keep-history to `helm uninstall`
	Timeout            string   ``                                   // Argument to pass to --timeout in applicable helm commands
	Chart              string   ``                                   // Chart argument to use in applicable helm commands
	Release            string   ``                                   // Release argument to use in applicable helm commands
	Force              bool     `envconfig:"force_upgrade"`          // Pass --force to applicable helm commands
	AtomicUpgrade      bool     `split_words:"true"`                 // Pass --atomic to `helm upgrade`
	CleanupOnFail      bool     `envconfig:"CLEANUP_FAILED_UPGRADE"` // Pass --cleanup-on-fail to `helm upgrade`
	LintStrictly       bool     `split_words:"true"`                 // Pass --strict to `helm lint`

	Stdout io.Writer `ignored:"true"`
	Stderr io.Writer `ignored:"true"`
}

// NewConfig creates a Config and reads environment variables into it, accounting for several possible formats.
func NewConfig(stdout, stderr io.Writer) (*Config, error) {
	var aliases settingAliases
	if err := envconfig.Process("plugin", &aliases); err != nil {
		return nil, err
	}

	if err := envconfig.Process("", &aliases); err != nil {
		return nil, err
	}

	cfg := Config{
		Command:        aliases.Command,
		AddRepos:       aliases.AddRepos,
		APIServer:      aliases.APIServer,
		ServiceAccount: aliases.ServiceAccount,
		Wait:           aliases.Wait,
		Force:          aliases.Force,
		KubeToken:      aliases.KubeToken,
		Certificate:    aliases.Certificate,

		Stdout: stdout,
		Stderr: stderr,
	}
	if err := envconfig.Process("plugin", &cfg); err != nil {
		return nil, err
	}

	if err := envconfig.Process("", &cfg); err != nil {
		return nil, err
	}

	if justNumbers.MatchString(cfg.Timeout) {
		cfg.Timeout = fmt.Sprintf("%ss", cfg.Timeout)
	}

	if cfg.Debug && cfg.Stderr != nil {
		cfg.logDebug()
	}

	cfg.deprecationWarn()

	return &cfg, nil
}

func (cfg Config) logDebug() {
	if cfg.KubeToken != "" {
		cfg.KubeToken = "(redacted)"
	}
	fmt.Fprintf(cfg.Stderr, "Generated config: %+v\n", cfg)
}

func (cfg *Config) deprecationWarn() {
	for _, varname := range deprecatedVars {
		_, barePresent := os.LookupEnv(varname)
		_, prefixedPresent := os.LookupEnv("PLUGIN_" + varname)
		if barePresent || prefixedPresent {
			fmt.Fprintf(cfg.Stderr, "Warning: ignoring deprecated '%s' setting\n", strings.ToLower(varname))
		}
	}
}

type settingAliases struct {
	Command        string   `envconfig:"helm_command"`
	AddRepos       []string `envconfig:"helm_repos"`
	APIServer      string   `envconfig:"api_server"`
	ServiceAccount string   `split_words:"true"`
	Wait           bool     ``
	Force          bool     ``
	KubeToken      string   `envconfig:"kubernetes_token"`
	Certificate    string   `envconfig:"kubernetes_certificate"`
}

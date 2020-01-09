package helm

import (
	"fmt"
	"github.com/urfave/cli/v2"
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
	Command            string   // Helm command to run
	DroneEvent         string   // Drone event that invoked this plugin.
	UpdateDependencies bool     // Call `helm dependency update` before the main command
	AddRepos           []string // Call `helm repo add` before the main command
	Debug              bool     // Generate debug output and pass --debug to all helm commands
	Values             string   // Argument to pass to --set in applicable helm commands
	StringValues       string   // Argument to pass to --set-string in applicable helm commands
	ValuesFiles        []string // Arguments to pass to --values in applicable helm commands
	Namespace          string   // Kubernetes namespace for all helm commands
	KubeToken          string   // Kubernetes authentication token to put in .kube/config
	SkipTLSVerify      bool     // Put insecure-skip-tls-verify in .kube/config
	Certificate        string   // The Kubernetes cluster CA's self-signed certificate (must be base64-encoded)
	APIServer          string   // The Kubernetes cluster's API endpoint
	ServiceAccount     string   // Account to use for connecting to the Kubernetes cluster
	ChartVersion       string   // Specific chart version to use in `helm upgrade`
	DryRun             bool     // Pass --dry-run to applicable helm commands
	Wait               bool     // Pass --wait to applicable helm commands
	ReuseValues        bool     // Pass --reuse-values to `helm upgrade`
	KeepHistory        bool     // Pass --keep-history to `helm uninstall`
	Timeout            string   // Argument to pass to --timeout in applicable helm commands
	Chart              string   // Chart argument to use in applicable helm commands
	Release            string   // Release argument to use in applicable helm commands
	Force              bool     // Pass --force to applicable helm commands
	AtomicUpgrade      bool     // Pass --atomic to `helm upgrade`
	CleanupOnFail      bool     // Pass --cleanup-on-fail to `helm upgrade`
	LintStrictly       bool     // Pass --strict to `helm lint`

	Stdout io.Writer
	Stderr io.Writer
}

// NewConfig creates a Config and reads environment variables into it, accounting for several possible formats.
func NewConfig(stdout, stderr io.Writer, argv ...string) (*Config, error) {
	cfg := Config{
		Stdout: stdout,
		Stderr: stderr,
	}
	// cli doesn't support Destination for string slices, so we'll use bare
	// strings as an intermediate value and split them on commas ourselves.
	var addRepos, valuesFiles string
	app := &cli.App{
		Name:   "drone-helm3",
		Action: func(*cli.Context) error { return nil },
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:        "mode",
				Destination: &cfg.Command,
				EnvVars:     []string{"MODE", "PLUGIN_MODE", "HELM_COMMAND", "PLUGIN_HELM_COMMAND"},
			},
			&cli.StringFlag{
				Name:        "drone-event",
				Destination: &cfg.DroneEvent,
				EnvVars:     []string{"DRONE_BUILD_EVENT"},
			},
			&cli.BoolFlag{
				Name:        "update-dependencies",
				Destination: &cfg.UpdateDependencies,
				EnvVars:     []string{"UPDATE_DEPENDENCIES", "PLUGIN_UPDATE_DEPENDENCIES"},
			},
			&cli.StringFlag{
				Name:        "add-repos",
				Destination: &addRepos,
				EnvVars:     []string{"ADD_REPOS", "PLUGIN_ADD_REPOS", "HELM_REPOS", "PLUGIN_HELM_REPOS"},
			},
			&cli.BoolFlag{
				Name:        "debug",
				Destination: &cfg.Debug,
				EnvVars:     []string{"DEBUG", "PLUGIN_DEBUG"},
			},
			&cli.StringFlag{
				Name:        "values",
				Destination: &cfg.Values,
				EnvVars:     []string{"VALUES", "PLUGIN_VALUES"},
			},
			&cli.StringFlag{
				Name:        "string-values",
				Destination: &cfg.StringValues,
				EnvVars:     []string{"STRING_VALUES", "PLUGIN_STRING_VALUES"},
			},
			&cli.StringFlag{
				Name:        "values-files",
				Destination: &valuesFiles,
				EnvVars:     []string{"VALUES_FILES", "PLUGIN_VALUES_FILES"},
			},
			&cli.StringFlag{
				Name:        "namespace",
				Destination: &cfg.Namespace,
				EnvVars:     []string{"NAMESPACE", "PLUGIN_NAMESPACE"},
			},
			&cli.StringFlag{
				Name:        "kube-token",
				Destination: &cfg.KubeToken,
				EnvVars:     []string{"KUBE_TOKEN", "PLUGIN_KUBE_TOKEN", "KUBERNETES_TOKEN", "PLUGIN_KUBERNETES_TOKEN"},
			},
			&cli.BoolFlag{
				Name:        "skip-tls-verify",
				Destination: &cfg.SkipTLSVerify,
				EnvVars:     []string{"SKIP_TLS_VERIFY", "PLUGIN_SKIP_TLS_VERIFY"},
			},
			&cli.StringFlag{
				Name:        "kube-certificate",
				Destination: &cfg.Certificate,
				EnvVars:     []string{"KUBE_CERTIFICATE", "PLUGIN_KUBE_CERTIFICATE", "KUBERNETES_CERTIFICATE", "PLUGIN_KUBERNETES_CERTIFICATE"},
			},
			&cli.StringFlag{
				Name:        "kube-api-server",
				Destination: &cfg.APIServer,
				EnvVars:     []string{"KUBE_API_SERVER", "PLUGIN_KUBE_API_SERVER", "API_SERVER", "PLUGIN_API_SERVER"},
			},
			&cli.StringFlag{
				Name:        "service-account",
				Destination: &cfg.ServiceAccount,
				EnvVars:     []string{"KUBE_SERVICE_ACCOUNT", "PLUGIN_KUBE_SERVICE_ACCOUNT", "SERVICE_ACCOUNT", "PLUGIN_SERVICE_ACCOUNT"},
			},
			&cli.StringFlag{
				Name:        "chart-version",
				Destination: &cfg.ChartVersion,
				EnvVars:     []string{"CHART_VERSION", "PLUGIN_CHART_VERSION"},
			},
			&cli.BoolFlag{
				Name:        "dry-run",
				Destination: &cfg.DryRun,
				EnvVars:     []string{"DRY_RUN", "PLUGIN_DRY_RUN"},
			},
			&cli.BoolFlag{
				Name:        "wait-for-upgrade",
				Destination: &cfg.Wait,
				EnvVars:     []string{"WAIT_FOR_UPGRADE", "PLUGIN_WAIT_FOR_UPGRADE", "WAIT", "PLUGIN_WAIT"},
			},
			&cli.BoolFlag{
				Name:        "reuse-values",
				Destination: &cfg.ReuseValues,
				EnvVars:     []string{"REUSE_VALUES", "PLUGIN_REUSE_VALUES"},
			},
			&cli.BoolFlag{
				Name:        "keep-history",
				Destination: &cfg.KeepHistory,
				EnvVars:     []string{"KEEP_HISTORY", "PLUGIN_KEEP_HISTORY"},
			},
			&cli.StringFlag{
				Name:        "timeout",
				Destination: &cfg.Timeout,
				EnvVars:     []string{"TIMEOUT", "PLUGIN_TIMEOUT"},
			},
			&cli.StringFlag{
				Name:        "chart",
				Destination: &cfg.Chart,
				EnvVars:     []string{"CHART", "PLUGIN_CHART"},
			},
			&cli.StringFlag{
				Name:        "release",
				Destination: &cfg.Release,
				EnvVars:     []string{"RELEASE", "PLUGIN_RELEASE"},
			},
			&cli.BoolFlag{
				Name:        "force-upgrade",
				Destination: &cfg.Force,
				EnvVars:     []string{"FORCE_UPGRADE", "PLUGIN_FORCE_UPGRADE", "FORCE", "PLUGIN_FORCE"},
			},
			&cli.BoolFlag{
				Name:        "atomic-upgrade",
				Destination: &cfg.AtomicUpgrade,
				EnvVars:     []string{"ATOMIC_UPGRADE", "PLUGIN_ATOMIC_UPGRADE"},
			},
			&cli.BoolFlag{
				Name:        "cleanup-failed-upgrade",
				Destination: &cfg.CleanupOnFail,
				EnvVars:     []string{"CLEANUP_FAILED_UPGRADE", "PLUGIN_CLEANUP_FAILED_UPGRADE"},
			},
			&cli.BoolFlag{
				Name:        "lint-strictly",
				Destination: &cfg.LintStrictly,
				EnvVars:     []string{"LINT_STRICTLY", "PLUGIN_LINT_STRICTLY"},
			},
		},
	}
	if err := app.Run(argv); err != nil {
		return nil, err
	}
	if addRepos != "" {
		cfg.AddRepos = strings.Split(addRepos, ",")
	}
	if valuesFiles != "" {
		cfg.ValuesFiles = strings.Split(valuesFiles, ",")
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

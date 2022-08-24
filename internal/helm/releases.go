package helm

import (
	"fmt"
	"github.com/pelotech/drone-helm3/internal/env"
	"github.com/pelotech/drone-helm3/internal/run"
	"helm.sh/helm/v3/pkg/action"
	"helm.sh/helm/v3/pkg/cli"
	"log"
	"os"
)

const (
	kubeConfigTemplate = "/root/.kube/config.tpl"
	kubeConfigFile     = "/root/.kube/config"
)

type Release struct {
	Name      string
	Namespace string
}

func DetermineReleases(cfg env.Config) ([]Release, error) {
	var releases []Release

	if !cfg.SkipKubeconfig {
		initKube := run.NewInitKube(cfg, kubeConfigTemplate, kubeConfigFile)
		if cfg.Debug {
			fmt.Fprintf(cfg.Stderr, "calling %T. \n", initKube)
		}
		if err := initKube.Execute(); err != nil {
			err = fmt.Errorf("while during the kubernetes initial configuration: %w", err)
			return nil, err
		}
	}

	if cfg.ChartSelector != "" {
		settings := cli.New()
		actionConfig := new(action.Configuration)

		if err := actionConfig.Init(settings.RESTClientGetter(), "", os.Getenv("HELM_DRIVER"), log.Printf); err != nil {
			err = fmt.Errorf("while executing helm configuration initialization: %w", err)
			return nil, err
		}

		client := action.NewList(actionConfig)
		client.AllNamespaces = true
		client.Deployed = true
		client.Filter = fmt.Sprintf("^%s-[0-9]*.[0-9]*.[0-9]*", client.Filter)

		results, err := client.Run()

		if err != nil {
			err = fmt.Errorf("while executing helm releases list: %w", err)
			return nil, err
		}

		ignoredReleasesMap := make(map[string]string)
		for _, s := range cfg.IgnoreReleases {
			ignoredReleasesMap[s] = ""
		}

		for _, rel := range results {
			// TODO: Check the relase chart name
			if _, ignore := ignoredReleasesMap[rel.Name]; ignore {
				continue
			}
			r := Release{
				Name:      rel.Name,
				Namespace: rel.Namespace,
			}
			releases = append(releases, r)
		}
	} else {
		r := Release{
			Name:      cfg.Release,
			Namespace: cfg.Namespace,
		}
		releases = append(releases, r)
	}

	return releases, nil
}

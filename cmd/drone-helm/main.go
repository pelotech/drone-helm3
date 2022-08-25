package main

import (
	"fmt"
	"os"

	_ "github.com/joho/godotenv/autoload"
	"github.com/pelotech/drone-helm3/internal/env"
	"github.com/pelotech/drone-helm3/internal/helm"
)

func main() {
	cfg, err := env.NewConfig(os.Stdout, os.Stderr)

	if err != nil {
		fmt.Fprintf(os.Stderr, "%s\n", err.Error())
		return
	}

	releases, err := helm.DetermineReleases(*cfg)
	if err != nil {
		fmt.Fprintf(os.Stderr, "%s\n", err.Error())
		os.Exit(1)
	}

	for _, rel := range releases {
		// Make a new plan for each matched release
		cfg.Release = rel.Name
		cfg.Namespace = rel.Namespace

		plan, err := helm.NewPlan(*cfg)
		if err != nil {
			fmt.Fprintf(os.Stderr, "%s\n", err.Error())
			os.Exit(1)
		}

		// Execute the plan
		err = plan.Execute()

		// Expect the plan to go off the rails
		if err != nil {
			fmt.Fprintf(os.Stderr, "%s\n", err.Error())
			// Throw away the plan
			os.Exit(1)
		}
	}

}

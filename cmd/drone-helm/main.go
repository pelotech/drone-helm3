package main

import (
	"fmt"
	"os"

	"github.com/pelotech/drone-helm3/internal/helm"
)

func main() {
	cfg, err := helm.NewConfig(os.Stdout, os.Stderr)

	if err != nil {
		fmt.Fprintf(os.Stderr, "%s\n", err.Error())
		return
	}

	// Make the plan
	plan, err := helm.NewPlan(*cfg)
	if err != nil {
		fmt.Fprintf(os.Stderr, "%w\n", err)
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

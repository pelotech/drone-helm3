package main

import (
	"fmt"
	"github.com/kelseyhightower/envconfig"
	"os"

	"github.com/pelotech/drone-helm3/internal/helm"
)

func main() {
	var c helm.Config

	if err := envconfig.Process("plugin", &c); err != nil {
		fmt.Fprintf(os.Stderr, "%s\n", err.Error())
		return
	}

	// Make the plan
	plan, err := helm.NewPlan(c)
	if err != nil {
		fmt.Fprintf(os.Stderr, "%w\n", err)
		os.Exit(1)
	}

	// Execute the plan
	err = plan.Execute()

	// Expect the plan to go off the rails
	if err != nil {
		fmt.Fprintf(os.Stderr, err.Error())
		// Throw away the plan
		os.Exit(1)
	}
}

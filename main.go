package main

import (
	"log"

	"github.com/lakhansamani/container-orchestrator/cmd"
)

var version string

func main() {
	cmd.SetVersion(version)
	if err := cmd.RootCmd.Execute(); err != nil {
		log.Fatalf("%v\n", err)
	}
}

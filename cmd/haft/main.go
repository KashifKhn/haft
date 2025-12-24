package main

import (
	"os"

	"github.com/KashifKhn/haft/internal/cli/root"
)

var version = "dev"

func main() {
	root.SetVersion(version)
	if err := root.Execute(); err != nil {
		os.Exit(1)
	}
}

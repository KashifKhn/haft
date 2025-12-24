package main

import (
	"os"

	"github.com/KashifKhn/haft/internal/cli/root"
)

func main() {
	if err := root.Execute(); err != nil {
		os.Exit(1)
	}
}

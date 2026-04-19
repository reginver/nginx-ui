package main

import (
	"fmt"
	"os"

	"github.com/0xJacky/nginx-ui/cmd"
)

// Version information set by build flags
var (
	Version   = "dev"
	BuildTime = "unknown"
	GitCommit = "unknown"
)

func main() {
	// Inject version info into cmd package
	cmd.Version = Version
	cmd.BuildTime = BuildTime
	cmd.GitCommit = GitCommit

	if err := cmd.Execute(); err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}
}

package main

import (
	"github.com/spf13/pflag"
	"kc/pkg/cmd"
	"os"
)

func main() {
	flags := pflag.NewFlagSet("kubectl-ns", pflag.ExitOnError)
	pflag.CommandLine = flags

	root := cmd.NewCmdContextSwitcher()
	if err := root.Execute(); err != nil {
		os.Exit(1)
	}
}

package cmd

import (
	"fmt"
	cobalias "github.com/spf13/cobra"
	"runtime"
)

var (
	Version   = "v1.0.0"
	Commit    = "placeholder" // set via GH Actions WF
	BuildDate = "placeholder" // set via GH Actions WF
)

var versionCmd = &cobalias.Command{
	Use:   "version",
	Short: "Prints CLI version",
	Long:  "Usage: `executable version`",
	Run: func(cmd *cobalias.Command, args []string) {
		fmt.Printf("kctl version: %s\n", Version)
		fmt.Printf("Commit: %s\n", Commit)
		fmt.Printf("BuildDate: %s\n", BuildDate)
		fmt.Printf("Go version: %s\n", runtime.Version())
	},
}

func init() {
	rootCmd.AddCommand(versionCmd)
}

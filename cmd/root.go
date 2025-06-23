package cmd

import (
	"fmt"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	cobalias "github.com/spf13/cobra"
)

var rootCmd = &cobalias.Command{
	Use:   "kctl",
	Short: "MyCLI is a custom Kubernetes controller CLI",
	Long:  `MyCLI is a tool to test and run components of your custom Kubernetes controller.`,
	Run: func(cmd *cobalias.Command, args []string) {
		zerolog.TimeFieldFormat = zerolog.TimeFormatUnix
		zerolog.TimeFieldFormat = "2006-01-02 15:04:05"
		log.Info().Msgf("KCTL version: %s", Version)
	},
}

func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		fmt.Println("Execution failed: ", err)
	}
}

//func init() {
//	rootCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
//}

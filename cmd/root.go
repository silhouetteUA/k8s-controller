package cmd

import (
	"fmt"
	cobalias "github.com/spf13/cobra"
)

var rootCmd = &cobalias.Command{
	Use:   "kctl",
	Short: "MyCLI is a custom Kubernetes controller CLI",
	Long:  `MyCLI is a tool to test and run components of your custom Kubernetes controller.`,
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

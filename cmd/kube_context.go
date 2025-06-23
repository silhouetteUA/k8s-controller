package cmd

import (
	"fmt"
	"github.com/rs/zerolog/log"
	"os"

	"github.com/spf13/cobra"
	"k8s.io/client-go/tools/clientcmd"
)

var contextCmd = &cobra.Command{
	Use:   "context",
	Short: "Show current Kubernetes context and cluster info",
	Run: func(cmd *cobra.Command, args []string) {
		kubeconfig := os.Getenv("KUBECONFIG")
		if kubeconfig == "" {
			home, err := os.UserHomeDir()
			if err != nil {
				fmt.Println("Error getting home directory:", err)
				log.Error().Err(err).Msg("Error getting home directory")
				return
			}
			kubeconfig = home + "/.kube/config"
		}

		config, err := clientcmd.LoadFromFile(kubeconfig)
		if err != nil {
			fmt.Println("Failed to load kubeconfig:", err)
			return
		}

		currentContext := config.CurrentContext
		fmt.Println("Current Context:", currentContext)

		ctxConfig := config.Contexts[currentContext]
		if ctxConfig == nil {
			fmt.Println("Context config not found")
			log.Error().Msg("Context config not found")
			return
		}

		fmt.Println("Cluster:", ctxConfig.Cluster)
		fmt.Println("Namespace:", ctxConfig.Namespace)
	},
}

func init() {
	rootCmd.AddCommand(contextCmd)
}

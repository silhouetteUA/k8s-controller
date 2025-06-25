package cmd

import (
	"context"
	"fmt"
	"os"

	"github.com/rs/zerolog/log"
	"github.com/spf13/cobra"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

var deploymentName string

var deleteCmd = &cobra.Command{
	Use:   "delete",
	Short: "Delete a Kubernetes deployment in the specified namespace",
	Run: func(cmd *cobra.Command, args []string) {
		if deploymentName == "" {
			log.Error().Msg("Deployment name must be provided with --name flag")
			os.Exit(1)
		}

		clientset, err := getKubeClient(kubeconfig)
		if err != nil {
			log.Error().Err(err).Msg("Failed to create Kubernetes client")
			os.Exit(1)
		}

		err = clientset.AppsV1().Deployments(namespace).Delete(
			context.Background(),
			deploymentName,
			metav1.DeleteOptions{},
		)
		if err != nil {
			log.Error().Err(err).Msgf("Failed to delete deployment: %s", deploymentName)
			os.Exit(1)
		}

		fmt.Printf("Deployment %q deleted successfully from namespace %q\n", deploymentName, namespace)
	},
}

func init() {
	rootCmd.AddCommand(deleteCmd)
	deleteCmd.Flags().StringVar(&kubeconfig, "kubeconfig", "", "Path to the kubeconfig file")
	deleteCmd.Flags().StringVar(&namespace, "namespace", "default", "Kubernetes namespace")
	deleteCmd.Flags().StringVar(&deploymentName, "name", "", "Name of the deployment to delete")
}

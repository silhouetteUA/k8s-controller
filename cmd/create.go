package cmd

import (
	"context"
	"fmt"
	"os"

	"github.com/rs/zerolog/log"
	"github.com/spf13/cobra"
	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

var deploymentImage string

var createCmd = &cobra.Command{
	Use:   "create",
	Short: "Create a simple Kubernetes deployment",
	Run: func(cmd *cobra.Command, args []string) {
		if deploymentName == "" || deploymentImage == "" {
			log.Error().Msg("Both --name and --image are required")
			os.Exit(1)
		}

		clientset, err := getKubeClient(kubeconfig)
		if err != nil {
			log.Error().Err(err).Msg("Failed to create Kubernetes client")
			os.Exit(1)
		}

		deployment := &appsv1.Deployment{
			ObjectMeta: metav1.ObjectMeta{
				Name: deploymentName,
			},
			Spec: appsv1.DeploymentSpec{
				Replicas: int32Ptr(1),
				Selector: &metav1.LabelSelector{
					MatchLabels: map[string]string{"app": deploymentName},
				},
				Template: corev1.PodTemplateSpec{
					ObjectMeta: metav1.ObjectMeta{
						Labels: map[string]string{"app": deploymentName},
					},
					Spec: corev1.PodSpec{
						Containers: []corev1.Container{
							{
								Name:  deploymentName,
								Image: deploymentImage,
							},
						},
					},
				},
			},
		}

		_, err = clientset.AppsV1().Deployments(namespace).Create(context.Background(), deployment, metav1.CreateOptions{})
		if err != nil {
			log.Error().Err(err).Msg("Failed to create deployment")
			os.Exit(1)
		}

		fmt.Printf("Deployment %q with image %q created in namespace %q\n", deploymentName, deploymentImage, namespace)
	},
}

func init() {
	rootCmd.AddCommand(createCmd)
	createCmd.Flags().StringVar(&kubeconfig, "kubeconfig", "", "Path to the kubeconfig file")
	createCmd.Flags().StringVar(&namespace, "namespace", "default", "Kubernetes namespace")
	createCmd.Flags().StringVar(&deploymentName, "name", "", "Name of the deployment to create")
	createCmd.Flags().StringVar(&deploymentImage, "image", "", "Container image for the deployment")
}

func int32Ptr(i int32) *int32 {
	return &i
}

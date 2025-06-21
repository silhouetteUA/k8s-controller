package cmd

import (
	"crypto/tls"
	"fmt"
	"github.com/spf13/cobra"
	"net/http"
	"time"
)

var apiServerURL string

var healthCmd = &cobra.Command{
	Use:   "health",
	Short: "Perform a health check of the Kubernetes API",
	Long:  "Use flag `api-server=value` to specify the Kubernetes API Server URL and execute health check against it",
	Run: func(cmd *cobra.Command, args []string) {
		//client := http.Client{Timeout: 3 * time.Second}
		client := &http.Client{
			Timeout: 3 * time.Second,
			Transport: &http.Transport{
				TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
			},
		}
		resp, err := client.Get(apiServerURL + "/healthz")
		if err != nil {
			fmt.Println("Status: Unhealthy")
			fmt.Printf("K8s API reachable: no (%v)\n", err)
			return
		}
		defer resp.Body.Close()
		if resp.StatusCode == 200 {
			fmt.Println("Status: Healthy")
			fmt.Println("K8s API reachable: yes")
		} else {
			fmt.Println("Status: Unhealthy")
			fmt.Printf("K8s API reachable: no (status code %d)\n", resp.StatusCode)
		}
	},
}

func init() {
	rootCmd.AddCommand(healthCmd)
	healthCmd.Flags().StringVar(&apiServerURL, "api-server", "https://kubernetes.default.svc", "Health check external Kubernetes API server")
}

package cmd

import (
	"context"
	"fmt"
	"github.com/google/uuid"
	"github.com/rs/zerolog/log"
	"github.com/silhouetteUA/k8s-controller/pkg/informer"
	"github.com/spf13/cobra"
	"github.com/valyala/fasthttp"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"
	"os"
)

var serverPort int
var serverKubeconfig string
var serverInCluster bool

var serverCmd = &cobra.Command{
	Use:   "server",
	Short: "Start a FastHTTP server",
	Run: func(cmd *cobra.Command, args []string) {
		clientset, err := getServerKubeClient(serverKubeconfig, serverInCluster)
		if err != nil {
			log.Error().Err(err).Msg("Failed to create Kubernetes client")
			os.Exit(1)
		}
		ctx := context.Background()
		go informer.StartInformerFactory(ctx, clientset, namespace)
		handler := func(ctx *fasthttp.RequestCtx) {
			uuid := uuid.New().String()
			switch string(ctx.Path()) {
			case "/version":
				log.Info().
					Str("request_id", uuid).
					Str("method", string(ctx.Method())).
					Str("path", string(ctx.Path())).
					Str("remote_addr", ctx.RemoteAddr().String()).
					Msg("Check version request")
				ctx.Response.Header.SetContentType("application/json")
				ctx.Response.Header.Set("X-Request-ID", uuid)
				_, err := fmt.Fprintf(ctx, `{"version": "%s", "commit": "%s", "date": "%s", "requestID": "%s"}`, Version, Commit, BuildDate, uuid)
				if err != nil {
					return
				}
			case "/deployments":
				log.Info().Msg("Deployments request received")
				ctx.Response.Header.Set("Content-Type", "application/json")
				deployments := informer.GetDeploymentNames()
				log.Info().Msgf("Deployments: %v", deployments)
				ctx.SetStatusCode(200)
				ctx.Write([]byte("["))
				for i, name := range deployments {
					ctx.WriteString("\"")
					ctx.WriteString(name)
					ctx.WriteString("\"")
					if i < len(deployments)-1 {
						ctx.WriteString(",")
					}
				}
				ctx.Write([]byte("]"))
				return
			case "/secrets":
				log.Info().Msg("Secrets request received")
				ctx.Response.Header.Set("Content-Type", "application/json")
				secrets := informer.GetSecretNames()
				log.Info().Msgf("Secrets: %v", secrets)
				ctx.SetStatusCode(200)
				ctx.Write([]byte("["))
				for i, name := range secrets {
					ctx.WriteString("\"")
					ctx.WriteString(name)
					ctx.WriteString("\"")
					if i < len(secrets)-1 {
						ctx.WriteString(",")
					}
				}
				ctx.Write([]byte("]"))
				return
			default:
				log.Info().
					Str("request_id", uuid).
					Str("method", string(ctx.Method())).
					Str("path", string(ctx.Path())).
					Str("remote_addr", ctx.RemoteAddr().String()).
					Msg("Incoming request")
				ctx.Response.Header.Set("X-Request-ID", uuid)
				_, err := fmt.Fprintf(ctx, `{"message:" "FastHTTP welcomes you, traveller!", "requestID": "%s"}`, uuid)
				if err != nil {
					return
				}
			}
		}
		addr := fmt.Sprintf(":%d", serverPort)
		log.Info().Msgf("Starting FastHTTP server on %s port", addr)
		if err := fasthttp.ListenAndServe(addr, handler); err != nil {
			log.Error().Err(err).Msg("Error starting FastHTTP server")
			os.Exit(1)
		}
	},
}

func getServerKubeClient(kubeconfigPath string, inCluster bool) (*kubernetes.Clientset, error) {
	var config *rest.Config
	var err error
	if inCluster {
		config, err = rest.InClusterConfig()
	} else {
		config, err = clientcmd.BuildConfigFromFlags("", kubeconfigPath)
	}
	if err != nil {
		return nil, err
	}
	return kubernetes.NewForConfig(config)
}

func init() {
	rootCmd.AddCommand(serverCmd)
	serverCmd.Flags().IntVar(&serverPort, "port", 8080, "Port to run the server on")
	serverCmd.Flags().StringVar(&serverKubeconfig, "kubeconfig", "", "Path to the kubeconfig file")
	serverCmd.Flags().BoolVar(&serverInCluster, "in-cluster", false, "Use in-cluster Kubernetes config")
	serverCmd.Flags().StringVar(&namespace, "watch-ns", "default", "Define the namespace to be watched by the informer, otherwise the default namespace is used")
}

package cmd

import (
	"context"
	"fmt"
	"github.com/buaazp/fasthttprouter"
	"github.com/go-logr/zerologr"
	"github.com/rs/zerolog/log"
	"github.com/silhouetteUA/k8s-controller/pkg/api"
	frontendv1alpha2 "github.com/silhouetteUA/k8s-controller/pkg/api/frontend/frontendBackup"
	frontendv1alpha1 "github.com/silhouetteUA/k8s-controller/pkg/api/frontend/v1alpha1"
	"github.com/silhouetteUA/k8s-controller/pkg/controller"
	"github.com/silhouetteUA/k8s-controller/pkg/informer"
	"github.com/spf13/cobra"
	"github.com/valyala/fasthttp"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/client-go/kubernetes"
	clientgoscheme "k8s.io/client-go/kubernetes/scheme"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"
	"os"
	ctrlruntime "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/cache"
	logf "sigs.k8s.io/controller-runtime/pkg/log"
	"sigs.k8s.io/controller-runtime/pkg/log/zap"
	"sigs.k8s.io/controller-runtime/pkg/manager"
	"sigs.k8s.io/controller-runtime/pkg/metrics/server"
)

var serverPort int
var serverKubeconfig string
var serverInCluster bool
var enableLeaderElection bool
var metricsPort int

var serverCmd = &cobra.Command{
	Use:   "server",
	Short: "Start a FastHTTP server",
	Run: func(cmd *cobra.Command, args []string) {
		level := parseLogLevel(logLevel)
		configureLogger(level)
		clientset, err := getServerKubeClient(serverKubeconfig, serverInCluster)
		if err != nil {
			log.Error().Err(err).Msg("Failed to create Kubernetes client")
			os.Exit(1)
		}
		ctx := context.Background()
		go informer.StartInformerFactory(ctx, clientset, namespace)
		logf.SetLogger(zap.New(zap.UseDevMode(true)))
		logf.SetLogger(zerologr.New(&log.Logger))
		// Start controller-runtime manager and controller
		scheme := runtime.NewScheme()
		if err := clientgoscheme.AddToScheme(scheme); err != nil {
			log.Error().Err(err).Msg("Failed to add client-go scheme")
			os.Exit(1)
		}
		if err := frontendv1alpha1.AddToScheme(scheme); err != nil {
			log.Error().Err(err).Msg("Failed to add FrontendPage scheme")
			os.Exit(1)
		}
		if err := frontendv1alpha2.AddToScheme(scheme); err != nil {
			log.Error().Err(err).Msg("Failed to add FrontendPageBackup scheme")
			os.Exit(1)
		}
		mgr, err := ctrlruntime.NewManager(ctrlruntime.GetConfigOrDie(), manager.Options{
			Scheme:                  scheme, // ADD YOUR OWN SCHEME, NOT A DEFAULT ONE !!!!!!
			LeaderElection:          enableLeaderElection,
			LeaderElectionID:        "k8s-controller-leader-election",
			LeaderElectionNamespace: namespace,
			Metrics:                 server.Options{BindAddress: fmt.Sprintf(":%d", metricsPort)},
			Cache:                   cache.Options{DefaultNamespaces: map[string]cache.Config{namespace: {}}},
		})
		if err != nil {
			log.Error().Err(err).Msg("Failed to create controller-runtime manager")
			os.Exit(1)
		}
		if err := controller.AddDeploymentController(mgr); err != nil {
			log.Error().Err(err).Msg("Failed to add deployment controller")
			os.Exit(1)
		}
		if err := controller.AddSecretController(mgr); err != nil {
			log.Error().Err(err).Msg("Failed to add secret controller")
			os.Exit(1)
		}
		if err := controller.AddFrontendController(mgr); err != nil {
			log.Error().Err(err).Msg("Failed to add frontend controller")
			os.Exit(1)
		}
		if err := controller.AddFrontendPageBackupController(mgr); err != nil {
			log.Error().Err(err).Msg("Failed to add frontend controller")
			os.Exit(1)
		}
		go func() {
			log.Info().Msg("Starting controller-runtime manager ... --watch-namespace=" + namespace)
			if err := mgr.Start(cmd.Context()); err != nil {
				log.Error().Err(err).Msg("Manager exited with error")
				os.Exit(1)
			}
		}()
		router := fasthttprouter.New()
		frontendAPI := &api.FrontendPageAPI{
			K8sClient: mgr.GetClient(),
			Namespace: namespace,
		}
		router.GET("/api/frontendpages", frontendAPI.ListFrontendPages)
		//curl -X POST -H "Content-Type: application/json" --data-binary "@config/crd/frontendPage_post.json" http://localhost:8080/api/frontendpages
		router.POST("/api/frontendpages", frontendAPI.CreateFrontendPage)
		router.GET("/api/frontendpages/:name", frontendAPI.GetFrontendPage)
		router.PUT("/api/frontendpages/:name", frontendAPI.UpdateFrontendPage)
		router.DELETE("/api/frontendpages/:name", frontendAPI.DeleteFrontendPage)
		//OLD way, can just parse the methods
		//handler := func(ctx *fasthttp.RequestCtx) {
		//	uuid := uuid.New().String()
		//	switch string(ctx.Path()) {
		//	case "/version":
		//		log.Info().
		//			Str("request_id", uuid).
		//			Str("method", string(ctx.Method())).
		//			Str("path", string(ctx.Path())).
		//			Str("remote_addr", ctx.RemoteAddr().String()).
		//			Msg("Check version request")
		//		ctx.Response.Header.SetContentType("application/json")
		//		ctx.Response.Header.Set("X-Request-ID", uuid)
		//		_, err := fmt.Fprintf(ctx, `{"version": "%s", "commit": "%s", "date": "%s", "requestID": "%s"}`, Version, Commit, BuildDate, uuid)
		//		if err != nil {
		//			return
		//		}
		//	case "/deployments":
		//		log.Info().Msg("Deployments request received")
		//		ctx.Response.Header.Set("Content-Type", "application/json")
		//		deployments := informer.GetDeploymentNames()
		//		log.Info().Msgf("Deployments: %v", deployments)
		//		ctx.SetStatusCode(200)
		//		ctx.Write([]byte("[")) //nolint:errcheck
		//		for i, name := range deployments {
		//			ctx.WriteString("\"") //nolint:errcheck
		//			ctx.WriteString(name) //nolint:errcheck
		//			ctx.WriteString("\"") //nolint:errcheck
		//			if i < len(deployments)-1 {
		//				ctx.WriteString(",") //nolint:errcheck
		//			}
		//		}
		//		ctx.Write([]byte("]")) //nolint:errcheck
		//		return
		//	case "/secrets":
		//		log.Info().Msg("Secrets request received")
		//		ctx.Response.Header.Set("Content-Type", "application/json")
		//		secrets := informer.GetSecretNames()
		//		log.Info().Msgf("Secrets: %v", secrets)
		//		ctx.SetStatusCode(200)
		//		ctx.Write([]byte("[")) //nolint:errcheck
		//		for i, name := range secrets {
		//			ctx.WriteString("\"") //nolint:errcheck
		//			ctx.WriteString(name) //nolint:errcheck
		//			ctx.WriteString("\"") //nolint:errcheck
		//			if i < len(secrets)-1 {
		//				ctx.WriteString(",") //nolint:errcheck
		//			}
		//		}
		//		ctx.Write([]byte("]")) //nolint:errcheck
		//		return
		//	default:
		//		log.Info().
		//			Str("request_id", uuid).
		//			Str("method", string(ctx.Method())).
		//			Str("path", string(ctx.Path())).
		//			Str("remote_addr", ctx.RemoteAddr().String()).
		//			Msg("Incoming request")
		//		ctx.Response.Header.Set("X-Request-ID", uuid)
		//		_, err := fmt.Fprintf(ctx, `{"message:" "FastHTTP welcomes you, traveller!", "requestID": "%s"}`, uuid)
		//		if err != nil {
		//			return
		//		}
		//	}
		//}
		addr := fmt.Sprintf(":%d", serverPort)
		log.Info().Msgf("Starting FastHTTP server on %s port", addr)
		if err := fasthttp.ListenAndServe(addr, router.Handler); err != nil { // here  you can switch between the handlers old=handler and new=router.Handler
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
	serverCmd.Flags().BoolVar(&enableLeaderElection, "enable-leader-election", true, "Enable leader election for controller manager")
	serverCmd.Flags().IntVar(&metricsPort, "metrics-port", 8081, "Port for controller manager metrics")
}

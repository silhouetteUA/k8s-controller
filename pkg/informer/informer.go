package informer

import (
	"context"
	"os"
	"time"

	"github.com/rs/zerolog/log"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/fields"
	"k8s.io/client-go/informers"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/cache"
)

func StartDeploymentInformer(ctx context.Context, clientset *kubernetes.Clientset, namespace string) {
	factory := informers.NewSharedInformerFactoryWithOptions(
		clientset,
		30*time.Second,
		informers.WithNamespace(namespace),
		informers.WithTweakListOptions(func(options *metav1.ListOptions) {
			options.FieldSelector = fields.Everything().String()
		}),
	)

	deploymentsInformer := factory.Apps().V1().Deployments().Informer()
	secretsInformer := factory.Core().V1().Secrets().Informer()

	addResourceHandlers(deploymentsInformer, "Deployment")
	addResourceHandlers(secretsInformer, "Secret")

	log.Info().Msg("Starting informers...")
	factory.Start(ctx.Done())

	for resource, synced := range factory.WaitForCacheSync(ctx.Done()) {
		if !synced {
			log.Error().Msgf("Failed to sync informer for %v", resource)
			os.Exit(1)
		}
	}

	log.Info().Msg("Informers cache synced. Watching for events...")
	<-ctx.Done()
}

func addResourceHandlers(informer cache.SharedIndexInformer, resourceType string) {
	_, err := informer.AddEventHandler(cache.ResourceEventHandlerFuncs{
		AddFunc: func(obj interface{}) {
			log.Info().Msgf("%s added: %s", resourceType, getObjectName(obj))
		},
		UpdateFunc: func(oldObj, newObj interface{}) {
			log.Info().Msgf("%s updated: %s", resourceType, getObjectName(newObj))
		},
		DeleteFunc: func(obj interface{}) {
			log.Info().Msgf("%s deleted: %s", resourceType, getObjectName(obj))
		},
	})
	if err != nil {
		return
	}
}

func getObjectName(obj interface{}) string {
	if o, ok := obj.(metav1.Object); ok {
		return o.GetName()
	}
	return "unknown"
}

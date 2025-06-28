package controller

import (
	context "context"
	v1 "k8s.io/api/core/v1"

	"github.com/rs/zerolog/log"
	"k8s.io/apimachinery/pkg/runtime"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/controller"
	"sigs.k8s.io/controller-runtime/pkg/manager"
)

type SecretReconciler struct {
	client.Client
	Scheme *runtime.Scheme
}

func (r *SecretReconciler) Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error) {
	log.Info().Msgf("Reconciling Secret: %s/%s", req.Namespace, req.Name)
	return ctrl.Result{}, nil
}

func AddSecretController(mgr manager.Manager) error {
	r := &SecretReconciler{
		Client: mgr.GetClient(),
		Scheme: mgr.GetScheme(),
	}
	return ctrl.NewControllerManagedBy(mgr).
		For(&v1.Secret{}).
		WithOptions(controller.Options{MaxConcurrentReconciles: 1}).
		Complete(r)
}

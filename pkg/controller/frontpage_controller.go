package controller

import (
	context "context"
	"reflect"

	"github.com/rs/zerolog/log"

	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"

	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/manager"

	frontendv1alpha1 "github.com/silhouetteUA/k8s-controller/pkg/api/frontend/v1alpha1"
)

type FrontendPageReconciler struct {
	client.Client
	Scheme *runtime.Scheme
}

func buildConfigMap(page *frontendv1alpha1.FrontendPage) *corev1.ConfigMap {
	return &corev1.ConfigMap{
		ObjectMeta: metav1.ObjectMeta{
			Name:      page.Name,
			Namespace: page.Namespace,
		},
		Data: map[string]string{
			"contents": page.Spec.Contents,
		},
	}
}

func buildDeployment(page *frontendv1alpha1.FrontendPage) *appsv1.Deployment {
	replicas := int32(page.Spec.Replicas)
	return &appsv1.Deployment{
		ObjectMeta: metav1.ObjectMeta{
			Name:      page.Name,
			Namespace: page.Namespace,
		},
		Spec: appsv1.DeploymentSpec{
			Replicas: &replicas,
			Selector: &metav1.LabelSelector{
				MatchLabels: map[string]string{"app": page.Name},
			},
			Template: corev1.PodTemplateSpec{
				ObjectMeta: metav1.ObjectMeta{
					Labels: map[string]string{"app": page.Name},
				},
				Spec: corev1.PodSpec{
					Containers: []corev1.Container{{
						Name:  "frontend",
						Image: page.Spec.Image,
						VolumeMounts: []corev1.VolumeMount{{
							Name:      "contents",
							MountPath: "/data",
						}},
					}},
					Volumes: []corev1.Volume{{
						Name: "contents",
						VolumeSource: corev1.VolumeSource{
							ConfigMap: &corev1.ConfigMapVolumeSource{
								LocalObjectReference: corev1.LocalObjectReference{
									Name: page.Name,
								},
							},
						},
					}},
				},
			},
		},
	}
}

func (r *FrontendPageReconciler) Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error) {
	var page frontendv1alpha1.FrontendPage
	err := r.Get(ctx, req.NamespacedName, &page)
	if err != nil {
		if client.IgnoreNotFound(err) == nil {
			// FrontendPage deleted: clean up resources
			log.Info().Msgf("FrontendPage deleted: %s %s", req.Name, req.Namespace)
			var cm corev1.ConfigMap
			cm.Name = req.Name
			cm.Namespace = req.Namespace
			_ = r.Delete(ctx, &cm) // ignore errors if not found
			var dep appsv1.Deployment
			dep.Name = req.Name
			dep.Namespace = req.Namespace
			_ = r.Delete(ctx, &dep)
			return ctrl.Result{}, nil
		}
		return ctrl.Result{}, err
	}

	// 1. Ensure ConfigMap exists and is up to date
	cm := buildConfigMap(&page)
	if err := ctrl.SetControllerReference(&page, cm, r.Scheme); err != nil {
		return ctrl.Result{}, err
	}

	log.Info().Msgf("Reconciling ConfigMap for FrontendPage: %s %s", cm.Name, cm.Namespace)
	var existingCM corev1.ConfigMap

	if err := r.Get(ctx, req.NamespacedName, &existingCM); err != nil {
		if !errors.IsNotFound(err) {
			return ctrl.Result{}, err
		}

		if err := r.Create(ctx, cm); err != nil {
			return ctrl.Result{}, err
		}
	} else if !reflect.DeepEqual(existingCM.Data, cm.Data) {
		existingCM.Data = cm.Data
		if err := r.Update(ctx, &existingCM); err != nil {
			return ctrl.Result{}, err
		}
	}

	// 2. Ensure Deployment exists and is up to date
	dep := buildDeployment(&page)
	if err := ctrl.SetControllerReference(&page, dep, r.Scheme); err != nil {
		return ctrl.Result{}, err
	}

	log.Info().Msgf("Reconciling Deployment for FrontendPage: %s %s", dep.Name, dep.Namespace)
	var existingDep appsv1.Deployment

	if err := r.Get(ctx, req.NamespacedName, &existingDep); err != nil {
		if !errors.IsNotFound(err) {
			return ctrl.Result{}, err
		}

		if err := r.Create(ctx, dep); err != nil {
			return ctrl.Result{}, err
		}
	} else {
		updated := false

		if *existingDep.Spec.Replicas != *dep.Spec.Replicas {
			existingDep.Spec.Replicas = dep.Spec.Replicas
			updated = true
		}

		if existingDep.Spec.Template.Spec.Containers[0].Image != dep.Spec.Template.Spec.Containers[0].Image {
			existingDep.Spec.Template.Spec.Containers[0].Image = dep.Spec.Template.Spec.Containers[0].Image
			updated = true
		}

		if updated {
			if err := r.Update(ctx, &existingDep); err != nil {
				if errors.IsConflict(err) {
					// Requeue to try again with the latest version
					return ctrl.Result{Requeue: true}, nil
				}
				return ctrl.Result{}, err
			}
		}
	}

	return ctrl.Result{}, nil
}

func AddFrontendController(mgr manager.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		For(&frontendv1alpha1.FrontendPage{}).
		Owns(&appsv1.Deployment{}).
		Owns(&corev1.ConfigMap{}).
		Complete(&FrontendPageReconciler{
			Client: mgr.GetClient(),
			Scheme: mgr.GetScheme(),
		})
}

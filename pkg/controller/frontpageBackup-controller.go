package controller

import (
	"context"
	"fmt"

	frontendv1alpha2 "github.com/silhouetteUA/k8s-controller/pkg/api/frontend/frontendBackup"
	frontendv1alpha1 "github.com/silhouetteUA/k8s-controller/pkg/api/frontend/v1alpha1"
	batchv1 "k8s.io/api/batch/v1"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

type FrontendPageBackupReconciler struct {
	client.Client
	Scheme *runtime.Scheme
}

func (r *FrontendPageBackupReconciler) Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error) {
	var backup frontendv1alpha2.FrontendPageBackup
	if err := r.Get(ctx, req.NamespacedName, &backup); err != nil {
		return ctrl.Result{}, client.IgnoreNotFound(err)
	}

	// Get the referenced FrontendPage
	var page frontendv1alpha1.FrontendPage
	if err := r.Get(ctx, client.ObjectKey{Name: backup.Spec.FrontendPageRef, Namespace: req.Namespace}, &page); err != nil {
		return ctrl.Result{}, err
	}

	cron := buildCronJob(&backup, &page)

	// Set owner
	if err := ctrl.SetControllerReference(&backup, cron, r.Scheme); err != nil {
		return ctrl.Result{}, err
	}

	var existing batchv1.CronJob
	err := r.Get(ctx, client.ObjectKey{Name: cron.Name, Namespace: cron.Namespace}, &existing)
	if errors.IsNotFound(err) {
		if err := r.Create(ctx, cron); err != nil {
			return ctrl.Result{}, err
		}
	} else if err != nil {
		return ctrl.Result{}, err
	}

	return ctrl.Result{}, nil
}

func buildCronJob(backup *frontendv1alpha2.FrontendPageBackup, page *frontendv1alpha1.FrontendPage) *batchv1.CronJob {
	jobName := fmt.Sprintf("backup-%s", page.Name)

	return &batchv1.CronJob{
		ObjectMeta: metav1.ObjectMeta{
			Name:      jobName,
			Namespace: backup.Namespace,
		},
		Spec: batchv1.CronJobSpec{
			Schedule: backup.Spec.Schedule,
			JobTemplate: batchv1.JobTemplateSpec{
				Spec: batchv1.JobSpec{
					Template: corev1.PodTemplateSpec{
						Spec: corev1.PodSpec{
							RestartPolicy: corev1.RestartPolicyNever,
							Containers: []corev1.Container{{
								Name:    "backup",
								Image:   "alpine",
								Command: []string{"sh", "-c"},
								Args: []string{
									fmt.Sprintf(`echo "%s" > /backup/%s.txt`, page.Spec.Contents, page.Name),
								},
								VolumeMounts: []corev1.VolumeMount{{
									Name:      "backup-vol",
									MountPath: "/backup",
								}},
							}},
							Volumes: []corev1.Volume{{
								Name: "backup-vol",
								VolumeSource: corev1.VolumeSource{
									EmptyDir: &corev1.EmptyDirVolumeSource{},
								},
							}},
						},
					},
				},
			},
		},
	}
}

func AddFrontendPageBackupController(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		For(&frontendv1alpha2.FrontendPageBackup{}).
		Owns(&batchv1.CronJob{}).
		Complete(&FrontendPageBackupReconciler{
			Client: mgr.GetClient(),
			Scheme: mgr.GetScheme(),
		})
}

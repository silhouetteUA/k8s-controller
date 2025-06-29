package controller

import (
	"context"
	"testing"
	"time"

	testutil "github.com/silhouetteUA/k8s-controller/pkg/testutil"
	"github.com/stretchr/testify/require"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

func TestSecretReconciler_BasicFlow(t *testing.T) {
	mgr, k8sClient, _, cleanup := testutil.StartTestManager(t)
	defer cleanup()

	// Register the Secret controller before starting the manager
	err := AddSecretController(mgr)
	require.NoError(t, err)

	go func() {
		_ = mgr.Start(context.Background())
	}()

	ns := "default"
	ctx := context.Background()
	name := "test-secret"

	secret := &corev1.Secret{
		ObjectMeta: metav1.ObjectMeta{
			Name:      name,
			Namespace: ns,
		},
		StringData: map[string]string{"username": "admin", "password": "secret"},
		Type:       corev1.SecretTypeOpaque,
	}
	if err := k8sClient.Create(ctx, secret); err != nil {
		t.Fatalf("Failed to create Secret: %v", err)
	}

	// Wait a bit to allow reconcile to be triggered
	time.Sleep(1 * time.Second)

	// Check that the secret still exists (reconcile didn't error or delete it)
	var got corev1.Secret
	err = k8sClient.Get(ctx, client.ObjectKey{Name: name, Namespace: ns}, &got)
	require.NoError(t, err)
	require.Equal(t, "admin", string(got.Data["username"]))
	require.Equal(t, "secret", string(got.Data["password"]))
}

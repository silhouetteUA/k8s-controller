package v1alpha2

import (
	"github.com/silhouetteUA/k8s-controller/pkg/api/frontend/v1alpha1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// FrontendPageBackupSpec defines the backup configuration
type FrontendPageBackupSpec struct {
	FrontendPageRef string `json:"frontendPageRef"` // name of the FrontendPage to back up
	Schedule        string `json:"schedule"`        // cron expression like "*/5 * * * *"
}

// FrontendPageBackupStatus shows backup progress
type FrontendPageBackupStatus struct {
	LastBackupTime *metav1.Time `json:"lastBackupTime,omitempty"`
	LastBackupPath string       `json:"lastBackupPath,omitempty"`
	Status         string       `json:"status,omitempty"` // success/failed/etc.
}

// +kubebuilder:object:root=true
// +kubebuilder:subresource:status
// +kubebuilder:resource:shortName=fpb,singular=frontendpagebackup,path=frontendpagebackups,scope=Namespaced
type FrontendPageBackup struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   FrontendPageBackupSpec   `json:"spec"`
	Status FrontendPageBackupStatus `json:"status,omitempty"`
}

// +kubebuilder:object:root=true

// FrontendPageBackupList contains a list of FrontendPageBackup
type FrontendPageBackupList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []FrontendPageBackup `json:"items"`
}

func init() {
	v1alpha1.SchemeBuilder.Register(&FrontendPageBackup{}, &FrontendPageBackupList{})
}

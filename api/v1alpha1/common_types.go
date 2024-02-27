package v1alpha1

import (
	corev1 "k8s.io/api/core/v1"
	"sigs.k8s.io/controller-runtime/pkg/client"

	velerov1 "github.com/vmware-tanzu/velero/pkg/apis/velero/v1"
)

type ClusterName string

// +kubebuilder:object:generate=false
type VeleroOrigin interface {
	client.Object
	*velerov1.Backup | *velerov1.Restore | *velerov1.Schedule
}

// +kubebuilder:object:generate=false
type VeleroProxy[T VeleroOrigin] interface {
	client.Object
	GetInstallRef() *corev1.ObjectReference
	SetClusterStatus(ClusterName, T)
}

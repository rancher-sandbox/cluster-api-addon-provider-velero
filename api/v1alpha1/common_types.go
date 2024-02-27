package v1alpha1

import (
	corev1 "k8s.io/api/core/v1"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

type ClusterName string

// +kubebuilder:object:generate=false
type InstalledObject interface {
	client.Object
	GetInstallRef() *corev1.ObjectReference
}

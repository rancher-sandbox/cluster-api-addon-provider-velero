package v1alpha1

import (
	"cmp"
	"strings"

	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/types"
	"sigs.k8s.io/controller-runtime/pkg/client"

	velerov1 "github.com/vmware-tanzu/velero/pkg/apis/velero/v1"
)

type NamespaceName string

func (c NamespaceName) ObjectKey() client.ObjectKey {
	clusterName := string(cmp.Or(c, "/"))
	namespaceName := strings.SplitN(clusterName, "/", 2)
	namespace, name := namespaceName[0], namespaceName[1]

	return types.NamespacedName{
		Namespace: namespace,
		Name:      name,
	}
}

func ToNamespaceName(obj client.Object) NamespaceName {
	return NamespaceName(client.ObjectKeyFromObject(obj).String())
}

func RefToNamespaceName(ref *corev1.ObjectReference) NamespaceName {
	if ref == nil {
		return NamespaceName("")
	}

	return NamespaceName(types.NamespacedName{
		Name:      ref.Name,
		Namespace: ref.Namespace,
	}.String())
}

// +kubebuilder:object:generate=false
type VeleroOrigin interface {
	client.Object
	*velerov1.Backup | *velerov1.Restore | *velerov1.Schedule
}

// +kubebuilder:object:generate=false
type VeleroProxy[V VeleroOrigin] interface {
	client.Object
	GetInstallRef() *corev1.ObjectReference
	SetClusterStatus(cluster NamespaceName, remote V)
}

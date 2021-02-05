package zitadel

import (
	"github.com/caos/orbos/pkg/kubernetes"
	"github.com/caos/orbos/pkg/tree"
	"github.com/caos/zitadel/operator/api/core"
	zitadelv1 "github.com/caos/zitadel/operator/api/zitadel/v1"
)

func ReadCrd(k8sClient kubernetes.ClientInt, namespace string, name string) (*tree.Tree, error) {
	unstruct, err := k8sClient.GetNamespacedCRDResource(zitadelv1.GroupVersion.Group, zitadelv1.GroupVersion.Version, "Zitadel", namespace, name)
	if err != nil {
		return nil, err
	}

	return core.UnmarshalUnstructuredSpec(unstruct)
}

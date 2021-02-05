package database

import (
	"github.com/caos/orbos/pkg/kubernetes"
	"github.com/caos/orbos/pkg/tree"
	"github.com/caos/zitadel/operator/api/core"
	databasev1 "github.com/caos/zitadel/operator/api/database/v1"
)

func ReadCrd(k8sClient kubernetes.ClientInt, namespace string, name string) (*tree.Tree, error) {
	unstruct, err := k8sClient.GetNamespacedCRDResource(databasev1.GroupVersion.Group, databasev1.GroupVersion.Version, "Database", namespace, name)
	if err != nil {
		return nil, err
	}

	return core.UnmarshalUnstructuredSpec(unstruct)
}

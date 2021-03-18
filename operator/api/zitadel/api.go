package zitadel

import (
	"github.com/caos/orbos/pkg/kubernetes"
	"github.com/caos/orbos/pkg/tree"
	"github.com/caos/zitadel/operator/api/core"
	zitadelv1 "github.com/caos/zitadel/operator/api/zitadel/v1"
	macherrs "k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/api/meta"
)

const (
	namespace  = "caos-system"
	kind       = "Zitadel"
	apiVersion = "caos.ch/v1"
	name       = "zitadel"
)

func ReadCrd(k8sClient kubernetes.ClientInt) (*tree.Tree, error) {
	unstruct, err := k8sClient.GetNamespacedCRDResource(zitadelv1.GroupVersion.Group, zitadelv1.GroupVersion.Version, kind, namespace, name)
	if err != nil {
		if macherrs.IsNotFound(err) || meta.IsNoMatchError(err) {
			return nil, nil
		}
		return nil, err
	}

	return core.UnmarshalUnstructuredSpec(unstruct)
}

func WriteCrd(k8sClient kubernetes.ClientInt, t *tree.Tree) error {

	unstruct, err := core.MarshalToUnstructuredSpec(t)
	if err != nil {
		return err
	}

	unstruct.SetName(name)
	unstruct.SetNamespace(namespace)
	unstruct.SetKind(kind)
	unstruct.SetAPIVersion(apiVersion)

	return k8sClient.ApplyNamespacedCRDResource(zitadelv1.GroupVersion.Group, zitadelv1.GroupVersion.Version, kind, namespace, name, unstruct)
}

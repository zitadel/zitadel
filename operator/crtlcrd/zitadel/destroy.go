package zitadel

import (
	"github.com/caos/orbos/mntr"
	"github.com/caos/orbos/pkg/kubernetes"
	"github.com/caos/orbos/pkg/tree"
	"github.com/caos/zitadel/operator/api/core"
	"github.com/caos/zitadel/operator/api/zitadel"
	orbz "github.com/caos/zitadel/operator/zitadel/kinds/orb"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
)

func Destroy(monitor mntr.Monitor, k8sClient kubernetes.ClientInt, version string) error {
	desired, err := zitadel.ReadCrd(k8sClient)
	if err != nil {
		return err
	}

	if desired == nil {
		unstruct := &unstructured.Unstructured{
			Object: map[string]interface{}{
				"apiVersion": "caos.ch/v1",
				"kind":       "Zitadel",
				"spec": map[string]interface{}{
					"kind":    "zitadel.caos.ch/Orb",
					"version": "v0",
					"spec":    map[string]interface{}{},
					"iam": map[string]interface{}{
						"kind":    "zitadel.caos.ch/ZITADEL",
						"version": "v0",
					},
				},
			},
		}

		desiredT, err := core.UnmarshalUnstructuredSpec(unstruct)
		if err != nil {
			return err
		}
		desired = desiredT
	}

	_, destroy, _, _, _, _, err := orbz.AdaptFunc(nil, "ensure", &version, false, []string{"operator", "iam"})(monitor, desired, &tree.Tree{})
	if err != nil {
		return err
	}

	return destroy(k8sClient)
}

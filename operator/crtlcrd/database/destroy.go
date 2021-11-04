package database

import (
	"github.com/caos/orbos/mntr"
	"github.com/caos/orbos/pkg/kubernetes"
	"github.com/caos/orbos/pkg/tree"
	"github.com/caos/zitadel/operator/api/core"
	"github.com/caos/zitadel/operator/api/database"
	orbdb "github.com/caos/zitadel/operator/database/kinds/orb"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
)

func Destroy(monitor mntr.Monitor, k8sClient kubernetes.ClientInt, version string) error {
	desired, err := database.ReadCrd(k8sClient)
	if err != nil {
		return err
	}

	if desired == nil {
		unstruct := &unstructured.Unstructured{
			Object: map[string]interface{}{
				"apiVersion": "caos.ch/v1",
				"kind":       "Database",
				"spec": map[string]interface{}{
					"kind":    "databases.caos.ch/Orb",
					"version": "v0",
					"spec":    map[string]interface{}{},
					"database": map[string]interface{}{
						"kind":    "databases.caos.ch/CockroachDB",
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

	_, destroy, _, _, _, _, err := orbdb.AdaptFunc("", &version, false, "operator", "database", "backup")(monitor, desired, &tree.Tree{})
	if err != nil {
		return err
	}

	if err := destroy(k8sClient); err != nil {
		return err
	}
	return nil
}

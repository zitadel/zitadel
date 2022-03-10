package zitadel

import (
	"github.com/caos/orbos/mntr"
	"github.com/caos/orbos/pkg/kubernetes"
	"github.com/caos/orbos/pkg/tree"
	"github.com/caos/zitadel/operator/api/zitadel"
	orbz "github.com/caos/zitadel/operator/zitadel/kinds/orb"
)

func Destroy(monitor mntr.Monitor, k8sClient kubernetes.ClientInt, version string) error {
	desired, err := zitadel.ReadCrd(k8sClient)
	if err != nil {
		return err
	}

	if desired != nil {
		_, destroy, _, _, _, _, err := orbz.AdaptFunc(nil, "ensure", &version, false, []string{"operator", "iam"})(monitor, desired, &tree.Tree{})
		if err != nil {
			return err
		}

		return destroy(k8sClient)
	}
	return nil
}

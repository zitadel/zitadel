package database

import (
	"github.com/caos/orbos/mntr"
	"github.com/caos/orbos/pkg/kubernetes"
	"github.com/caos/orbos/pkg/tree"
	"github.com/caos/zitadel/operator/api/database"
	orbdb "github.com/caos/zitadel/operator/database/kinds/orb"
)

func Destroy(monitor mntr.Monitor, k8sClient kubernetes.ClientInt, version string) error {
	desired, err := database.ReadCrd(k8sClient)
	if err != nil {
		return err
	}

	if desired != nil {
		_, destroy, _, _, _, _, err := orbdb.AdaptFunc("", &version, false, "database")(monitor, desired, &tree.Tree{})
		if err != nil {
			return err
		}

		return destroy(k8sClient)
	}
	return nil
}

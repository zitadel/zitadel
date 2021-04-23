package crtlcrd

import (
	"github.com/caos/orbos/mntr"
	"github.com/caos/orbos/pkg/kubernetes"
	"github.com/caos/orbos/pkg/tree"
	"github.com/caos/zitadel/operator/api/database"
	orbdb "github.com/caos/zitadel/operator/database/kinds/orb"
)

func Restore(monitor mntr.Monitor, k8sClient *kubernetes.Client, backup string, binaryVersion *string) error {
	desired, err := database.ReadCrd(k8sClient)
	if err != nil {
		return err
	}

	query, _, _, _, _, _, err := orbdb.AdaptFunc(backup, binaryVersion, false, "restore")(monitor, desired, &tree.Tree{})
	if err != nil {
		return err
	}

	ensure, err := query(k8sClient, map[string]interface{}{})
	if err != nil {
		return err
	}

	if err := ensure(k8sClient); err != nil {
		monitor.Error(err)
		return err
	}

	return nil
}

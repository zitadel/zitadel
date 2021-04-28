package crtlcrd

import (
	"github.com/caos/orbos/mntr"
	"github.com/caos/orbos/pkg/kubernetes"
	"github.com/caos/zitadel/pkg/databases"
)

func Restore(
	monitor mntr.Monitor,
	k8sClient *kubernetes.Client,
	backup string,
) error {
	if err := databases.CrdClear(monitor, k8sClient); err != nil {
		return err
	}

	if err := databases.CrdRestore(
		monitor,
		k8sClient,
		backup,
	); err != nil {
		return err
	}

	return nil
}

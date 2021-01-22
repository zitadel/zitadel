package backup

import (
	"github.com/caos/orbos/mntr"
	"github.com/caos/orbos/pkg/kubernetes"
	"github.com/caos/zitadel/operator"
	"github.com/pkg/errors"
)

func getCleanupFunc(monitor mntr.Monitor, namespace string, name string) operator.EnsureFunc {
	return func(k8sClient kubernetes.ClientInt) error {
		monitor.Info("waiting for backup to be completed")
		if err := k8sClient.WaitUntilJobCompleted(namespace, name, timeout); err != nil {
			monitor.Error(errors.Wrap(err, "error while waiting for backup to be completed"))
			return err
		}
		monitor.Info("backup is completed, cleanup")
		if err := k8sClient.DeleteJob(namespace, name); err != nil {
			monitor.Error(errors.Wrap(err, "error while trying to cleanup backup"))
			return err
		}
		monitor.Info("restore backup is completed")
		return nil
	}
}

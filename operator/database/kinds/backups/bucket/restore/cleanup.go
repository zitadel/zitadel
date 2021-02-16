package restore

import (
	"github.com/caos/orbos/mntr"
	"github.com/caos/orbos/pkg/kubernetes"
	"github.com/caos/zitadel/operator"
	"github.com/pkg/errors"
)

func getCleanupFunc(monitor mntr.Monitor, namespace, jobName string) operator.EnsureFunc {
	return func(k8sClient kubernetes.ClientInt) error {
		monitor.Info("waiting for restore to be completed")
		if err := k8sClient.WaitUntilJobCompleted(namespace, jobName, timeout); err != nil {
			monitor.Error(errors.Wrap(err, "error while waiting for restore to be completed"))
			return err
		}
		monitor.Info("restore is completed, cleanup")
		if err := k8sClient.DeleteJob(namespace, jobName); err != nil {
			monitor.Error(errors.Wrap(err, "error while trying to cleanup restore"))
			return err
		}
		monitor.Info("restore cleanup is completed")
		return nil
	}
}

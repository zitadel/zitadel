package clean

import (
	"github.com/caos/orbos/mntr"
	"github.com/caos/orbos/pkg/kubernetes"
	"github.com/caos/zitadel/operator"
	"github.com/pkg/errors"
)

func GetCleanupFunc(
	monitor mntr.Monitor,
	namespace string,
	backupName string,
) operator.EnsureFunc {
	return func(k8sClient kubernetes.ClientInt) error {
		monitor.Info("waiting for clean to be completed")
		if err := k8sClient.WaitUntilJobCompleted(namespace, GetJobName(backupName), timeout); err != nil {
			monitor.Error(errors.Wrap(err, "error while waiting for clean to be completed"))
			return err
		}
		monitor.Info("clean is completed, cleanup")
		if err := k8sClient.DeleteJob(namespace, GetJobName(backupName)); err != nil {
			monitor.Error(errors.Wrap(err, "error while trying to cleanup clean"))
			return err
		}
		monitor.Info("clean cleanup is completed")
		return nil
	}
}

package backup

import (
	"fmt"

	"github.com/caos/orbos/mntr"
	"github.com/caos/orbos/pkg/kubernetes"

	"github.com/caos/zitadel/operator"
)

func GetCleanupFunc(
	monitor mntr.Monitor,
	namespace string,
	backupName string,
) operator.EnsureFunc {
	return func(k8sClient kubernetes.ClientInt) error {
		monitor.Info("waiting for backup to be completed")
		if err := k8sClient.WaitUntilJobCompleted(namespace, GetJobName(backupName), timeout); err != nil {
			return fmt.Errorf("error while waiting for backup to be completed: %w", err)
		}
		monitor.Info("backup is completed, cleanup")
		if err := k8sClient.DeleteJob(namespace, GetJobName(backupName)); err != nil {
			return fmt.Errorf("error while trying to cleanup backup: %w", err)
		}
		monitor.Info("restore backup is completed")
		return nil
	}
}

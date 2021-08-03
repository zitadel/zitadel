package clean

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
		monitor.Info("waiting for clean to be completed")
		if err := k8sClient.WaitUntilJobCompleted(namespace, GetJobName(backupName), timeout); err != nil {
			return fmt.Errorf("error while waiting for clean to be completed: %w", err)
		}
		monitor.Info("clean is completed, cleanup")
		if err := k8sClient.DeleteJob(namespace, GetJobName(backupName)); err != nil {
			return fmt.Errorf("error while trying to cleanup clean: %w", err)
		}
		monitor.Info("clean cleanup is completed")
		return nil
	}
}

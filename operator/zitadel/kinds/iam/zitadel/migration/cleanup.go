package migration

import (
	"fmt"

	"github.com/caos/orbos/mntr"
	"github.com/caos/orbos/pkg/kubernetes"
	"github.com/caos/zitadel/operator"
)

func GetCleanupFunc(
	monitor mntr.Monitor,
	namespace string,
	reason string,
) operator.EnsureFunc {
	return func(k8sClient kubernetes.ClientInt) error {
		monitor.Info("cleanup migration job")
		if err := k8sClient.DeleteJob(namespace, getJobName(reason)); err != nil {
			return fmt.Errorf("error during job deletion: %w", err)
		}
		monitor.Info("migration cleanup is completed")
		return nil
	}
}

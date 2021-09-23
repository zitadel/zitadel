package migration

import (
	"fmt"
	"time"

	"github.com/caos/orbos/mntr"
	"github.com/caos/orbos/pkg/kubernetes"

	"github.com/caos/zitadel/operator"
)

const (
	timeout = 20 * time.Minute
)

func GetDoneFunc(
	monitor mntr.Monitor,
	namespace string,
	reason string,
) operator.EnsureFunc {
	return func(k8sClient kubernetes.ClientInt) error {
		monitor.Info("waiting for migration to be completed")
		if err := k8sClient.WaitUntilJobCompleted(namespace, getJobName(reason), timeout); err != nil {
			return fmt.Errorf("error while waiting for migration to be completed: %w", err)
		}
		monitor.Info("migration is completed")
		return nil
	}
}

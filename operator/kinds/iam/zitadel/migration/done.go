package migration

import (
	"github.com/caos/orbos/mntr"
	"github.com/caos/orbos/pkg/kubernetes"
	"github.com/caos/zitadel/operator"
	"github.com/pkg/errors"
	"time"
)

const (
	timeout time.Duration = 300
)

func GetDoneFunc(
	monitor mntr.Monitor,
	namespace string,
	reason string,
) operator.EnsureFunc {
	return func(k8sClient kubernetes.ClientInt) error {
		monitor.Info("waiting for migration to be completed")
		if err := k8sClient.WaitUntilJobCompleted(namespace, getJobName(reason), timeout); err != nil {
			monitor.Error(errors.Wrap(err, "error while waiting for migration to be completed"))
			return err
		}
		monitor.Info("migration is completed")
		return nil
	}
}

package setup

import (
	"github.com/caos/orbos/mntr"
	"github.com/caos/orbos/pkg/kubernetes"
	"github.com/caos/zitadel/operator"
	"github.com/pkg/errors"
)

func GetDoneFunc(
	monitor mntr.Monitor,
	namespace string,
	reason string,
) operator.EnsureFunc {
	return func(k8sClient kubernetes.ClientInt) error {
		monitor.Info("waiting for setup to be completed")
		if err := k8sClient.WaitUntilJobCompleted(namespace, getJobName(reason), timeout); err != nil {
			monitor.Error(errors.Wrap(err, "error while waiting for setup to be completed"))
			return err
		}
		monitor.Info("migration is completed")
		return nil
	}
}

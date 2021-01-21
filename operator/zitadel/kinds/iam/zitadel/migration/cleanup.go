package migration

import (
	"github.com/caos/orbos/mntr"
	"github.com/caos/orbos/pkg/kubernetes"
	"github.com/caos/zitadel/operator"
	"github.com/pkg/errors"
)

func GetCleanupFunc(
	monitor mntr.Monitor,
	namespace string,
	reason string,
) operator.EnsureFunc {
	return func(k8sClient kubernetes.ClientInt) error {
		monitor.Info("cleanup migration job")
		if err := k8sClient.DeleteJob(namespace, getJobName(reason)); err != nil {
			monitor.Error(errors.Wrap(err, "error during job deletion"))
			return err
		}
		monitor.Info("migration cleanup is completed")
		return nil
	}
}

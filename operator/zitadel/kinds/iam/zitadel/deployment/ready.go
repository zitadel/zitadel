package deployment

import (
	"github.com/caos/orbos/mntr"
	"github.com/caos/orbos/pkg/kubernetes"
	"github.com/caos/orbos/pkg/labels"
	"github.com/caos/zitadel/operator"
	"github.com/pkg/errors"
)

func GetReadyFunc(monitor mntr.Monitor, namespace string, name *labels.Name) operator.EnsureFunc {
	return func(k8sClient kubernetes.ClientInt) error {
		monitor.Info("waiting for deployment to be ready")
		if err := k8sClient.WaitUntilDeploymentReady(namespace, name.Name(), true, true, 60); err != nil {
			monitor.Error(errors.Wrap(err, "error while waiting for deployment to be ready"))
			return err
		}
		monitor.Info("deployment is ready")
		return nil
	}
}

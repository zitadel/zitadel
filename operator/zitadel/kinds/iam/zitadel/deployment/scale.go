package deployment

import (
	"github.com/caos/orbos/mntr"
	"github.com/caos/orbos/pkg/kubernetes"
	"github.com/caos/orbos/pkg/labels"
	"github.com/caos/zitadel/operator"
)

func GetScaleFunc(monitor mntr.Monitor, namespace string, name *labels.Name) func(replicaCount int) operator.EnsureFunc {
	return func(replicaCount int) operator.EnsureFunc {
		return func(k8sClient kubernetes.ClientInt) error {
			monitor.Info("Scaling deployment")
			return k8sClient.ScaleDeployment(namespace, name.Name(), replicaCount)
		}
	}
}

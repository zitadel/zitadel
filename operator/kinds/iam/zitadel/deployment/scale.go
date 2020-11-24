package deployment

import (
	"github.com/caos/orbos/mntr"
	"github.com/caos/orbos/pkg/kubernetes"
	"github.com/caos/zitadel/operator"
)

func GetScaleFunc(monitor mntr.Monitor, namespace string) func(replicaCount int) operator.EnsureFunc {
	return func(replicaCount int) operator.EnsureFunc {
		return func(k8sClient kubernetes.ClientInt) error {
			monitor.Info("Scaling deployment")
			return k8sClient.ScaleDeployment(namespace, deployName, replicaCount)
		}
	}
}

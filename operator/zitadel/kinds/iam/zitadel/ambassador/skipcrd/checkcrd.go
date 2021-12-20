package skipcrd

import (
	"github.com/caos/orbos/mntr"
	"github.com/caos/orbos/pkg/kubernetes"
	"github.com/caos/zitadel/operator"
)

func EnsureFunc(monitor mntr.Monitor, k8sClient kubernetes.ClientInt, def string) (operator.EnsureFunc, error) {
	_, ok, err := k8sClient.CheckCRD(def)
	if err != nil {
		return nil, err
	}
	if !ok {
		return func(k8sClient kubernetes.ClientInt) error {
			monitor.WithField("crd", def).Info("Skipped applying ambassador CRD, as definition doesn't exist in cluster")
			return nil
		}, nil
	}
	return nil, nil
}

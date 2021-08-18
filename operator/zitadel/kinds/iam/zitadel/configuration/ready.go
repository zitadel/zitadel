package configuration

import (
	"fmt"

	"github.com/caos/orbos/mntr"
	"github.com/caos/orbos/pkg/kubernetes"
)

func GetReadyFunc(
	monitor mntr.Monitor,
	namespace string,
	secretName string,
	secretVarsName string,
	secretPasswordName string,
	cmName string,
	consoleCMName string,
) func(k8sClient kubernetes.ClientInt) error {
	return func(k8sClient kubernetes.ClientInt) error {
		monitor.Debug("Waiting for configuration to be created")
		if err := k8sClient.WaitForSecret(namespace, secretName, timeout); err != nil {
			return fmt.Errorf("error while waiting for secret: %w", err)
		}

		if err := k8sClient.WaitForSecret(namespace, secretVarsName, timeout); err != nil {
			return fmt.Errorf("error while waiting for vars secret: %w", err)
		}

		if err := k8sClient.WaitForSecret(namespace, secretPasswordName, timeout); err != nil {
			return fmt.Errorf("error while waiting for password secret: %w", err)
		}

		if err := k8sClient.WaitForConfigMap(namespace, cmName, timeout); err != nil {
			return fmt.Errorf("error while waiting for configmap: %w", err)
		}

		if err := k8sClient.WaitForConfigMap(namespace, consoleCMName, timeout); err != nil {
			return fmt.Errorf("error while waiting for console configmap: %w", err)
		}
		monitor.Debug("configuration is created")
		return nil
	}
}

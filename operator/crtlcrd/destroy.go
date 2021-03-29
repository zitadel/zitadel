package crtlcrd

import (
	"github.com/caos/orbos/mntr"
	"github.com/caos/orbos/pkg/kubernetes"
	"github.com/caos/zitadel/operator/crtlcrd/database"
	"github.com/caos/zitadel/operator/crtlcrd/zitadel"
)

func Destroy(monitor mntr.Monitor, k8sClient kubernetes.ClientInt, features ...string) error {
	for _, feature := range features {
		switch feature {
		case Zitadel:
			if err := zitadel.Destroy(monitor, k8sClient); err != nil {
				return err
			}
		case Database:
			if err := database.Destroy(monitor, k8sClient); err != nil {
				return err
			}
		}
	}
	return nil
}

package zitadel

import (
	"github.com/caos/orbos/mntr"
	"github.com/caos/orbos/pkg/kubernetes"
	"github.com/caos/zitadel/operator/api/zitadel"
	"github.com/caos/zitadel/operator/zitadel/kinds/orb"
	macherrs "k8s.io/apimachinery/pkg/api/errors"
)

func ScaleDown(
	monitor mntr.Monitor,
	k8sClient *kubernetes.Client,
	version *string,
) (bool, error) {

	desired, err := zitadel.ReadCrd(k8sClient)
	if err != nil {
		return false, err
	}

	noZitadel := false
	if err := Takeoff(monitor, k8sClient, desired, orb.AdaptFunc(nil, "scaledown", version, false, []string{"scaledown"})); err != nil {
		if macherrs.IsNotFound(err) {
			noZitadel = true
		} else {
			return noZitadel, err
		}
	}
	return noZitadel, nil
}

func ScaleUp(
	monitor mntr.Monitor,
	k8sClient *kubernetes.Client,
	version *string,
) error {
	desired, err := zitadel.ReadCrd(k8sClient)
	if err != nil {
		return err
	}

	return Takeoff(monitor, k8sClient, desired, orb.AdaptFunc(nil, "scaleup", version, false, []string{"scaleup"}))
}

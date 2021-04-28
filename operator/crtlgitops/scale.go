package crtlgitops

import (
	"github.com/caos/orbos/mntr"
	"github.com/caos/orbos/pkg/git"
	"github.com/caos/orbos/pkg/kubernetes"
	orbconfig "github.com/caos/orbos/pkg/orb"
	"github.com/caos/zitadel/operator/zitadel"
	"github.com/caos/zitadel/operator/zitadel/kinds/orb"
	macherrs "k8s.io/apimachinery/pkg/api/errors"
)

func ScaleDown(
	monitor mntr.Monitor,
	gitClient *git.Client,
	k8sClient *kubernetes.Client,
	orbCfg *orbconfig.Orb,
	version *string,
	gitops bool,
) (bool, error) {
	noZitadel := false
	if err := zitadel.Takeoff(monitor, gitClient, orb.AdaptFunc(orbCfg, "scaledown", version, gitops, []string{"scaledown"}), k8sClient)(); err != nil {
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
	gitClient *git.Client,
	k8sClient *kubernetes.Client,
	orbCfg *orbconfig.Orb,
	version *string,
	gitops bool,
) error {
	return zitadel.Takeoff(monitor, gitClient, orb.AdaptFunc(orbCfg, "scaleup", version, gitops, []string{"scaleup"}), k8sClient)()
}

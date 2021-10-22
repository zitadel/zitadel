package cmds

import (
	"github.com/caos/orbos/mntr"
	"github.com/caos/orbos/pkg/git"
	"github.com/caos/orbos/pkg/kubernetes"
	orbconfig "github.com/caos/orbos/pkg/orb"
	"github.com/caos/zitadel/operator/crtlcrd/zitadel"
	"github.com/caos/zitadel/operator/crtlgitops"
	kubernetes2 "github.com/caos/zitadel/pkg/kubernetes"
	macherrs "k8s.io/apimachinery/pkg/api/errors"
)

func scaleForFunction(
	monitor mntr.Monitor,
	gitClient *git.Client,
	orbCfg *orbconfig.Orb,
	k8sClient *kubernetes.Client,
	version *string,
	gitops bool,
	ensureFunc func() error,
) error {
	noOperator := false
	if err := kubernetes2.ScaleZitadelOperator(monitor, k8sClient, 0); err != nil {
		if macherrs.IsNotFound(err) {
			noOperator = true
		} else {
			return err
		}
	}

	noZitadel := false
	if gitops {
		noZitadelT, err := crtlgitops.ScaleDown(monitor, gitClient, k8sClient, orbCfg, version, gitops)
		if err != nil {
			return err
		}
		noZitadel = noZitadelT
	} else {
		noZitadelT, err := zitadel.ScaleDown(monitor, k8sClient, version)
		if err != nil {
			return err
		}
		noZitadel = noZitadelT
	}

	noDatabase := false
	if err := kubernetes2.ScaleDatabaseOperator(monitor, k8sClient, 0); err != nil {
		if macherrs.IsNotFound(err) {
			noDatabase = true
		} else {
			return err
		}
	}

	if err := ensureFunc(); err != nil {
		return err
	}

	if !noDatabase {
		if err := kubernetes2.ScaleDatabaseOperator(monitor, k8sClient, 1); err != nil {
			return err
		}
	}

	if !noZitadel {
		if gitops {
			if err := crtlgitops.ScaleUp(monitor, gitClient, k8sClient, orbCfg, version, gitops); err != nil {
				return err
			}
		} else {
			if err := zitadel.ScaleUp(monitor, k8sClient, version); err != nil {
				return err
			}
		}
	}

	if !noOperator {
		if err := kubernetes2.ScaleZitadelOperator(monitor, k8sClient, 1); err != nil {
			return err
		}
	}

	return nil
}

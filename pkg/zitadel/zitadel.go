package zitadel

import (
	"github.com/caos/orbos/mntr"
	"github.com/caos/orbos/pkg/git"
	"github.com/caos/orbos/pkg/kubernetes"
	orbconfig "github.com/caos/orbos/pkg/orb"
	"github.com/caos/orbos/pkg/tree"
	"github.com/caos/zitadel/operator/api"
	"github.com/caos/zitadel/operator/api/zitadel"
	"github.com/caos/zitadel/operator/zitadel/kinds/orb"
)

func CrdScaleDown(
	monitor mntr.Monitor,
	k8sClient *kubernetes.Client,
	version *string,
) error {
	desired, err := zitadel.ReadCrd(k8sClient)
	if err != nil {
		monitor.Error(err)
		return err
	}

	return scaleDown(monitor, nil, k8sClient, false, version, desired)
}

func GitOpsScaleDown(
	monitor mntr.Monitor,
	orbCfg *orbconfig.Orb,
	gitClient *git.Client,
	k8sClient *kubernetes.Client,
	version *string,
) error {
	desired, err := api.ReadZitadelYml(gitClient)
	if err != nil {
		monitor.Error(err)
		return err
	}

	return scaleDown(monitor, orbCfg, k8sClient, true, version, desired)
}

//Take care! to use this function you have to include migration files into the binary
func scaleDown(
	monitor mntr.Monitor,
	orbCfg *orbconfig.Orb,
	k8sClient *kubernetes.Client,
	gitops bool,
	version *string,
	desired *tree.Tree,
) error {
	current := &tree.Tree{}
	query, _, _, _, _, err := orb.AdaptFunc(orbCfg, "scaledown", version, gitops, []string{"scaledown"})(monitor, desired, current)
	if err != nil {
		return err
	}

	ensure, err := query(k8sClient, map[string]interface{}{})
	if err != nil {
		return err
	}

	if err := ensure(k8sClient); err != nil {
		return err
	}
	return nil
}

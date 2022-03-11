package zitadel

import (
	"github.com/caos/orbos/mntr"
	"github.com/caos/orbos/pkg/git"
	"github.com/caos/orbos/pkg/kubernetes"
	"github.com/caos/orbos/pkg/tree"
	"github.com/caos/zitadel/operator/api/zitadel"
	"github.com/caos/zitadel/operator/zitadel/kinds/orb"
	"github.com/caos/zitadel/pkg/databases/db"
)

func CrdScaleDown(
	monitor mntr.Monitor,
	k8sClient *kubernetes.Client,
	version *string,
	dbClient db.Client,
) error {
	desired, err := zitadel.ReadCrd(k8sClient)
	if err != nil {
		return err
	}

	return scaleDown(monitor, k8sClient, false, version, desired, dbClient)
}

func GitOpsScaleDown(
	monitor mntr.Monitor,
	gitClient *git.Client,
	k8sClient *kubernetes.Client,
	version *string,
	dbClient db.Client,
) error {
	desired, err := gitClient.ReadTree(git.ZitadelFile)
	if err != nil {
		return err
	}

	return scaleDown(monitor, k8sClient, true, version, desired, dbClient)
}

//Take care! to use this function you have to include migration files into the binary
func scaleDown(
	monitor mntr.Monitor,
	k8sClient *kubernetes.Client,
	gitops bool,
	version *string,
	desired *tree.Tree,
	dbClient db.Client,
) error {
	current := &tree.Tree{}
	query, _, _, _, _, _, err := orb.AdaptFunc("scaledown", version, gitops, []string{"scaledown"}, dbClient)(monitor, desired, current)
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

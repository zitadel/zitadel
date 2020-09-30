package iam

import (
	"fmt"
	"github.com/caos/orbos/pkg/orb"

	"github.com/caos/orbos/mntr"
	"github.com/caos/orbos/pkg/secret"
	"github.com/caos/orbos/pkg/tree"
	"github.com/caos/zitadel/operator"
	"github.com/caos/zitadel/operator/kinds/iam/zitadel"
	"github.com/pkg/errors"
	core "k8s.io/api/core/v1"
)

func GetQueryAndDestroyFuncs(
	monitor mntr.Monitor,
	desiredTree *tree.Tree,
	currentTree *tree.Tree,
	nodeselector map[string]string,
	tolerations []core.Toleration,
	orbconfig *orb.Orb,
) (
	query operator.QueryFunc,
	destroy operator.DestroyFunc,
	err error,
) {

	defer func() {
		if err != nil {
			err = fmt.Errorf("adapting %s failed: %w", desiredTree.Common.Kind, err)
		}
	}()

	switch desiredTree.Common.Kind {
	case "zitadel.caos.ch/Zitadel":
		return zitadel.AdaptFunc(nodeselector, tolerations, orbconfig)(monitor, desiredTree, currentTree)
	default:
		return nil, nil, errors.Errorf("unknown iam kind %s", desiredTree.Common.Kind)
	}
}

func GetSecrets(
	monitor mntr.Monitor,
	desiredTree *tree.Tree,
) (
	map[string]*secret.Secret,
	error,
) {

	switch desiredTree.Common.Kind {
	case "zitadel.caos.ch/Zitadel":
		return zitadel.SecretsFunc()(monitor, desiredTree)
	default:
		return nil, errors.Errorf("unknown iam kind %s", desiredTree.Common.Kind)
	}
}

func GetBackupList(
	monitor mntr.Monitor,
	desiredTree *tree.Tree,
) (
	[]string,
	error,
) {
	switch desiredTree.Common.Kind {
	case "zitadel.caos.ch/Zitadel":
		return zitadel.BackupListFunc()(monitor, desiredTree)
	default:
		return nil, errors.Errorf("unknown iam kind %s", desiredTree.Common.Kind)
	}
}

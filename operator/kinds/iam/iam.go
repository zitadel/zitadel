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
	action string,
	migrationsPath string,
	version string,
	features []string,
) (
	query operator.QueryFunc,
	destroy operator.DestroyFunc,
	secrets map[string]*secret.Secret,
	err error,
) {

	defer func() {
		if err != nil {
			err = fmt.Errorf("adapting %s failed: %w", desiredTree.Common.Kind, err)
		}
	}()

	switch desiredTree.Common.Kind {
	case "zitadel.caos.ch/Zitadel":
		return zitadel.AdaptFunc(nodeselector, tolerations, orbconfig, action, migrationsPath, version, features)(monitor, desiredTree, currentTree)
	default:
		return nil, nil, nil, errors.Errorf("unknown iam kind %s", desiredTree.Common.Kind)
	}
}

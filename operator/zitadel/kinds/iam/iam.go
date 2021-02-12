package iam

import (
	"fmt"
	"github.com/caos/zitadel/operator/zitadel/kinds/iam/zitadel/database"

	"github.com/caos/orbos/mntr"
	"github.com/caos/orbos/pkg/labels"
	"github.com/caos/orbos/pkg/secret"
	"github.com/caos/orbos/pkg/tree"
	"github.com/caos/zitadel/operator"
	"github.com/caos/zitadel/operator/zitadel/kinds/iam/zitadel"
	"github.com/pkg/errors"
	core "k8s.io/api/core/v1"
)

func GetQueryAndDestroyFuncs(
	monitor mntr.Monitor,
	operatorLabels *labels.Operator,
	desiredTree *tree.Tree,
	currentTree *tree.Tree,
	nodeselector map[string]string,
	tolerations []core.Toleration,
	dbClient database.Client,
	action string,
	version *string,
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
	case "zitadel.caos.ch/ZITADEL":
		apiLabels := labels.MustForAPI(operatorLabels, "ZITADEL", desiredTree.Common.Version)
		return zitadel.AdaptFunc(apiLabels, nodeselector, tolerations, dbClient, action, version, features)(monitor, desiredTree, currentTree)
	default:
		return nil, nil, nil, errors.Errorf("unknown iam kind %s", desiredTree.Common.Kind)
	}
}

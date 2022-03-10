package iam

import (
	"fmt"

	core "k8s.io/api/core/v1"

	"github.com/caos/orbos/mntr"
	"github.com/caos/orbos/pkg/labels"
	"github.com/caos/orbos/pkg/secret"
	"github.com/caos/orbos/pkg/tree"

	"github.com/caos/zitadel/operator"
	"github.com/caos/zitadel/operator/zitadel/kinds/iam/zitadel"
	"github.com/caos/zitadel/operator/zitadel/kinds/iam/zitadel/database"
)

func Adapt(
	monitor mntr.Monitor,
	operatorLabels *labels.Operator,
	desiredTree *tree.Tree,
	currentTree *tree.Tree,
	nodeselector map[string]string,
	tolerations []core.Toleration,
	dbClient database.Client,
	namespace string,
	action string,
	version *string,
	features []string,
	customImageRegistry string,
) (
	query operator.QueryFunc,
	destroy operator.DestroyFunc,
	configure operator.ConfigureFunc,
	secrets map[string]*secret.Secret,
	existing map[string]*secret.Existing,
	migrate bool,
	err error,
) {

	defer func() {
		if err != nil {
			err = fmt.Errorf("adapting %s failed: %w", desiredTree.Common.Kind, err)
		}
	}()

	switch desiredTree.Common.Kind {
	case "zitadel.caos.ch/ZITADEL":
		apiLabels := labels.MustForAPI(operatorLabels, "ZITADEL", desiredTree.Common.Version())
		return zitadel.AdaptFunc(
			apiLabels,
			nodeselector,
			tolerations,
			dbClient,
			namespace,
			action,
			version,
			features,
			customImageRegistry,
		)(monitor, desiredTree, currentTree)
	default:
		return nil, nil, nil, nil, nil, false, mntr.ToUserError(fmt.Errorf("unknown iam kind %s", desiredTree.Common.Kind))
	}
}

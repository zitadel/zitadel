package cockroachdb

import (
	"fmt"

	"github.com/caos/orbos/pkg/kubernetes"

	"github.com/caos/zitadel/operator/zitadel/kinds/dbconn/cockroachdb"

	"github.com/caos/orbos/mntr"
	"github.com/caos/orbos/pkg/labels"
	"github.com/caos/orbos/pkg/secret"
	"github.com/caos/orbos/pkg/tree"

	"github.com/caos/zitadel/operator"
)

func Adapt(
	monitor mntr.Monitor,
	operatorLabels *labels.Operator,
	desiredTree *tree.Tree,
	currentTree *tree.Tree,
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

	if desiredTree == nil {
		return func(_ kubernetes.ClientInt, _ map[string]interface{}) (operator.EnsureFunc, error) {
				return func(_ kubernetes.ClientInt) error { return nil }, nil
			},
			func(_ kubernetes.ClientInt) error { return nil },
			func(_ kubernetes.ClientInt, _ map[string]interface{}, _ bool) error { return nil }, nil, nil, false, err
	}

	switch desiredTree.Common.Kind {
	case "zitadel.caos.ch/CockroachDB":
		apiLabels := labels.MustForAPI(operatorLabels, "CockroachDB", desiredTree.Common.Version())
		return cockroachdb.Adapter(apiLabels)(monitor, desiredTree, currentTree)
	default:
		return nil, nil, nil, nil, nil, false, mntr.ToUserError(fmt.Errorf("unknown iam kind %s", desiredTree.Common.Kind))
	}
}

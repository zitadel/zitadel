package databases

import (
	"fmt"
	"github.com/caos/orbos/mntr"
	"github.com/caos/orbos/pkg/labels"
	"github.com/caos/orbos/pkg/secret"
	"github.com/caos/orbos/pkg/tree"
	"github.com/caos/zitadel/operator"
	"github.com/caos/zitadel/operator/database/kinds/databases/managed"
	"github.com/caos/zitadel/operator/database/kinds/databases/provided"
)

const (
	component = "database"
)

func ComponentSelector() *labels.Selector {
	return labels.OpenComponentSelector("ZITADEL", component)
}

func Adapt(
	monitor mntr.Monitor,
	desiredTree *tree.Tree,
	currentTree *tree.Tree,
	namespace string,
	apiLabels *labels.API,
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
	componentLabels := labels.MustForComponent(apiLabels, component)
	internalMonitor := monitor.WithField("component", component)

	switch desiredTree.Common.Kind {
	case "databases.caos.ch/CockroachDB":
		return managed.Adapter(componentLabels, namespace, features, customImageRegistry)(internalMonitor, desiredTree, currentTree)
	case "databases.caos.ch/ProvidedDatabase":
		return provided.Adapter()(internalMonitor, desiredTree, currentTree)
	default:
		return nil, nil, nil, nil, nil, false, mntr.ToUserError(fmt.Errorf("unknown database kind %s: %w", desiredTree.Common.Kind, err))
	}
}

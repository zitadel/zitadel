package databases

import (
	"fmt"

	core "k8s.io/api/core/v1"

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
	timestamp string,
	nodeselector map[string]string,
	tolerations []core.Toleration,
	version string,
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
		return managed.Adapter(
			componentLabels,
			namespace,
			timestamp,
			nodeselector,
			tolerations,
			version,
			features,
			customImageRegistry,
		)(internalMonitor, desiredTree, currentTree)
	case "databases.caos.ch/ProvidedDatabase":
		return provided.Adapter()(internalMonitor, desiredTree, currentTree)
	default:
		return nil, nil, nil, nil, nil, false, mntr.ToUserError(fmt.Errorf("unknown database kind %s: %w", desiredTree.Common.Kind, err))
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
	case "databases.caos.ch/CockroachDB":
		return managed.BackupList()(monitor, desiredTree)
	case "databases.caos.ch/ProvidedDatabse":
		return nil, mntr.ToUserError(fmt.Errorf("no backups supported for database kind %s", desiredTree.Common.Kind))
	default:
		return nil, mntr.ToUserError(fmt.Errorf("unknown database kind %s", desiredTree.Common.Kind))
	}
}

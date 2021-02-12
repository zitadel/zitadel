package databases

import (
	"github.com/caos/orbos/mntr"
	"github.com/caos/orbos/pkg/labels"
	"github.com/caos/orbos/pkg/secret"
	"github.com/caos/orbos/pkg/tree"
	"github.com/caos/zitadel/operator"
	"github.com/caos/zitadel/operator/database/kinds/databases/managed"
	"github.com/caos/zitadel/operator/database/kinds/databases/provided"
	"github.com/pkg/errors"
	core "k8s.io/api/core/v1"
)

const (
	component = "database"
)

func ComponentSelector() *labels.Selector {
	return labels.OpenComponentSelector("ZITADEL", component)
}

func GetQueryAndDestroyFuncs(
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
) (
	query operator.QueryFunc,
	destroy operator.DestroyFunc,
	secrets map[string]*secret.Secret,
	err error,
) {
	componentLabels := labels.MustForComponent(apiLabels, component)
	internalMonitor := monitor.WithField("component", component)

	switch desiredTree.Common.Kind {
	case "databases.caos.ch/CockroachDB":
		return managed.AdaptFunc(componentLabels, namespace, timestamp, nodeselector, tolerations, version, features)(internalMonitor, desiredTree, currentTree)
	case "databases.caos.ch/ProvidedDatabase":
		return provided.AdaptFunc()(internalMonitor, desiredTree, currentTree)
	default:
		return nil, nil, nil, errors.Errorf("unknown database kind %s", desiredTree.Common.Kind)
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
		return nil, errors.Errorf("no backups supported for database kind %s", desiredTree.Common.Kind)
	default:
		return nil, errors.Errorf("unknown database kind %s", desiredTree.Common.Kind)
	}
}

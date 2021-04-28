package orb

import (
	"github.com/caos/orbos/mntr"
	"github.com/caos/orbos/pkg/kubernetes"
	"github.com/caos/orbos/pkg/tree"
	"github.com/caos/zitadel/operator/database/kinds/databases"
	"github.com/pkg/errors"
)

func BackupListFunc() func(monitor mntr.Monitor, k8sClient kubernetes.ClientInt, desiredTree *tree.Tree) (strings []string, err error) {
	return func(monitor mntr.Monitor, k8sClient kubernetes.ClientInt, desiredTree *tree.Tree) (strings []string, err error) {
		desiredKind, err := ParseDesiredV0(desiredTree)
		if err != nil {
			return nil, errors.Wrap(err, "parsing desired state failed")
		}
		desiredTree.Parsed = desiredKind

		if desiredKind.Spec.Verbose && !monitor.IsVerbose() {
			monitor = monitor.Verbose()
		}

		return databases.GetBackupList(monitor, k8sClient, desiredKind.Database)
	}
}

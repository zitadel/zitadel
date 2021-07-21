package orb

import (
	"fmt"

	"github.com/caos/orbos/mntr"
	"github.com/caos/orbos/pkg/tree"

	"github.com/caos/zitadel/operator/database/kinds/databases"
)

func BackupListFunc() func(monitor mntr.Monitor, desiredTree *tree.Tree) (strings []string, err error) {
	return func(monitor mntr.Monitor, desiredTree *tree.Tree) (strings []string, err error) {
		desiredKind, err := ParseDesiredV0(desiredTree)
		if err != nil {
			return nil, fmt.Errorf("parsing desired state failed: %w", err)
		}
		desiredTree.Parsed = desiredKind

		if desiredKind.Spec.Verbose && !monitor.IsVerbose() {
			monitor = monitor.Verbose()
		}

		return databases.GetBackupList(monitor, desiredKind.Database)
	}
}

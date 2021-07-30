package managed

import (
	"fmt"

	"github.com/caos/orbos/mntr"
	"github.com/caos/orbos/pkg/kubernetes"
	"github.com/caos/orbos/pkg/tree"

	"github.com/caos/zitadel/operator/database/kinds/backups"
)

func BackupList() func(monitor mntr.Monitor, k8sClient kubernetes.ClientInt, desired *tree.Tree) ([]string, error) {
	return func(monitor mntr.Monitor, k8sClient kubernetes.ClientInt, desired *tree.Tree) ([]string, error) {
		desiredKind, err := parseDesiredV0(desired)
		if err != nil {
			return nil, fmt.Errorf("parsing desired state failed: %w", err)
		}
		desired.Parsed = desiredKind

		if !monitor.IsVerbose() && desiredKind.Spec.Verbose {
			monitor.Verbose()
		}

		backuplists := make([]string, 0)
		if desiredKind.Spec.Backups != nil {
			for name, def := range desiredKind.Spec.Backups {
				backuplist, err := backups.GetBackupList(monitor, k8sClient, name, def)
				if err != nil {
					return nil, err
				}
				for _, backup := range backuplist {
					backuplists = append(backuplists, name+"."+backup)
				}
			}
		}
		return backuplists, nil
	}
}

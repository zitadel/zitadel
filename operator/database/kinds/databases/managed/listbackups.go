package managed

import (
	"github.com/caos/orbos/mntr"
	"github.com/caos/orbos/pkg/kubernetes"
	"github.com/caos/orbos/pkg/tree"
	"github.com/caos/zitadel/operator/database/kinds/backups"
	"github.com/pkg/errors"
)

func BackupList() func(monitor mntr.Monitor, k8sClient kubernetes.ClientInt, desired *tree.Tree) ([]string, error) {
	return func(monitor mntr.Monitor, k8sClient kubernetes.ClientInt, desired *tree.Tree) ([]string, error) {
		desiredKind, err := parseDesiredV0(desired)
		if err != nil {
			return nil, errors.Wrap(err, "parsing desired state failed")
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

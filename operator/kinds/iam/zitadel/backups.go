package zitadel

import (
	"github.com/caos/orbos/mntr"
	"github.com/caos/orbos/pkg/tree"
	"github.com/pkg/errors"
)

func BackupListFunc() func(monitor mntr.Monitor, desired *tree.Tree) (strings []string, err error) {
	return func(monitor mntr.Monitor, desired *tree.Tree) (strings []string, err error) {
		desiredKind, err := parseDesiredV0(desired)
		if err != nil {
			return nil, errors.Wrap(err, "parsing desired state failed")
		}
		desired.Parsed = desiredKind

		if !monitor.IsVerbose() && desiredKind.Spec.Verbose {
			monitor.Verbose()
		}

		//return databases.GetBackupList(monitor, desiredKind.Database)
		//TODO:
		return []string{}, nil
	}
}

package core

import (
	"github.com/caos/orbos/mntr"
	"github.com/caos/orbos/pkg/tree"
)

type BackupListFunc func(monitor mntr.Monitor, name string, desired *tree.Tree) ([]string, error)

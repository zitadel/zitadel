package core

import (
	"github.com/caos/orbos/mntr"
	"github.com/caos/orbos/pkg/kubernetes"
	"github.com/caos/orbos/pkg/tree"
)

type BackupListFunc func(monitor mntr.Monitor, k8sClient kubernetes.ClientInt, name string, desired *tree.Tree) ([]string, error)

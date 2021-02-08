package managed

import (
	"github.com/caos/orbos/pkg/kubernetes/k8s"
	"github.com/caos/orbos/pkg/tree"
	"github.com/pkg/errors"
	corev1 "k8s.io/api/core/v1"
)

type DesiredV0 struct {
	Common *tree.Common `yaml:",inline"`
	Spec   Spec
}

type Spec struct {
	Verbose         bool
	Force           bool                  `yaml:"force,omitempty"`
	ReplicaCount    int                   `yaml:"replicaCount,omitempty"`
	StorageCapacity string                `yaml:"storageCapacity,omitempty"`
	StorageClass    string                `yaml:"storageClass,omitempty"`
	NodeSelector    map[string]string     `yaml:"nodeSelector,omitempty"`
	Tolerations     []corev1.Toleration   `yaml:"tolerations,omitempty"`
	ClusterDns      string                `yaml:"clusterDNS,omitempty"`
	Backups         map[string]*tree.Tree `yaml:"backups,omitempty"`
	Resources       *k8s.Resources        `yaml:"resources,omitempty"`
}

func parseDesiredV0(desiredTree *tree.Tree) (*DesiredV0, error) {
	desiredKind := &DesiredV0{
		Common: desiredTree.Common,
		Spec:   Spec{},
	}

	if err := desiredTree.Original.Decode(desiredKind); err != nil {
		return nil, errors.Wrap(err, "parsing desired state failed")
	}

	return desiredKind, nil
}

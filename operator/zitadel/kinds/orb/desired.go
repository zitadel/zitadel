package orb

import (
	"github.com/caos/orbos/pkg/tree"
	"github.com/pkg/errors"
	corev1 "k8s.io/api/core/v1"
)

type DesiredV0 struct {
	Common *tree.Common `yaml:",inline"`
	Spec   struct {
		Verbose         bool
		NodeSelector    map[string]string   `yaml:"nodeSelector,omitempty"`
		Tolerations     []corev1.Toleration `yaml:"tolerations,omitempty"`
		Version         string              `yaml:"version,omitempty"`
		SelfReconciling bool                `yaml:"selfReconciling"`
		//Use this registry to pull the database operator image from
		//@default: ghcr.io
		CustomImageRegistry string `json:"customImageRegistry,omitempty" yaml:"customImageRegistry,omitempty"`
	}
	IAM *tree.Tree
}

func parseDesiredV0(desiredTree *tree.Tree) (*DesiredV0, error) {
	desiredKind := &DesiredV0{Common: desiredTree.Common}

	if err := desiredTree.Original.Decode(desiredKind); err != nil {
		return nil, errors.Wrap(err, "parsing desired state failed")
	}
	desiredKind.Common.Version = "v0"

	return desiredKind, nil
}

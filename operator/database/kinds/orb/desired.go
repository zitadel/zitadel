package orb

import (
	"fmt"

	"github.com/caos/orbos/mntr"
	"github.com/caos/orbos/pkg/tree"
	corev1 "k8s.io/api/core/v1"
)

type DesiredV0 struct {
	Common   *tree.Common `json:",inline" yaml:",inline"`
	Spec     *Spec        `json:"spec" yaml:"spec"`
	Database *tree.Tree
}

// +kubebuilder:object:generate=true
type Spec struct {
	Verbose         bool                `json:"verbose" json:"verbose"`
	NodeSelector    map[string]string   `json:"nodeSelector,omitempty" yaml:"nodeSelector,omitempty"`
	Tolerations     []corev1.Toleration `json:"tolerations,omitempty" yaml:"tolerations,omitempty"`
	Version         string              `json:"version,omitempty" yaml:"version,omitempty"`
	SelfReconciling bool                `json:"selfReconciling" yaml:"selfReconciling"`
	//Use this registry to pull container images from
	//@default: <multiple public registries>
	CustomImageRegistry string `json:"customImageRegistry,omitempty" yaml:"customImageRegistry,omitempty"`
}

func ParseDesiredV0(desiredTree *tree.Tree) (*DesiredV0, error) {
	desiredKind := &DesiredV0{Common: desiredTree.Common}

	if err := desiredTree.Original.Decode(desiredKind); err != nil {
		return nil, mntr.ToUserError(fmt.Errorf("parsing desired state failed: %w", err))
	}
	desiredKind.Common.OverwriteVersion("v0")

	return desiredKind, nil
}

package bucket

import (
	secret2 "github.com/caos/orbos/pkg/secret"
	"github.com/caos/orbos/pkg/tree"
	"github.com/pkg/errors"
)

type DesiredV0 struct {
	Common *tree.Common `yaml:",inline"`
	Spec   *Spec
}

type Spec struct {
	Verbose            bool
	Cron               string          `yaml:"cron,omitempty"`
	Bucket             string          `yaml:"bucket,omitempty"`
	ServiceAccountJSON *secret2.Secret `yaml:"serviceAccountJSON,omitempty"`
}

func (s *Spec) IsZero() bool {
	if (s.ServiceAccountJSON == nil || s.ServiceAccountJSON.IsZero()) &&
		!s.Verbose &&
		s.Cron == "" &&
		s.Bucket == "" {
		return true
	}
	return false
}

func ParseDesiredV0(desiredTree *tree.Tree) (*DesiredV0, error) {
	desiredKind := &DesiredV0{
		Common: desiredTree.Common,
		Spec:   &Spec{},
	}

	if err := desiredTree.Original.Decode(desiredKind); err != nil {
		return nil, errors.Wrap(err, "parsing desired state failed")
	}

	return desiredKind, nil
}

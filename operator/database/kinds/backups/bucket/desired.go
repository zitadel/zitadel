package bucket

import (
	"fmt"

	"github.com/caos/orbos/pkg/secret"
	"github.com/caos/orbos/pkg/tree"
	"github.com/pkg/errors"
)

type DesiredV0 struct {
	Common *tree.Common `yaml:",inline"`
	Spec   *Spec
}

type Spec struct {
	Verbose                    bool
	Cron                       string           `yaml:"cron,omitempty"`
	Bucket                     string           `yaml:"bucket,omitempty"`
	ServiceAccountJSON         *secret.Secret   `yaml:"serviceAccountJSON,omitempty"`
	ExistingServiceAccountJSON *secret.Existing `yaml:"existingServiceAccountJSON,omitempty"`
}

func (s *Spec) IsZero() bool {
	if ((s.ServiceAccountJSON == nil || s.ServiceAccountJSON.IsZero()) && (s.ExistingServiceAccountJSON == nil || s.ExistingServiceAccountJSON.IsZero())) &&
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

func (d *DesiredV0) validateSecrets() error {
	if err := secret.ValidateSecret(d.Spec.ServiceAccountJSON, d.Spec.ExistingServiceAccountJSON); err != nil {
		return fmt.Errorf("validating api key failed: %w", err)
	}
}

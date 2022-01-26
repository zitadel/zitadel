package cockroachdb

import (
	"errors"
	"fmt"

	"github.com/caos/orbos/pkg/secret"

	"github.com/caos/orbos/mntr"
	"github.com/caos/orbos/pkg/tree"
)

type DesiredV0 struct {
	Common *tree.Common `yaml:",inline"`
	Spec   *Spec
}

type Spec struct {
	Verbose             bool
	URL                 string
	Port                string
	Certificate         *secret.Secret   `yaml:"certificate,omitempty"`
	ExistingCertificate *secret.Existing `yaml:"existingCertificate,omitempty"`
}

func parseDesiredV0(desiredTree *tree.Tree) (*DesiredV0, error) {
	desiredKind := &DesiredV0{
		Common: desiredTree.Common,
		Spec:   &Spec{},
	}

	if err := desiredTree.Original.Decode(desiredKind); err != nil {
		return nil, mntr.ToUserError(fmt.Errorf("parsing desired state failed: %w", err))
	}

	return desiredKind, nil
}

func (d *DesiredV0) validate() (err error) {

	defer func() {
		if err != nil {
			err = fmt.Errorf("database connection spec is invalid: %w", err)
		}
	}()

	if d.Spec == nil {
		return errors.New("spec is empty")
	}

	if d.Spec.URL == "" {
		return errors.New("url is empty")
	}

	if d.Spec.Port == "" {
		return errors.New("port is empty")
	}

	return nil
}

func (d *DesiredV0) validateSecrets() (err error) {
	if err := secret.ValidateSecret(d.Spec.Certificate, d.Spec.ExistingCertificate); err != nil {
		return fmt.Errorf("validating certificate failed: %w", err)
	}
	return nil
}

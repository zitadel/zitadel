package s3

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
	Verbose                 bool
	Cron                    string           `yaml:"cron,omitempty"`
	Bucket                  string           `yaml:"bucket,omitempty"`
	Endpoint                string           `yaml:"endpoint,omitempty"`
	Region                  string           `yaml:"region,omitempty"`
	AccessKeyID             *secret.Secret   `yaml:"accessKeyID,omitempty"`
	ExistingAccessKeyID     *secret.Existing `yaml:"existingAccessKeyID,omitempty"`
	SecretAccessKey         *secret.Secret   `yaml:"secretAccessKey,omitempty"`
	ExistingSecretAccessKey *secret.Existing `yaml:"existingSecretAccessKey,omitempty"`
	SessionToken            *secret.Secret   `yaml:"sessionToken,omitempty"`
	ExistingSessionToken    *secret.Existing `yaml:"existingSessionToken,omitempty"`
}

func (s *Spec) IsZero() bool {
	if ((s.AccessKeyID == nil || s.AccessKeyID.IsZero()) && (s.ExistingAccessKeyID == nil || s.ExistingAccessKeyID.IsZero())) &&
		((s.SecretAccessKey == nil || s.SecretAccessKey.IsZero()) && (s.ExistingSecretAccessKey == nil || s.ExistingSecretAccessKey.IsZero())) &&
		((s.SessionToken == nil || s.SessionToken.IsZero()) && (s.ExistingSessionToken == nil || s.ExistingSessionToken.IsZero())) &&
		!s.Verbose &&
		s.Bucket == "" &&
		s.Cron == "" &&
		s.Endpoint == "" &&
		s.Region == "" {
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
	if err := secret.ValidateSecret(d.Spec.AccessKeyID, d.Spec.ExistingAccessKeyID); err != nil {
		return fmt.Errorf("validating access key id failed: %w", err)
	}
	if err := secret.ValidateSecret(d.Spec.SecretAccessKey, d.Spec.ExistingSecretAccessKey); err != nil {
		return fmt.Errorf("validating secret access key failed: %w", err)
	}
	if err := secret.ValidateSecret(d.Spec.SessionToken, d.Spec.ExistingSessionToken); err != nil {
		return fmt.Errorf("validating session token failed: %w", err)
	}
	return nil
}

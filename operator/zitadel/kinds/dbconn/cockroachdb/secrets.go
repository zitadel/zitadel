package cockroachdb

import (
	"github.com/caos/orbos/pkg/secret"
)

func getSecretsMap(desiredKind *DesiredV0) (
	map[string]*secret.Secret,
	map[string]*secret.Existing,
) {

	var (
		secrets  = map[string]*secret.Secret{}
		existing = map[string]*secret.Existing{}
	)

	if desiredKind.Spec == nil {
		desiredKind.Spec = &Spec{}
	}

	if desiredKind.Spec.Certificate == nil {
		desiredKind.Spec.Certificate = &secret.Secret{}
	}
	if desiredKind.Spec.ExistingCertificate == nil {
		desiredKind.Spec.ExistingCertificate = &secret.Existing{}
	}
	certKey := "certificate"
	secrets[certKey] = desiredKind.Spec.Certificate
	existing[certKey] = desiredKind.Spec.ExistingCertificate

	return secrets, existing
}

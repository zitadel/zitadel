package managed

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

	if desiredKind.Spec.ZitadelUserPassword == nil {
		desiredKind.Spec.ZitadelUserPassword = &secret.Secret{}
	}

	if desiredKind.Spec.ZitadelUserPasswordExisting == nil {
		desiredKind.Spec.ZitadelUserPasswordExisting = &secret.Existing{}
	}

	pwKey := "zitadeluserpassword"
	secrets[pwKey] = desiredKind.Spec.ZitadelUserPassword
	existing[pwKey] = desiredKind.Spec.ZitadelUserPasswordExisting

	return secrets, existing
}

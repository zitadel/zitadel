package bucket

import (
	"github.com/caos/orbos/pkg/secret"
)

func getSecretsMap(desiredKind *DesiredV0) (map[string]*secret.Secret, map[string]*secret.Existing) {

	var (
		secrets  = make(map[string]*secret.Secret, 0)
		existing = make(map[string]*secret.Existing, 0)
	)
	if desiredKind.Spec == nil {
		desiredKind.Spec = &Spec{}
	}

	if desiredKind.Spec.ServiceAccountJSON == nil {
		desiredKind.Spec.ServiceAccountJSON = &secret.Secret{}
	}

	if desiredKind.Spec.ExistingServiceAccountJSON == nil {
		desiredKind.Spec.ExistingServiceAccountJSON = &secret.Existing{}
	}

	sakey := "serviceaccountjson"
	secrets[sakey] = desiredKind.Spec.ServiceAccountJSON
	existing[sakey] = desiredKind.Spec.ExistingServiceAccountJSON

	return secrets, existing
}

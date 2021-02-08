package bucket

import (
	"github.com/caos/orbos/pkg/secret"
)

func getSecretsMap(desiredKind *DesiredV0) map[string]*secret.Secret {
	secrets := make(map[string]*secret.Secret, 0)
	if desiredKind.Spec == nil {
		desiredKind.Spec = &Spec{}
	}

	if desiredKind.Spec.ServiceAccountJSON == nil {
		desiredKind.Spec.ServiceAccountJSON = &secret.Secret{}
	}
	secrets["serviceaccountjson"] = desiredKind.Spec.ServiceAccountJSON

	return secrets
}

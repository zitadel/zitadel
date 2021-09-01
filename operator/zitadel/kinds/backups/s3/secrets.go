package s3

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

	if desiredKind.Spec.AccessKeyID == nil {
		desiredKind.Spec.AccessKeyID = &secret.Secret{}
	}

	if desiredKind.Spec.ExistingAccessKeyID == nil {
		desiredKind.Spec.ExistingAccessKeyID = &secret.Existing{}
	}

	akikey := "accesskeyid"
	secrets[akikey] = desiredKind.Spec.AccessKeyID
	existing[akikey] = desiredKind.Spec.ExistingAccessKeyID

	if desiredKind.Spec.SecretAccessKey == nil {
		desiredKind.Spec.SecretAccessKey = &secret.Secret{}
	}

	if desiredKind.Spec.ExistingSecretAccessKey == nil {
		desiredKind.Spec.ExistingSecretAccessKey = &secret.Existing{}
	}

	sakkey := "secretaccesskey"
	secrets[sakkey] = desiredKind.Spec.SecretAccessKey
	existing[sakkey] = desiredKind.Spec.ExistingSecretAccessKey

	if desiredKind.Spec.SessionToken == nil {
		desiredKind.Spec.SessionToken = &secret.Secret{}
	}

	if desiredKind.Spec.ExistingSessionToken == nil {
		desiredKind.Spec.ExistingSessionToken = &secret.Existing{}
	}

	stkey := "sessiontoken"
	secrets[stkey] = desiredKind.Spec.SessionToken
	existing[stkey] = desiredKind.Spec.ExistingSessionToken

	return secrets, existing
}

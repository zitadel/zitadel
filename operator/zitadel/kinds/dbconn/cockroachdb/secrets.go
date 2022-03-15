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
	certEntry := "certificate"
	secrets[certEntry] = desiredKind.Spec.Certificate
	existing[certEntry] = desiredKind.Spec.ExistingCertificate

	if desiredKind.Spec.CertificateKey == nil {
		desiredKind.Spec.CertificateKey = &secret.Secret{}
	}
	if desiredKind.Spec.ExistingCertificateKey == nil {
		desiredKind.Spec.ExistingCertificateKey = &secret.Existing{}
	}
	certKeyEntry := "certificatekey"
	secrets[certKeyEntry] = desiredKind.Spec.CertificateKey
	existing[certKeyEntry] = desiredKind.Spec.ExistingCertificateKey

	if desiredKind.Spec.Password == nil {
		desiredKind.Spec.Password = &secret.Secret{}
	}
	if desiredKind.Spec.ExistingPassword == nil {
		desiredKind.Spec.ExistingPassword = &secret.Existing{}
	}
	pwEntry := "password"
	secrets[pwEntry] = desiredKind.Spec.Password
	existing[pwEntry] = desiredKind.Spec.ExistingPassword

	return secrets, existing
}

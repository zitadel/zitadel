package saml

import (
	"encoding/base64"
	"github.com/caos/zitadel/internal/api/saml/serviceprovider"
	"github.com/caos/zitadel/internal/api/saml/signature"
	"github.com/caos/zitadel/internal/api/saml/xml/md"
	"github.com/caos/zitadel/internal/api/saml/xml/samlp"
	"github.com/caos/zitadel/internal/api/saml/xml/xml_dsig"
	"reflect"
)

func signaturePostProvided(
	signatureF func() *xml_dsig.SignatureType,
) func() bool {
	return func() bool {
		signatureV := signatureF()

		return signatureV != nil &&
			!reflect.DeepEqual(signatureV.SignatureValue, xml_dsig.SignatureValueType{}) &&
			signatureV.SignatureValue.Text != ""
	}
}
func signaturePostVerificationNecessary(
	idpMetadataF func() *md.IDPSSODescriptorType,
	spMetadataF func() *md.EntityDescriptorType,
	signatureF func() *xml_dsig.SignatureType,
	protocolBinding func() string,
) func() bool {
	return func() bool {
		spMeta := spMetadataF()
		idpMeta := idpMetadataF()

		return ((spMeta == nil || spMeta.SPSSODescriptor == nil || spMeta.SPSSODescriptor.AuthnRequestsSigned == "true") ||
			(idpMeta == nil || idpMeta.WantAuthnRequestsSigned == "true") ||
			signaturePostProvided(signatureF)()) &&
			protocolBinding() == PostBinding
	}
}

func verifyPostSignature(
	authRequestF func() string,
	spF func() *serviceprovider.ServiceProvider,
	errF func(error),
) func() error {
	return func() error {
		sp := spF()

		data, err := base64.StdEncoding.DecodeString(authRequestF())
		if err != nil {
			errF(err)
			return err
		}

		if err := sp.ValidatePostSignature(string(data)); err != nil {
			errF(err)
			return err
		}
		return nil
	}
}

func createPostSignature(
	samlResponse *samlp.ResponseType,
	idp *IdentityProvider,
) error {
	sig, err := signature.Create(idp.signer, samlResponse.Assertion)
	if err != nil {
		return err
	}

	samlResponse.Assertion.Signature = sig
	return nil
}

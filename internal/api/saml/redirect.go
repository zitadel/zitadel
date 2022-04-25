package saml

import (
	"encoding/base64"
	"fmt"
	"github.com/caos/zitadel/internal/api/saml/serviceprovider"
	"github.com/caos/zitadel/internal/api/saml/signature"
	"github.com/caos/zitadel/internal/api/saml/xml"
	"github.com/caos/zitadel/internal/api/saml/xml/md"
	"github.com/caos/zitadel/internal/api/saml/xml/samlp"
	"net/url"
)

func signatureRedirectVerificationNecessary(
	idpMetadataF func() *md.IDPSSODescriptorType,
	spMetadataF func() *md.EntityDescriptorType,
	signatureF func() string,
	protocolBinding func() string,
) func() bool {
	return func() bool {
		spMeta := spMetadataF()
		idpMeta := idpMetadataF()

		return ((spMeta == nil || spMeta.SPSSODescriptor == nil || spMeta.SPSSODescriptor.AuthnRequestsSigned == "true") ||
			(idpMeta == nil || idpMeta.WantAuthnRequestsSigned == "true") ||
			signatureF() != "") &&
			protocolBinding() == RedirectBinding
	}
}

func verifyRedirectSignature(
	authRequest func() string,
	relayState func() string,
	sig func() string,
	sigAlg func() string,
	sp func() *serviceprovider.ServiceProvider,
	errF func(error),
) func() error {
	return func() error {
		if authRequest() == "" {
			return fmt.Errorf("no authrequest provided but required")
		}
		if relayState() == "" {
			return fmt.Errorf("no relaystate provided but required")
		}
		if sig() == "" {
			return fmt.Errorf("no signature provided but required")
		}
		if sigAlg() == "" {
			return fmt.Errorf("no signature algorithm provided but required")
		}

		spInstance := sp()
		if sp == nil {
			return fmt.Errorf("no service provider instance provided but required")
		}

		err := spInstance.ValidateRedirectSignature(
			authRequest(),
			relayState(),
			sigAlg(),
			sig(),
		)
		errF(err)
		return err
	}
}

func createRedirectSignature(
	samlResponse *samlp.ResponseType,
	idp *IdentityProvider,
	response *Response,
) error {
	respStr, err := xml.Marshal(samlResponse)
	if err != nil {
		return err
	}

	respData, err := xml.DeflateAndBase64([]byte(respStr))
	if err != nil {
		return err
	}

	sigAlg := idp.signingContext.GetSignatureMethodIdentifier()
	sig, err := signature.CreateRedirect(idp.signingContext, buildRedirectQuery(string(respData), response.RelayState, sigAlg, ""))
	if err != nil {
		return err
	}

	response.Signature = url.QueryEscape(base64.StdEncoding.EncodeToString(sig))
	response.SigAlg = url.QueryEscape(base64.StdEncoding.EncodeToString([]byte(sigAlg)))
	return nil
}

func buildRedirectQuery(
	response string,
	relayState string,
	sigAlg string,
	sig string,
) string {
	query := "SAMLResponse=" + url.QueryEscape(response)
	if relayState != "" {
		query += "&RelayState=" + url.QueryEscape(relayState)
	}
	if sig != "" {
		query += "&Signature=" + url.QueryEscape(sig)
	}
	if sigAlg != "" {
		query += "&SigAlg=" + url.QueryEscape(sigAlg)
	}

	return query
}
